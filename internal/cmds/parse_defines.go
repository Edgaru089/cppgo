package cmds

import "strings"

func parseDefines(strs []string, defines map[string]string) {

	for _, s := range strs {
		parts := strings.Split(s, ",")
		for _, s := range parts {
			i := strings.IndexByte(s, '=')
			if i == -1 {
				defines[s] = "1"
			} else {
				defines[s[:i]] = s[i+1:]
			}
		}
	}
}
