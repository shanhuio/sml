package smake // import "shanhu.io/sml/smake"

import (
	"fmt"
	"os"
)

func smake() error {
	gopath, err := absGOPATH()
	if err != nil {
		return err
	}
	c := newContext(gopath, ".")
	return c.smake()
}

// Main is the entry point for smake.
func Main() {
	if err := smake(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
