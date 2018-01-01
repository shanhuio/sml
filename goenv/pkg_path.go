package goenv

import (
	"path/filepath"
)

// PkgPath returns the file path of a sub directory (e.g. src, bin, etc.)
// for a particular package.
func PkgPath(subDir, pkg string) string {
	return filepath.Join(subDir, filepath.FromSlash(pkg))
}

// SrcDir returns the Go language source directory.
func SrcDir(pkg string) string {
	return PkgPath("src", pkg)
}
