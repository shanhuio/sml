package smake

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type context struct {
	gopath string
	dir    string
	env    []string
	errLog io.Writer
}

func newContext(gopath, dir string, gomod bool) *context {
	env := []string{fmt.Sprintf("GOPATH=%s", gopath)}
	if s := os.Getenv("PATH"); s != "" {
		env = append(env, fmt.Sprintf("PATH=%s", s))
	}
	if s := os.Getenv("HOME"); s != "" {
		env = append(env, fmt.Sprintf("HOME=%s", s))
	}
	if !gomod {
		env = append(env, "GO111MODULE=off")
	} else {
		env = append(env, "GO111MODULE=on")
	}

	return &context{
		gopath: gopath,
		dir:    dir,
		env:    env,
		errLog: os.Stderr,
	}
}

func (c *context) srcRoot() string {
	return filepath.Join(c.gopath, "src")
}

func (c *context) execPkgs(pkgs []*relPkg, tasks [][]string) error {
	for _, args := range tasks {
		if len(args) == 0 {
			continue
		}

		line := strings.Join(args, " ")
		fmt.Println(line)

		if len(pkgs) > 0 {
			for _, pkg := range pkgs {
				args = append(args, pkg.rel)
			}
		}
		p, err := exec.LookPath(args[0])
		if err != nil {
			return err
		}
		cmd := exec.Cmd{
			Path:   p,
			Args:   args,
			Dir:    c.dir,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Env:    c.env,
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *context) logf(f string, args ...interface{}) {
	fmt.Fprintf(c.errLog, f, args...)
}

func (c *context) logln(args ...interface{}) {
	fmt.Fprintln(c.errLog, args...)
}
