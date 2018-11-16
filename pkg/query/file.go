package query

import (
	"go/ast"
)

type FileMap map[string]*File // +map

type File struct {
	*ast.File
}

func FromFile(f *ast.File) *File {
	return &File{f}
}

func (f *File) Dump() string {
	return AstDump(f.File)
}

func (f *File) Tags() Tags {
	return ExtractTags(f.File.Doc)
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
	for ty := range f.TypeIter() {
		if ty.Name() == name {
			return ty
		}
	}

	return nil
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
				c <- &InterfaceDef{ty.TypeSpec, ty.AsInterface()}
			}
		}
	}()

	return c
}

func (f *File) Interface(name string) *InterfaceDef {
	for intf := range f.InterfaceIter() {
		if intf.Name() == name {
			return intf
		}
	}

	return nil
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
				c <- &StructDef{ty.TypeSpec, ty.AsStruct()}
			}
		}
	}()

	return c
}

func (f *File) Struct(name string) *StructDef {
	for s := range f.StructIter() {
		if s.Name() == name {
			return s
		}
	}

	return nil
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
				c <- &FuncDecl{decl}
			}
		}
	}()

	return c
}

func (f *File) Func(name string) *FuncDecl {
	for fd := range f.FuncIter() {
		if fd.Name() == name {
			return fd
		}
	}

	return nil
}

func (f *File) Funcs() FuncDeclMap {
	items := make(FuncDeclMap)

	for fd := range f.FuncIter() {
		items[fd.Name()] = fd
	}

	return items
}

func (f *File) ImportIter() ImportDeclChan {
	c := make(chan *ImportDecl)

	go func() {
		defer close(c)

		for decl := range f.GenDeclIter() {
			if decl.IsImport() {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.ImportSpec); ok {
						c <- &ImportDecl{f, decl, &ImportSpec{spec}}
					}
				}
			}
		}
	}()

	return c
}

func (f *File) Import(path string) *ImportDecl {
	for i := range f.ImportIter() {
		if i.Path() == path {
			return i
		}
	}

	return nil
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
						c <- &ConstDecl{f, decl, &ValueSpec{spec}}
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
						c <- &VarDecl{f, decl, &ValueSpec{spec}}
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
