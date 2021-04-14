package smake // import "shanhu.io/sml/smake"

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"shanhu.io/misc/errcode"
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

func usingGoMod() bool {
	v, ok := os.LookupEnv("GO111MODULE")
	if !ok {
		return false
	}
	return strings.ToLower(v) != "off"
}

func run() error {
	wd, err := workDir()
	if err != nil {
		return err
	}

	mod := usingGoMod()
	if mod {
		root, err := findGoModuleRoot(wd)
		if err != nil {
			return errcode.Annotate(err, "find module root")
		}
		wd = root
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}

	gopath, err := absGOPATH()
	if err != nil {
		return err
	}

	c := newContext(gopath, wd, usingGoMod())
	return smake(c)
}

// Main is the entry point for smake.
func Main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
