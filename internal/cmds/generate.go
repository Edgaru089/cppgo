package cmds

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// The generate command.
	cmdGenerate = &cobra.Command{
		Use:   "generate [options] [path]",
		Short: "Generate a module/package with cppgo source files",
		Long: `generate updates all the changed cppgo source files in the go module path
given. See "cppgo doc" for more info.`,
		DisableFlagsInUseLine: true,
		PreRun:                func(cmd *cobra.Command, args []string) {},
		RunE:                  RunGenerate,
	}

	rawDefines    []string
	defines       map[string]string
	workdir       string
	generateFiles chan string
)

func init() {
	cmdGenerate.Flags().StringSliceVarP(&rawDefines, "define", "D", nil, "pre-defined macros (as NAME=VAL,NAME=VAL)")
}

// RunGenerate executes the build command.
func RunGenerate(_ *cobra.Command, args []string) error {

	defines = make(map[string]string)
	parseDefines(rawDefines, defines)

	// Insert the default macros
	if _, ok := defines[defaultGOOS()]; !ok {
		defines[defaultGOOS()] = "1"
	}
	if _, ok := defines[defaultGOARCH()]; !ok {
		defines[defaultGOARCH()] = "1"
	}
	if _, ok := defines["CGO_ENABLED"]; !ok {
		defines["CGO_ENABLED"] = defaultCGO_ENABLED()
	}
	if _, ok := defines["GOVERSION"]; !ok {
		defines["GOVERSION"] = defaultGOVERSION()
	}

	if Flags.Verbose {
		fmt.Printf("generate: Verbose: %v, Extensions: %v, Macros:%v\n", Flags.Verbose, Flags.Extensions, defines)
	}

	// Get the working directory
	if len(args) > 0 {
		workdir = args[0]
	} else {
		var err error
		workdir, err = os.Getwd()
		if err != nil {
			return errors.New("generate: working directory cannot be retrieved and none specified")
		}
	}
	// Look into the parents for go.mod
	for parent := workdir; parent != filepath.Dir(parent); parent = filepath.Dir(parent) {
		if _, err := os.Stat(filepath.Join(parent, "go.mod")); err == nil {
			workdir = parent
			break
		}
	}
	if Flags.Verbose {
		fmt.Println("FinWorkdir:", workdir)
	}
	// We should now have the working directory

	generateFiles = make(chan string, 2)
	// Search the folder for files
	go func() {
		filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
			base := filepath.Base(path)
			if d.IsDir() && len(base) > 1 && base[0] == '.' {
				return filepath.SkipDir // (Should be) hidden
			}

			if !d.IsDir() && strings.IndexByte(base, '.') != -1 {
				ext := filepath.Ext(base)[1:]
				for _, e := range Flags.Extensions {
					if ext == e {
						generateFiles <- path
						break
					}
				}
			}

			return nil
		})
		close(generateFiles)
	}()

	// Receive & generate the files
	for file := range generateFiles {

		base := filepath.Base(file)
		dir := filepath.Dir(file)
		output := filepath.Join(dir, fmt.Sprintf("%s_%s.go", filepath.Ext(base)[1:], base[:strings.LastIndexByte(base, '.')]))

		if Flags.Verbose {
			fmt.Println("Generate:", file, "-->", filepath.Base(output))
		}

		// Construct the command-line
		args := []string{"-P", "-C", "-nostdinc", file, "-o", output}
		for key, val := range defines {
			args = append(args, "-D", key+"="+val)
		}
		cmd := exec.Command(Flags.CppBinary, args...)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			for range generateFiles {
				// Drain the channel
			}
			root.SilenceUsage = true
			return err
		}
	}

	return nil
}
