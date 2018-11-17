package query

import "go/ast"

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o pkg_dump.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o pkg_map.go

// Packages represents a map of Go package.
type Packages map[string]*Package // +map

// Package node represents a set of source files collectively building a Go package.
// +dump
type Package struct {
	*ast.Package
}

// FromPackage returns a queriable Package
func FromPackage(p *ast.Package) *Package {
	return &Package{p}
}

func FromPackages(pkgs map[string]*ast.Package) Packages {
	wrapped := make(map[string]*Package)

	for name, pkg := range pkgs {
		wrapped[name] = &Package{pkg}
	}

	return wrapped
}

func (p *Package) File(name string) *File {
	for filename, file := range p.Package.Files {
		if filename == name {
			return &File{file}
		}
	}

	return nil
}

func (p *Package) Files() FileMap {
	files := make(FileMap)

	for name, file := range p.Package.Files {
		files[name] = &File{file}
	}

	return files
}
