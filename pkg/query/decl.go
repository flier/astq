package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o decl_dump.go
//go:generate astgen -t ../../template/iter.gogo -p $GOFILE -o decl_iter.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o decl_map.go
//go:generate astgen -t ../../template/tag.gogo -p $GOFILE -o decl_tag.go

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
