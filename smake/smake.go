package smake

import (
	"fmt"

	lintpkg "golang.org/x/lint"
	"shanhu.io/misc/goload"
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

func smake(c *context) error {
	rootPkg, err := pkgFromDir(c.srcRoot(), c.dir)
	if err != nil {
		return err
	}

	scanRes, err := goload.ScanPkgs(rootPkg, nil)
	if err != nil {
		return err
	}

	pkgs, err := relPkgs(rootPkg, scanRes)
	if err != nil {
		return err
	}

	if err := c.execPkgs(pkgs, [][]string{
		{"gofmt", "-s", "-w", "-l"},
		{"go", "install", "-i"},
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
