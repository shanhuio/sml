package smake // import "shanhu.io/sml/smake"

import (
	"fmt"
	"os"
	"path/filepath"
)

func workDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}

func run() error {
	wd, err := workDir()
	if err != nil {
		return err
	}
	if err := os.Chdir(wd); err != nil {
		return err
	}

	gopath, err := absGOPATH()
	if err != nil {
		return err
	}
	c := newContext(gopath, wd)
	return smake(c)
}

// Main is the entry point for smake.
func Main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
