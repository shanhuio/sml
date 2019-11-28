package smake

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"shanhu.io/misc/goload"
)

type context struct {
	gopath string
	dir    string
	env    []string
}

func newContext(gopath, dir string) *context {
	env := []string{fmt.Sprintf("GOPATH=%s", gopath)}
	if s := os.Getenv("PATH"); s != "" {
		env = append(env, fmt.Sprintf("PATH=%s", s))
	}
	if s := os.Getenv("HOME"); s != "" {
		env = append(env, fmt.Sprintf("HOME=%s", s))
	}

	return &context{
		gopath: gopath,
		dir:    dir,
		env:    env,
	}
}

func (c *context) srcRoot() string {
	return filepath.Join(c.gopath, "src")
}

func (c *context) execPkgs(pkgs []string, tasks [][]string) error {
	for _, args := range tasks {
		line := strings.Join(args, " ")
		fmt.Println(line)

		if len(pkgs) > 0 {
			args = append(args, pkgs...)
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

func (c *context) smake() error {
	rootPkg, err := pkgFromDir(c.srcRoot(), c.dir)
	if err != nil {
		return err
	}

	pkgs, err := goload.ListPkgs(rootPkg)
	if err != nil {
		return err
	}

	pkgs, err = relPkgs(rootPkg, pkgs)
	if err != nil {
		return err
	}

	if err := c.execPkgs(pkgs, [][]string{
		{"go", "install"},
		{"gofmt", "-s", "-w", "-l"},
		{"go", "install", "-i"},
	}); err != nil {
		return err
	}

	if err := c.execPkgs(nil, [][]string{
		{"smlchk", fmt.Sprintf("-path=%s", rootPkg)},
	}); err != nil {
		return err
	}

	return c.execPkgs(pkgs, [][]string{
		{"golint"},
		{"go", "vet"},
		{"gotags", "-R", "-f=tags"},
	})
}
