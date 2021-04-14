package smake

import (
	"path/filepath"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/osutil"
)

func findGoModuleRoot(d string) (string, error) {
	for {
		modFile := filepath.Join(d, "go.mod")
		ok, err := osutil.IsRegular(modFile)
		if err != nil {
			return "", errcode.Annotate(err, "check go.mod file")
		}
		if !ok {
			if d == "/" {
				break
			}
			d = filepath.Dir(d)
			if d == "" {
				break
			}
			continue
		}
		return d, nil
	}

	return "", errcode.NotFoundf("go module not found for dir %q", d)
}
