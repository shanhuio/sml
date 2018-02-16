package goenv

import (
	"path/filepath"
)

func pkgPath(subDir, p string) string {
	return filepath.Join(subDir, filepath.FromSlash(p))
}

// PkgPath returns the file path of a sub directory (e.g. src, bin, etc.)
// for a particular package.
func PkgPath(subDir, pkg string) string {
	return pkgPath(subDir, pkg)
}

// SrcDir returns the Go language source directory.
func SrcDir(pkg string) string {
	return PkgPath("src", pkg)
}

// SrcFile returns the file path of a Go language source file.
func SrcFile(f string) string {
	return pkgPath("src", f)
}
