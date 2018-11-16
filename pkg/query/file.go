package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

//go:generate astgen -t ../../template/iter.gogo -p $GOFILE -o file_iter.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o file_map.go

type FileMap map[string]*File // +map

type File struct {
	*ast.File
}

func FromFile(f *ast.File) *File {
	return &File{f}
}

func (f *File) Dump() string {
	return astDump(f.File)
}

func (f *File) Tags() Tags {
	return extractTags(f.File.Doc)
}

type GenDeclIter <-chan *GenDecl // +iter

type GenDecl struct {
	*ast.GenDecl
}

func (d *GenDecl) IsImport() bool { return d.GenDecl.Tok == token.IMPORT }
func (d *GenDecl) IsConst() bool  { return d.GenDecl.Tok == token.CONST }
func (d *GenDecl) IsType() bool   { return d.GenDecl.Tok == token.TYPE }
func (d *GenDecl) IsVar() bool    { return d.GenDecl.Tok == token.VAR }

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

type TypeDeclIter <-chan *TypeDecl    // +iter
type TypeDeclMap map[string]*TypeDecl // +map

type TypeDecl struct {
	*File
	*GenDecl
	*TypeSpec
}

func (t *TypeDecl) Tags() Tags {
	var docs []*ast.CommentGroup

	if t.File != nil {
		docs = append(docs, t.File.Doc)
	}
	if t.GenDecl != nil {
		docs = append(docs, t.GenDecl.Doc)
	}
	if t.TypeSpec != nil {
		docs = append(docs, t.TypeSpec.TypeSpec.Doc, t.TypeSpec.TypeSpec.Comment)
	}

	return extractTags(docs...)
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

type InterfaceIter <-chan *InterfaceDef    // +iter
type InterfaceMap map[string]*InterfaceDef // +map

type InterfaceDef struct {
	*TypeSpec
	*InterfaceType
}

func (intf *InterfaceDef) String() string {
	return fmt.Sprintf("type %s %s", intf.Name(), intf.InterfaceType)
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

type StructIter <-chan *StructDef    // +iter
type StructMap map[string]*StructDef // +map

type StructDef struct {
	*TypeSpec
	*StructType
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

type FuncDeclIter <-chan *FuncDecl    // +iter
type FuncDeclMap map[string]*FuncDecl // +map

type FuncDecl struct {
	*ast.FuncDecl
	*FuncType
}

func (f *FuncDecl) Name() string {
	return f.FuncDecl.Name.Name
}

func (f *FuncDecl) IsFunc() bool {
	return f.FuncDecl.Recv == nil
}

func (f *FuncDecl) IsMethod() bool {
	return f.FuncDecl.Recv != nil
}

func (f *FuncDecl) Recv() *NamedField {
	recv := f.FuncDecl.Recv

	if recv != nil && len(recv.List) > 0 {
		field := recv.List[0]

		var ident *ast.Ident
		if len(field.Names) > 0 {
			ident = field.Names[0]
		}

		return &NamedField{&Field{field}, ident}
	}

	return nil
}

func (f *FuncDecl) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("func ")

	if recv := f.Recv(); recv != nil {
		buf.WriteString(fmt.Sprintf("(%s) ", recv))
	}

	buf.WriteString(f.Name())
	buf.WriteString(f.FuncType.String())

	return buf.String()
}

func (f *FuncDecl) Dump() string {
	return astDump(f.FuncDecl)
}

func (f *File) FuncIter() FuncDeclIter {
	c := make(chan *FuncDecl)

	go func() {
		defer close(c)

		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				c <- &FuncDecl{decl, &FuncType{decl.Type}}
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

type ImportDeclIter <-chan *ImportDecl    // +iter
type ImportDeclMap map[string]*ImportDecl // +map

type ImportDecl struct {
	*File
	*GenDecl
	*ImportSpec
}

func (f *File) ImportIter() ImportDeclIter {
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

type ConstDeclIter <-chan *ConstDecl    // +iter
type ConstDeclMap map[string]*ConstDecl // +map

type ConstDecl struct {
	*File
	*GenDecl
	*ValueSpec
}

func (c *ConstDecl) String() string {
	return "const " + c.ValueSpec.String()
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

type VarDeclIter <-chan *VarDecl    // +iter
type VarDeclMap map[string]*VarDecl // +map

type VarDecl struct {
	*File
	*GenDecl
	*ValueSpec
}

func (v *VarDecl) String() string {
	return "var " + v.ValueSpec.String()
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
