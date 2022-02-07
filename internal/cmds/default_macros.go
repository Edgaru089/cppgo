package cmds

import (
	"bytes"
	"fmt"
	"go/build"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// defaultGOOS returns the name of the GOOS_XXX macro from the environment.
//
// It uses the "go/build" package.
func defaultGOOS() string {
	return "GOOS_" + strings.ToUpper(build.Default.GOOS)
}

// defaultGOARCH returns the name of the GOARCH_XXX macro from the environment.
//
// It uses the "go/build" package.
func defaultGOARCH() string {
	return "GOARCH_" + strings.ToUpper(build.Default.GOARCH)
}

// defaultCGO_ENABLED returns the value of CGO_ENABLED (either "0" or "1").
//
// It uses the "go/build" package.
func defaultCGO_ENABLED() string {
	if build.Default.CgoEnabled {
		return "1"
	} else {
		return "0"
	}
}

var (
	goversion struct {
		string
		sync.RWMutex
	}
)

// defaultGOVERSION returns the value of GOVERSION (like "11706" for 1.17.6, or 10803 for 1.8.3)
func defaultGOVERSION() string {
	// Oh no, we have to run "go version" and parse its output...

	goversion.RLock()
	if len(goversion.string) == 0 {
		goversion.RUnlock()
		goversion.Lock()
		defer goversion.Unlock()

		cmd := exec.Command(Flags.GoBinary, "version")
		outbytes, err := cmd.Output()
		if err != nil {
			cobra.CheckErr("cmd.Output(): " + err.Error())
		}

		// Find the first number in the output
		i := bytes.IndexFunc(outbytes, func(r rune) bool { return r >= '0' && r <= '9' })

		// Get every digit
		var versions [3]int
		for j := 0; i < len(outbytes) && outbytes[i] != ' '; i++ {
			if outbytes[i] == '.' {
				j++
			} else {
				versions[j] = versions[j]*10 + int(outbytes[i]) - '0'
			}
		}

		// Good-old sprintf()
		goversion.string = fmt.Sprintf("%d%02d%02d", versions[0], versions[1], versions[2])
		return goversion.string
	}

	defer goversion.RUnlock()
	return goversion.string
}
