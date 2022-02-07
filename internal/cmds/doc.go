package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// The doc command.
	cmdDoc = &cobra.Command{
		Use:                   "doc",
		Short:                 "Tell me everything about this cursed and condemned program",
		DisableFlagsInUseLine: true,
		Run:                   RunDoc,
	}
)

// RunDoc executes the doc command.
func RunDoc(cmd *cobra.Command, args []string) {
	fmt.Println(`This is just bad. Don't use this unless you've used macros in C/C++ for 3+ years and
you're absolutely crazy (like me). While easy to abuse, the C preprocessor surely
helps solving simple problems like inline condiction-compiling and repeated code
(like if err!=nil {return err}).

When invoked, the "generate" command locates the go.mod file in the current or parent 
directories (if none can be found, it assumes the current directory as working root).
From there, it searches the directories for .cppgo/.pgo files (or whatever extension
you may specify), and, if its timestamp is newer than the generated Go source file,
it invokes the C preprocessor ("cpp"):

    $ cpp -P -C -nostdinc <input file> -o <output file> [-D <predefined macros> ... ]

The -C flag preserves comments (important!), and the -P flag disables the "#line" directives.
Note that the "github.com/spf13/pflag" package works a little differently, and you need to
write "-D MACRO" or "-D MACRO=1234" instead of "-DMACRO".

The generated Go source files take the name of <extension>_<filename>.go in
the same directory as their originals. for example, "source.cppgo" is preprocessed
into "cppgo_source.go".

cppgo defines these macros for you, that you can (but shouldn't) override:
    - GOARCH_<arch> and GOOS_<os>: <arch> and <os> are from GOARCH and GOOS,
      but are in all caps: like GOARCH_AMD64 and GOOS_LINUX.
	- GOVERSION: "go" and "." are removed, and the last 2 numbers are padded to
	  2 digits: like 10807 for go1.8.7, and 11705 for go1.17.5.
	- CGO_ENABLED: Either 0 or 1.


The "build" command takes and does everything just like "generate", but it invokes
"go build" after everything's generated. It also passes every argument after "--"
to "go build" as-is.

The "run" command works just like "build", except it invokes "go run" instead of
"go build".

The "clean" command removes all generated Go source files, if an original can be
found in the same directory.

Don't, just don't, use this in any public code. If you really wish to do so however,
you can either put "cppgo_*.go" into .gitignore, or write a //go:generate line (or both).`)
}
