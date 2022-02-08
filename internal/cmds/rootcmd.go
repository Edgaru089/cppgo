package cmds

import (
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "cppgo [flags] command ...",
		Short: "Use the C preprocessor in Go source files",
		Long: `This tool invokes the C preprocessor ("cpp") for files in the go module path
with the extension "cppgo", among other things. Run "cppgo doc" for more info.`,
		DisableFlagsInUseLine: true,
	}

	commands = []*cobra.Command{
		cmdGenerate,
		cmdBuild,
		cmdRun,
		cmdInfo,
	}
)

func init() {
	root.PersistentFlags().BoolVarP(&Flags.Verbose, "verbose", "v", false, "be verbose whenever cppgo does something")
	root.PersistentFlags().StringSliceVar(&Flags.Extensions, "extensions", []string{"cppgo", "pgo"}, "file extensions cppgo should preprocess (comma-separated)")
	root.PersistentFlags().StringVar(&Flags.GoBinary, "go", "go", "use this go command instead")
	root.PersistentFlags().StringVar(&Flags.CppBinary, "cpp", "cpp", "use this C preprocessor instead")

	root.AddCommand(commands...)
	root.SilenceUsage = true
}

func Execute() error {
	return root.Execute()
}
