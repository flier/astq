package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

var (
	whitespace = regexp.MustCompile("\\s")
)

//go:generate astgen -t ../../template/map.tpl -p $GOFILE -o types_map.go

type Named interface {
	Name() string
}

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

type TypeDeclMap map[string]*TypeDecl // +map

type TypeDecl struct {
	*File
	*ast.GenDecl
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

	return ExtractTags(docs...)
}

type TypeSpec struct {
	*ast.TypeSpec
}

func (t *TypeSpec) Name() string {
	return t.TypeSpec.Name.Name
}

func (t *TypeSpec) Type() Expr {
	return AsExpr(t.TypeSpec.Type)
}

func (t *TypeSpec) String() string {
	return fmt.Sprintf("type %s %s", t.Name(), t.Type())
}

func (t *TypeSpec) Dump() string {
	return AstDump(t.TypeSpec)
}

func (t *TypeSpec) Doc() (doc []string) {
	if t.TypeSpec.Doc != nil {
		for _, comment := range t.TypeSpec.Doc.List {
			doc = append(doc, comment.Text)
		}
	}

	return
}

func (t *TypeSpec) Comment() (doc []string) {
	if t.TypeSpec.Comment != nil {
		for _, comment := range t.TypeSpec.Comment.List {
			doc = append(doc, comment.Text)
		}
	}

	return
}

type Array struct {
	*ast.ArrayType
}

func (a *Array) Len() Expr {
	return AsExpr(a.ArrayType.Len)
}

func (a *Array) Elem() Expr {
	return AsExpr(a.ArrayType.Elt)
}

func (a *Array) String() string {
	return fmt.Sprintf("[]%s", a.Elem())
}

func (a *Array) Dump() string {
	return AstDump(a.ArrayType)
}

type Map struct {
	*ast.MapType
}

func (m *Map) Key() Expr {
	return AsExpr(m.MapType.Key)
}

func (m *Map) Value() Expr {
	return AsExpr(m.MapType.Value)
}

func (m *Map) String() string {
	return fmt.Sprintf("map[%s]%s", m.Key(), m.Value())
}

func (m *Map) Dump() string {
	return AstDump(m.MapType)
}

type FuncType struct {
	*ast.FuncType
}

func (f *FuncType) String() string {
	return fmt.Sprintf("func ()")
}

func (f *FuncType) Dump() string {
	return AstDump(f.FuncType)
}

type ChanType struct {
	*ast.ChanType
}

func (c *ChanType) Dir() ast.ChanDir {
	return c.ChanType.Dir
}

func (c *ChanType) Value() Expr {
	return AsExpr(c.ChanType.Value)
}

func (c *ChanType) String() string {
	switch c.ChanType.Dir {
	case ast.SEND:
		return fmt.Sprintf("chan<- %s", c.Value())
	case ast.RECV:
		return fmt.Sprintf("<-chan %s", c.Value())
	default:
		return fmt.Sprintf("chan %s", c.Value())
	}
}

func (c *ChanType) Dump() string {
	return AstDump(c.ChanType)
}

type InterfaceMap map[string]*InterfaceDef // +map

type InterfaceDef struct {
	*TypeSpec
	*InterfaceType
}

type InterfaceType struct {
	*ast.InterfaceType
}

func (intf *InterfaceType) String() string {
	return fmt.Sprintf("interface {}")
}

func (intf *InterfaceType) Dump() string {
	return AstDump(intf.InterfaceType)
}

func (intf *InterfaceType) Method(name string) *Method {
	for _, field := range intf.InterfaceType.Methods.List {
		if ty, ok := field.Type.(*ast.FuncType); ok {
			for _, method := range field.Names {
				if method.Name == name {
					return &Method{
						intf, field, method, ty,
					}
				}
			}
		}
	}

	return nil
}

func (intf *InterfaceType) Methods() MethodMap {
	items := make(MethodMap)

	for _, field := range intf.InterfaceType.Methods.List {
		if ty, ok := field.Type.(*ast.FuncType); ok {
			for _, ident := range field.Names {
				items[ident.Name] = &Method{
					intf, field, ident, ty,
				}
			}
		}
	}

	return items
}

type MethodMap map[string]*Method // +map

type Method struct {
	*InterfaceType
	*ast.Field
	*ast.Ident
	*ast.FuncType
}

func (m *Method) Name() string {
	return m.Ident.Name
}

func (m *Method) Tag() reflect.StructTag {
	if m.Field.Tag == nil {
		return ""
	}

	return reflect.StructTag(m.Field.Tag.Value)
}

type StructMap map[string]*StructDef // +map

type StructDef struct {
	*TypeSpec
	*StructType
}

type StructType struct {
	*ast.StructType
}

func (s *StructType) String() string {
	return fmt.Sprintf("struct {}")
}

func (s *StructType) Dump() string {
	return AstDump(s.StructType)
}

func (s *StructType) Field(name string) *Field {
	for _, field := range s.StructType.Fields.List {
		if len(field.Names) > 0 {
			for _, ident := range field.Names {
				if ident.Name == name {
					return &Field{s, field, ident}
				}
			}
		} else {
			f := &Field{s, field, nil}

			if f.Name() == name {
				return f
			}
		}
	}

	return nil
}

func (s *StructType) Fields() FieldMap {
	items := make(FieldMap)

	for _, field := range s.StructType.Fields.List {
		if len(field.Names) > 0 {
			for _, ident := range field.Names {
				items[ident.Name] = &Field{s, field, ident}
			}
		} else {
			f := &Field{s, field, nil}
			items[f.Name()] = f
		}
	}

	return items
}

type FieldMap map[string]*Field // +map

type Field struct {
	*StructType
	*ast.Field
	*ast.Ident
}

func (f *Field) Path() *Path {
	return &Path{f.Field.Type}
}

func (f *Field) Name() string {
	if f.Ident != nil {
		return f.Ident.Name
	}

	return f.Path().Last()
}

func (f *Field) Type() Expr {
	return AsExpr(f.Field.Type)
}

func (f *Field) Tag() reflect.StructTag {
	if f.Field.Tag == nil {
		return ""
	}

	return reflect.StructTag(strings.Trim(f.Field.Tag.Value, "`"))
}

type ImportSpec struct {
	*ast.ImportSpec
}

type FuncDecl struct {
	*ast.FuncDecl
}

type Labeled struct {
	*ast.LabeledStmt
}

type Tags map[string]string // +map

func ExtractTags(groups ...*ast.CommentGroup) Tags {
	tags := make(Tags)

	for _, group := range groups {
		if group == nil {
			continue
		}

		for _, comment := range group.List {
			if comment := strings.TrimSpace(strings.TrimLeft(comment.Text, "/")); strings.HasPrefix(comment, "+") {
				parts := whitespace.Split(comment, 2)

				if len(parts) > 1 {
					key := strings.TrimLeft(parts[0], "+")
					value := parts[1]

					tags[key] = value
				} else {
					key := strings.TrimLeft(comment, "+")

					tags[key] = ""
				}
			}

		}
	}

	return tags
}

func AstDump(x interface{}) string {
	var buf bytes.Buffer

	if err := ast.Fprint(&buf, nil, x, ast.NotNilFilter); err != nil {
		panic(err)
	}

	return buf.String()
}
