package smake

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
)

func listFiles(imported *build.Package) ([]string, error) {
	var files []string
	files = append(files, imported.GoFiles...)
	files = append(files, imported.CgoFiles...)
	files = append(files, imported.TestGoFiles...)
	return files, nil
}

func fileSourceMap(pkg string) (map[string][]byte, error) {
	imported, err := build.Import(pkg, "", 0)
	if err != nil {
		return nil, err
	}

	files, err := listFiles(imported)
	if err != nil {
		return nil, err
	}

	fileMap := make(map[string][]byte)

	for _, f := range files {
		path := filepath.Join(imported.Dir, f)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %s", path, err)
		}
		fileMap[path] = src
	}

	return fileMap, nil
}
