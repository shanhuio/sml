package smake

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
)

func listFiles(pkg *build.Package) []string {
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	return files
}

func listAbsFiles(pkg *build.Package) []string {
	files := listFiles(pkg)
	for i, f := range files {
		files[i] = filepath.Join(pkg.Dir, f)
	}
	return files
}

func fileSourceMap(pkg *relPkg) (map[string][]byte, error) {
	files := listFiles(pkg.pkg)
	fileMap := make(map[string][]byte)

	for _, f := range files {
		path := filepath.Join(pkg.pkg.Dir, f)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %s", path, err)
		}
		fileMap[path] = src
	}

	return fileMap, nil
}
