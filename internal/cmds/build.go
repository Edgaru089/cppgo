package cmds

import (
	"fmt"

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
		Run:                   RunBuild,
	}
)

// RunBuild executes the build command.
func RunBuild(cmd *cobra.Command, args []string) {
	fmt.Printf("Hi!\n  Verbose: %v, Extensions: %v\n  Args:", Flags.Verbose, Flags.Extensions)
	for _, s := range args {
		fmt.Printf(` "%s"`, s)
	}
	fmt.Println()
}
