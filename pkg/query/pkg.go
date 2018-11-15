package query

import "go/ast"

type Packages map[string]*Package // +map

type Package struct {
	*ast.Package
}

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

func (p *Package) Dump() string {
	return AstDump(p.Package)
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
