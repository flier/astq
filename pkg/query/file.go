package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o file_dump.go
//go:generate astgen -t ../../template/iter.gogo -p $GOFILE -o file_iter.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o file_map.go
//go:generate astgen -t ../../template/tag.gogo -p $GOFILE -o file_tag.go

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

type GenDeclIter <-chan *GenDecl // +tag iter:"" tag:""

// +tag dump:""
type GenDecl struct {
	*ast.GenDecl
}

func (d *GenDecl) Tags() Tags {
	return extractTags(d.GenDecl.Doc)
}

func (d *GenDecl) IsImport() bool { return d.GenDecl.Tok == token.IMPORT }
func (d *GenDecl) IsConst() bool  { return d.GenDecl.Tok == token.CONST }
func (d *GenDecl) IsType() bool   { return d.GenDecl.Tok == token.TYPE }
func (d *GenDecl) IsVar() bool    { return d.GenDecl.Tok == token.VAR }

func (d *GenDecl) Imports() (decls []*ImportDecl) {
	if d.GenDecl.Tok == token.IMPORT {
		for _, spec := range d.GenDecl.Specs {
			if spec, ok := spec.(*ast.ImportSpec); ok {
				decls = append(decls, &ImportDecl{d, &ImportSpec{spec}})
			}
		}
	}

	return
}

func (d *GenDecl) Consts() (decls []*ConstDecl) {
	if d.GenDecl.Tok == token.CONST {
		for _, spec := range d.GenDecl.Specs {
			if spec, ok := spec.(*ast.ValueSpec); ok {
				decls = append(decls, &ConstDecl{d, &ValueSpec{spec}})
			}
		}
	}

	return
}

func (d *GenDecl) Types() (decls []*TypeDecl) {
	if d.GenDecl.Tok == token.TYPE {
		for _, spec := range d.GenDecl.Specs {
			if spec, ok := spec.(*ast.TypeSpec); ok {
				decls = append(decls, &TypeDecl{nil, d, &TypeSpec{spec}})
			}
		}
	}

	return
}

func (d *GenDecl) Vars() (decls []*VarDecl) {
	if d.GenDecl.Tok == token.VAR {
		for _, spec := range d.GenDecl.Specs {
			if spec, ok := spec.(*ast.ValueSpec); ok {
				decls = append(decls, &VarDecl{d, &ValueSpec{spec}})
			}
		}
	}

	return
}

func (d *GenDecl) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(d.Tok.String() + "{\n")
	for _, spec := range d.GenDecl.Specs {
		switch spec := spec.(type) {
		case *ast.ImportSpec:
			buf.WriteString((&ImportSpec{spec}).String())
		case *ast.TypeSpec:
			buf.WriteString((&TypeSpec{spec}).String())
		case *ast.ValueSpec:
			buf.WriteString((&ValueSpec{spec}).String())
		}
	}
	buf.WriteString("}")

	return buf.String()
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

type TypeDeclIter <-chan *TypeDecl    // +tag iter:"" tag:""
type TypeDeclMap map[string]*TypeDecl // +tag map:"" tag:""

// +tag dump:"TypeSpec"
type TypeDecl struct {
	*File
	*GenDecl
	*TypeSpec
}

func (t *TypeDecl) String() string {
	return "type " + t.TypeSpec.String()
}

func (t *TypeDecl) Tags() Tags {
	return extractTags(t.GenDecl.Doc, t.TypeSpec.TypeSpec.Doc, t.TypeSpec.TypeSpec.Comment)
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

type InterfaceIter <-chan *InterfaceDef    // +tag iter:"" tag:""
type InterfaceMap map[string]*InterfaceDef // +tag map:"" tag:""

// +tag dump:"InterfaceType"
type InterfaceDef struct {
	*TypeDecl
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

type StructIter <-chan *StructDef    // +tag iter:"" tag:""
type StructMap map[string]*StructDef // +tag map:"" tag:""

// +tag dump:"StructType"
type StructDef struct {
	*TypeDecl
	*StructType
}

func (s *StructDef) String() string {
	return fmt.Sprintf("type %s %s", s.Name(), s.StructType)
}

func (s *StructDef) MethodIter() MethodIter {
	c := make(chan *Method)

	go func() {
		defer close(c)

		for fn := range s.File.FuncIter() {
			recv := fn.Recv()

			if recv == nil {
				continue
			}

			switch expr := recv.Type().(type) {
			case *StarExpr:
				if expr.Target().String() == s.Name() {
					c <- &Method{fn.FuncType, fn.FuncDecl.Name}
				}
			case *Ident:
				if expr.Name() == s.Name() {
					c <- &Method{fn.FuncType, fn.FuncDecl.Name}
				}
			}
		}
	}()

	return c
}

func (s *StructDef) HasMethod(name string) bool {
	return s.Method(name) != nil
}

func (s *StructDef) Method(name string) *Method {
	for method := range s.MethodIter() {
		if method.Name() == name {
			return method
		}
	}

	return nil
}

func (s *StructDef) Methods() MethodMap {
	methods := make(MethodMap)

	for method := range s.MethodIter() {
		methods[method.Name()] = method
	}

	return methods
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

type FuncDeclIter <-chan *FuncDecl    // +tag iter:"" tag:""
type FuncDeclMap map[string]*FuncDecl // +tag map:"" tag:""

// +tag dump:"FuncDecl"
type FuncDecl struct {
	*File
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

func (f *FuncDecl) Tags() Tags {
	return extractTags(f.File.Doc, f.FuncDecl.Doc)
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
	buf.WriteString(f.FuncType.Signature().String())

	return buf.String()
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

type ImportDeclIter <-chan *ImportDecl    // +tag iter:"" tag:""
type ImportDeclMap map[string]*ImportDecl // +tag map:"" tag:""

// +tag dump:"ImportSpec"
type ImportDecl struct {
	*GenDecl
	*ImportSpec
}

func (i *ImportDecl) String() string {
	return "import " + i.ImportSpec.String()
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

type ConstDeclIter <-chan *ConstDecl    // +tag iter:"" tag:""
type ConstDeclMap map[string]*ConstDecl // +tag map:"" tag:""

// +tag dump:"ValueSpec"
type ConstDecl struct {
	*GenDecl
	*ValueSpec
}

func (c *ConstDecl) Type() Expr {
	ty := c.ValueSpec.Type()

	if ty == nil {
		for _, decl := range c.GenDecl.Specs {
			if spec, ok := decl.(*ast.ValueSpec); ok {
				if spec.Type != nil {
					ty = asExpr(spec.Type)
				}

				if c.ValueSpec.ValueSpec == spec {
					break
				}
			}
		}
	}

	return ty
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

type VarDeclIter <-chan *VarDecl    // +tag iter:"" tag:""
type VarDeclMap map[string]*VarDecl // +tag map:"" tag:""

// +tag dump:"ValueSpec"
type VarDecl struct {
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
