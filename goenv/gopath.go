package goenv

import (
	"os"
	"os/user"
	"path/filepath"
)

// GOPATH returns GOPATH reading from environment variables.
// If GOPATH is missing it returns $HOME/go.
func GOPATH() (string, error) {
	p := os.Getenv("GOPATH")
	if p != "" {
		return p, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, "go"), nil
}
