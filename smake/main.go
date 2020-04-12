package smake // import "shanhu.io/sml/smake"

import (
	"fmt"
	"os"
	"path/filepath"
)

func smake() error {
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(d)
	if err != nil {
		return err
	}
	realAbs, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return err
	}
	if err := os.Chdir(realAbs); err != nil {
		return err
	}

	gopath, err := absGOPATH()
	if err != nil {
		return err
	}
	c := newContext(gopath, realAbs)
	return c.smake()
}

// Main is the entry point for smake.
func Main() {
	if err := smake(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
