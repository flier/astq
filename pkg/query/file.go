package query

import (
	"go/ast"
	"go/token"
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

func (f *File) TypeDecl(name string) *TypeDecl {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok && spec.Name.Name == name {
					return &TypeDecl{f, decl, &TypeSpec{spec}}
				}
			}
		}
	}

	return nil
}

func (f *File) TypeDecls() TypeDeclMap {
	items := make(TypeDeclMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					items[spec.Name.Name] = &TypeDecl{f, decl, &TypeSpec{spec}}
				}
			}
		}
	}

	return items
}

func (f *File) Interface(name string) *InterfaceDef {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if it, ok := spec.Type.(*ast.InterfaceType); ok && spec.Name.Name == name {
						return &InterfaceDef{&TypeSpec{spec}, &InterfaceType{it}}
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Interfaces() InterfaceMap {
	items := make(InterfaceMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if it, ok := spec.Type.(*ast.InterfaceType); ok {
						items[spec.Name.Name] = &InterfaceDef{&TypeSpec{spec}, &InterfaceType{it}}
					}
				}
			}
		}
	}

	return items
}

func (f *File) Struct(name string) *StructDef {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := spec.Type.(*ast.StructType); ok && spec.Name.Name == name {
						return &StructDef{&TypeSpec{spec}, &StructType{st}}
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Structs() StructMap {
	items := make(StructMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := spec.Type.(*ast.StructType); ok {
						items[spec.Name.Name] = &StructDef{&TypeSpec{spec}, &StructType{st}}
					}
				}
			}
		}
	}

	return items
}

func (f *File) Func(name string) *FuncDecl {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok && decl.Name.Name == name {
			return &FuncDecl{decl}
		}
	}

	return nil
}

func (f *File) Funcs() FuncDeclMap {
	items := make(FuncDeclMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			items[decl.Name.Name] = &FuncDecl{decl}
		}
	}

	return items
}

func (f *File) Import(path string) *ImportDecl {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ImportSpec); ok {
					is := &ImportDecl{f, decl, &ImportSpec{spec}}

					if is.Path() == path {
						return is
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Imports() ImportDeclMap {
	items := make(ImportDeclMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ImportSpec); ok {
					is := &ImportDecl{f, decl, &ImportSpec{spec}}

					items[is.Path()] = is
				}
			}
		}
	}

	return items
}

func (f *File) Const(name string) *ConstDecl {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.CONST {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok {
					for idx, ident := range spec.Names {
						if ident.Name == name {
							return &ConstDecl{f, decl, &ValueSpec{spec, idx}}
						}
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Consts() ConstDeclMap {
	items := make(ConstDeclMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.CONST {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok {
					for idx, ident := range spec.Names {
						items[ident.Name] = &ConstDecl{f, decl, &ValueSpec{spec, idx}}
					}
				}
			}
		}
	}

	return items
}

func (f *File) Var(name string) *VarDecl {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok {
					for idx, ident := range spec.Names {
						if ident.Name == name {
							return &VarDecl{f, decl, &ValueSpec{spec, idx}}
						}
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Vars() VarDeclMap {
	items := make(VarDeclMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok {
					for idx, ident := range spec.Names {
						items[ident.Name] = &VarDecl{f, decl, &ValueSpec{spec, idx}}
					}
				}
			}
		}
	}

	return items
}
