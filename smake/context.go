package smake

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/lint"
	"shanhu.io/misc/goload"
	"shanhu.io/tools/gocheck"
)

type context struct {
	gopath string
	dir    string
	env    []string
	errLog io.Writer
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
		errLog: os.Stderr,
	}
}

func (c *context) srcRoot() string {
	return filepath.Join(c.gopath, "src")
}

func (c *context) execPkgs(pkgs []*relPkg, tasks [][]string) error {
	for _, args := range tasks {
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

func (c *context) smlchk(pkgs []*relPkg) error {
	fmt.Fprintln(c.errLog, "smlchk")

	const textHeight = 300
	const textWidth = 80

	for _, pkg := range pkgs {
		errs := gocheck.CheckAll(pkg.abs, textHeight, textWidth)
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintln(c.errLog, err)
			}
			return fmt.Errorf("smlchk %q failed", pkg)
		}
	}
	return nil
}

func (c *context) lint(pkgs []*relPkg) error {
	fmt.Fprintln(c.errLog, "lint")

	const minConfidence = 0.8
	for _, pkg := range pkgs {
		files, err := fileSourceMap(pkg.abs)
		if err != nil {
			return err
		}

		l := new(lint.Linter)
		ps, err := l.LintFiles(files)
		if err != nil {
			return err
		}

		errCount := 0
		for _, p := range ps {
			if p.Confidence < minConfidence {
				continue
			}
			fmt.Fprintf(c.errLog, "%v: %s\n", p.Position, p.Text)
			errCount++
		}

		if errCount > 0 {
			return fmt.Errorf("lint %q failed", pkg.rel)
		}
	}
	return nil
}

func (c *context) smake() error {
	rootPkg, err := pkgFromDir(c.srcRoot(), c.dir)
	if err != nil {
		return err
	}

	absPkgs, err := goload.ListPkgs(rootPkg)
	if err != nil {
		return err
	}

	pkgs, err := relPkgs(rootPkg, absPkgs)
	if err != nil {
		return err
	}

	if err := c.execPkgs(pkgs, [][]string{
		{"gofmt", "-s", "-w", "-l"},
		{"go", "install", "-i"},
	}); err != nil {
		return err
	}

	if err := c.smlchk(pkgs); err != nil {
		return err
	}
	if err := c.lint(pkgs); err != nil {
		return err
	}

	return c.execPkgs(pkgs, [][]string{
		{"go", "vet"},
		{"gotags", "-R", "-f=tags"},
	})
}
