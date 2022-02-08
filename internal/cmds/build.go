package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// The build command.
	cmdBuild = &cobra.Command{
		Use:   "build [options] [path] [ -- [go-build options]]",
		Short: "Generate & build a package with cppgo source files",
		Long: `build updates all the changed cppgo source files in the go module path
given (just like "cppgo generate"), and invokes "go build" at the given
directory. Arguments after a double dash "--" are passed to "go build".`,
		DisableFlagsInUseLine: true,
		RunE:                  RunBuildRun,
	}

	cmdRun = &cobra.Command{
		Use:   "run [options] [path] [ -- [go-build options]]",
		Short: "Generate & run a package with cppgo source files",
		Long: `run updates all the changed cppgo source files in the go module path
given (just like "cppgo generate"), and invokes "go run ." at the given
directory. Arguments after a double dash "--" are passed to "go run".`,
		DisableFlagsInUseLine: true,
		RunE:                  RunBuildRun,
	}
)

func init() {
	cmdBuild.Flags().StringSliceVarP(&rawDefines, "define", "D", nil, "pre-defined macros (as NAME=VAL,NAME=VAL)")
	cmdRun.Flags().StringSliceVarP(&rawDefines, "define", "D", nil, "pre-defined macros (as NAME=VAL,NAME=VAL)")
}

// RunBuildRun executes the build and run command.
func RunBuildRun(cmd *cobra.Command, args []string) error {

	// Split the target directory and "go build" flags
	var dirs, flags []string
	if strings.HasPrefix(cmd.Use, "build") {
		flags = []string{"build"}
	} else if strings.HasPrefix(cmd.Use, "run") {
		flags = []string{"run", "."}
	}

	isdirs := true
	for _, str := range args {
		if len(str) > 0 {
			if str[0] == '-' {
				isdirs = false
			}

			if isdirs {
				dirs = append(dirs, str)
			} else {
				flags = append(flags, str)
			}
		}
	}

	// Generate
	err := RunGenerate(nil, dirs)
	if err != nil {
		return err
	}

	// Build
	var builddir string
	if len(dirs) > 0 {
		builddir = dirs[0]
	} else {
		builddir, err = os.Getwd()
		if err != nil {
			return errors.New("build: working directory cannot be retrieved and none specified")
		}
	}

	if Flags.Verbose {
		fmt.Printf("Build/Run, At: %s, Invoke: go %v\n", builddir, flags)
	}

	gocmd := exec.Command(Flags.GoBinary, flags...)
	gocmd.Dir = builddir
	gocmd.Stdout = os.Stdout
	gocmd.Stderr = os.Stderr
	cobra.CheckErr(gocmd.Run())

	return nil
}
