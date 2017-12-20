package subcmd

import (
	"strings"
)

func parseCmd(arg string) (string, string) {
	at := strings.Index(arg, "@")
	if at >= 0 {
		cmd := arg[:at]
		host := arg[at+1:]
		return cmd, host
	}

	return arg, ""
}
