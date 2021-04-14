package smake

import (
	"path/filepath"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/gomod"
	"shanhu.io/misc/osutil"
)

type goModule struct {
	dir  string
	name string
}

func findGoModule(dir string) (*goModule, error) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return nil, errcode.Annotatef(err, "absolute path of %q", dir)
	}

	for {
		modFile := filepath.Join(d, "go.mod")
		ok, err := osutil.IsRegular(modFile)
		if err != nil {
			return nil, errcode.Annotate(err, "check go.mod file")
		}
		if !ok {
			if d == "/" {
				break
			}
			d = filepath.Dir(d)
			if d == "" {
				break
			}
		}

		f, err := gomod.Parse(modFile)
		if err != nil {
			return nil, errcode.Annotatef(err, "pasre go.mod: %q", modFile)
		}

		return &goModule{dir: d, name: f.Name}, nil
	}

	return nil, errcode.NotFoundf("go module not found for dir %q", dir)
}
