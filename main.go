package main

import (
	"os"

	"github.com/Edgaru089/cppgo/internal/cmds"
)

func main() {
	err := cmds.Execute()
	if err != nil {
		os.Exit(1)
	}
}
