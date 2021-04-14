package smake

import (
	"fmt"
	"path/filepath"

	lintpkg "golang.org/x/lint"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/goload"
	"shanhu.io/misc/gomod"
	"shanhu.io/sml/gotags"
	"shanhu.io/tools/gocheck"
)

func smlchk(c *context, pkgs []*relPkg) error {
	c.logln("smlchk")

	const textHeight = 300
	const textWidth = 80

	for _, pkg := range pkgs {
		errs := gocheck.CheckAll(pkg.abs, textHeight, textWidth)
		if len(errs) != 0 {
			for _, err := range errs {
				c.logln(err)
			}
			return fmt.Errorf("smlchk %q failed", pkg.rel)
		}
	}
	return nil
}

func lint(c *context, pkgs []*relPkg) error {
	c.logln("lint")

	const minConfidence = 0.8
	for _, pkg := range pkgs {
		files, err := fileSourceMap(pkg)
		if err != nil {
			return err
		}

		l := new(lintpkg.Linter)
		ps, err := l.LintFiles(files)
		if err != nil {
			return err
		}

		errCount := 0
		for _, p := range ps {
			if p.Confidence < minConfidence {
				continue
			}
			c.logf("%v: %s\n", p.Position, p.Text)
			errCount++
		}

		if errCount > 0 {
			return fmt.Errorf("lint %q failed", pkg.rel)
		}
	}
	return nil
}

func tags(c *context, pkgs []*relPkg) error {
	c.logln("tags")

	var files []string
	for _, pkg := range pkgs {
		list := listAbsFiles(pkg.pkg)
		files = append(files, list...)
	}
	return gotags.Write(files, "tags")
}

func listModPkgs(c *context) ([]*relPkg, error) {
	root := c.workDir()
	modFile := filepath.Join(root, "go.mod")
	mod, err := gomod.Parse(modFile)
	if err != nil {
		return nil, errcode.Annotate(err, "parse go.mod")
	}

	scanRes, err := goload.ScanModPkgs(mod.Name, root, nil)
	if err != nil {
		return nil, errcode.Annotate(err, "scan packages")
	}
	return relPkgs(mod.Name, scanRes)
}

func listPkgs(c *context) ([]*relPkg, error) {
	if c.gomod() {
		return listModPkgs(c)
	}

	rootPkg, err := pkgFromDir(c.srcRoot(), c.dir)
	if err != nil {
		return nil, errcode.Annotate(err, "find root package")
	}

	scanRes, err := goload.ScanPkgs(rootPkg, nil)
	if err != nil {
		return nil, errcode.Annotate(err, "scan packages")
	}

	return relPkgs(rootPkg, scanRes)
}

func smake(c *context) error {
	pkgs, err := listPkgs(c)
	if err != nil {
		return errcode.Annotate(err, "list packages")
	}

	if len(pkgs) == 0 {
		c.logln("no packages found")
		return nil
	}

	installCmd := []string{"go", "install"}
	if !c.gomod() {
		installCmd = append(installCmd, "-i")
	}

	if err := c.execPkgs(pkgs, [][]string{
		{"gofmt", "-s", "-w", "-l"},
		installCmd,
	}); err != nil {
		return err
	}

	if err := smlchk(c, pkgs); err != nil {
		return err
	}
	if err := lint(c, pkgs); err != nil {
		return err
	}

	if err := c.execPkgs(pkgs, [][]string{
		{"go", "vet"},
	}); err != nil {
		return err
	}

	return tags(c, pkgs)
}
