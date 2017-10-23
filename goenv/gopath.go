package goenv

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// GOPATH returns GOPATH reading from environment variables.
// If GOPATH is missing it returns $HOME/go.
func GOPATH() (string, error) {
	p := os.Getenv("GOPATH")
	if p != "" {
		lst := filepath.SplitList(p)
		if len(lst) > 1 {
			return "", fmt.Errorf("GOPATH contains multiple folders")
		}
		return p, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, "go"), nil
}
