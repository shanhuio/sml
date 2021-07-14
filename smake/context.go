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
	mod    bool
}

func newContext(gopath, dir string, gomod bool) *context {
	env := []string{fmt.Sprintf("GOPATH=%s", gopath)}
	for _, v := range []string{
		"PATH", "HOME", "SSH_AUTH_SOCK",
	} {
		if s := os.Getenv(v); s != "" {
			env = append(env, fmt.Sprintf("%s=%s", v, s))
		}
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
		mod:    gomod,
	}
}

func (c *context) workDir() string { return c.dir }

func (c *context) gomod() bool { return c.mod }

func (c *context) srcRoot() string { return filepath.Join(c.gopath, "src") }

type execConfig struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (c *context) execPkgs(
	pkgs []*relPkg, args []string, config *execConfig,
) error {
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

	var stdout, stderr io.Writer
	if config != nil {
		stdout, stderr = config.Stdout, config.Stderr
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	cmd := exec.Cmd{
		Path:   p,
		Args:   args,
		Dir:    c.dir,
		Stdout: stdout,
		Stderr: stderr,
		Env:    c.env,
	}
	return cmd.Run()
}

func (c *context) logf(f string, args ...interface{}) {
	fmt.Fprintf(c.errLog, f, args...)
}

func (c *context) logln(args ...interface{}) {
	fmt.Fprintln(c.errLog, args...)
}
