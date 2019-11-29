package smake

import (
	"fmt"
	"go/build"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/goload"
	"shanhu.io/sml/goenv"
)

func absGOPATH() (string, error) {
	gopath, err := goenv.GOPATH()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(gopath)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func pkgFromDir(src, dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	p, err := filepath.Rel(src, abs)
	if err != nil {
		return "", err
	}
	return filepath.FromSlash(p), nil
}

type relPkg struct {
	abs string
	rel string
	pkg *build.Package
}

func relPkgs(rootPkg string, scanRes *goload.ScanResult) ([]*relPkg, error) {
	var pkgs []string
	for pkg := range scanRes.Pkgs {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	var ret []*relPkg
	prefix := rootPkg + "/"

	for _, pkg := range pkgs {
		rel := &relPkg{
			abs: pkg,
			pkg: scanRes.Pkgs[pkg].Build,
		}

		if pkg == rootPkg {
			rel.rel = "."
			ret = append(ret, rel)
			continue
		}

		if strings.HasPrefix(pkg, prefix) {
			rel.rel = "./" + strings.TrimPrefix(pkg, prefix)
			ret = append(ret, rel)
			continue
		}

		return nil, fmt.Errorf("%q is not in %q", pkg, rootPkg)
	}
	return ret, nil
}
