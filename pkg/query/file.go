package query

import (
	"go/ast"
)

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o file_dump.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o file_map.go

type FileMap map[string]*File // +tag map:""

// +tag dump:""
type File struct {
	*ast.File
}

func FromFile(f *ast.File) *File {
	return &File{f}
}

func (f *File) Tags() Tags {
	return extractTags(f.File.Doc)
}

func (f *File) GenDeclIter() GenDeclIter {
	c := make(chan *GenDecl)

	go func() {
		defer close(c)

		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok {
				c <- &GenDecl{decl}
			}
		}
	}()

	return c
}

func (f *File) TypeIter() TypeDeclIter {
	c := make(chan *TypeDecl)

	go func() {
		defer close(c)

		for decl := range f.GenDeclIter() {
			if decl.IsType() {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.TypeSpec); ok {
						c <- &TypeDecl{f, decl, &TypeSpec{spec}}
					}
				}
			}
		}
	}()

	return c
}

func (f *File) TypeDecl(name string) *TypeDecl {
	return f.TypeIter().Find(func(ty *TypeDecl) bool {
		return ty.Name() == name
	})
}

func (f *File) TypeDecls() TypeDeclMap {
	items := make(TypeDeclMap)

	for ty := range f.TypeIter() {
		items[ty.Name()] = ty
	}

	return items
}

func (f *File) InterfaceIter() InterfaceIter {
	c := make(chan *InterfaceDef)

	go func() {
		defer close(c)

		for ty := range f.TypeIter() {
			if ty.IsInterface() {
				c <- &InterfaceDef{ty, ty.AsInterface()}
			}
		}
	}()

	return c
}

func (f *File) Interface(name string) *InterfaceDef {
	return f.InterfaceIter().Find(func(intf *InterfaceDef) bool {
		return intf.Name() == name
	})
}

func (f *File) Interfaces() InterfaceMap {
	items := make(InterfaceMap)

	for intf := range f.InterfaceIter() {
		items[intf.Name()] = intf
	}

	return items
}

func (f *File) StructIter() StructIter {
	c := make(chan *StructDef)

	go func() {
		defer close(c)

		for ty := range f.TypeIter() {
			if ty.IsStruct() {
				c <- &StructDef{ty, ty.AsStruct()}
			}
		}
	}()

	return c
}

func (f *File) Struct(name string) *StructDef {
	return f.StructIter().Find(func(s *StructDef) bool {
		return s.Name() == name
	})
}

func (f *File) Structs() StructMap {
	items := make(StructMap)

	for s := range f.StructIter() {
		items[s.Name()] = s
	}

	return items
}

func (f *File) FuncIter() FuncDeclIter {
	c := make(chan *FuncDecl)

	go func() {
		defer close(c)

		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				c <- &FuncDecl{f, decl, &FuncType{decl.Type}}
			}
		}
	}()

	return c
}

func (f *File) Func(name string) *FuncDecl {
	return f.FuncIter().Find(func(f *FuncDecl) bool {
		return f.Name() == name
	})
}

func (f *File) Funcs() FuncDeclMap {
	items := make(FuncDeclMap)

	for fd := range f.FuncIter() {
		items[fd.Name()] = fd
	}

	return items
}

func (f *File) ImportIter() ImportDeclIter {
	c := make(chan *ImportDecl)

	go func() {
		defer close(c)

		for decl := range f.GenDeclIter() {
			if decl.IsImport() {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.ImportSpec); ok {
						c <- &ImportDecl{decl, &ImportSpec{spec}}
					}
				}
			}
		}
	}()

	return c
}

func (f *File) Import(path string) *ImportDecl {
	return f.ImportIter().Find(func(i *ImportDecl) bool {
		return i.Path() == path
	})
}

func (f *File) Imports() ImportDeclMap {
	items := make(ImportDeclMap)

	for i := range f.ImportIter() {
		items[i.Path()] = i
	}

	return items
}

func (f *File) ConstIter() ConstDeclIter {
	c := make(chan *ConstDecl)

	go func() {
		defer close(c)

		for decl := range f.GenDeclIter() {
			if decl.IsConst() {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.ValueSpec); ok {
						c <- &ConstDecl{decl, &ValueSpec{spec}}
					}
				}
			}
		}
	}()

	return c
}

func (f *File) Const(name string) *ConstDecl {
	for decl := range f.ConstIter() {
		for _, varName := range decl.Names() {
			if varName == name {
				return decl
			}
		}
	}

	return nil
}

func (f *File) Consts() ConstDeclMap {
	items := make(ConstDeclMap)

	for decl := range f.ConstIter() {
		for _, name := range decl.Names() {
			items[name] = decl
		}
	}

	return items
}

func (f *File) VarIter() VarDeclIter {
	c := make(chan *VarDecl)

	go func() {
		defer close(c)

		for decl := range f.GenDeclIter() {
			if decl.IsVar() {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.ValueSpec); ok {
						c <- &VarDecl{decl, &ValueSpec{spec}}
					}
				}
			}
		}
	}()

	return c
}

func (f *File) Var(name string) *VarDecl {
	for decl := range f.VarIter() {
		for _, varName := range decl.Names() {
			if varName == name {
				return decl
			}
		}
	}

	return nil
}

func (f *File) Vars() VarDeclMap {
	items := make(VarDeclMap)

	for decl := range f.VarIter() {
		for _, name := range decl.Names() {
			items[name] = decl
		}
	}

	return items
}
