package smake

import (
	"fmt"
	"strings"

	"path/filepath"
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
}

func relPkgs(rootPkg string, pkgs []string) ([]*relPkg, error) {
	var ret []*relPkg
	prefix := rootPkg + "/"
	for _, pkg := range pkgs {
		if pkg == rootPkg {
			ret = append(ret, &relPkg{abs: pkg, rel: "."})
			continue
		}

		if strings.HasPrefix(pkg, prefix) {
			rel := strings.TrimPrefix(pkg, prefix)
			ret = append(ret, &relPkg{abs: pkg, rel: "./" + rel})
			continue
		}

		return nil, fmt.Errorf("%q is not in %q", pkg, rootPkg)
	}
	return ret, nil
}
