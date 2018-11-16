package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var (
	whitespace = regexp.MustCompile("\\s")
)

//go:generate astgen -t ../../template/iter.tpl -p $GOFILE -o types_iter.go
//go:generate astgen -t ../../template/map.tpl -p $GOFILE -o types_map.go

type Named interface {
	Name() string
}

type GenDeclIter <-chan *GenDecl // +iter

type GenDecl struct {
	*ast.GenDecl
}

func (d *GenDecl) IsImport() bool { return d.GenDecl.Tok == token.IMPORT }
func (d *GenDecl) IsConst() bool  { return d.GenDecl.Tok == token.CONST }
func (d *GenDecl) IsType() bool   { return d.GenDecl.Tok == token.TYPE }
func (d *GenDecl) IsVar() bool    { return d.GenDecl.Tok == token.VAR }

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

func (t *TypeSpec) IsInterface() bool {
	_, ok := t.TypeSpec.Type.(*ast.InterfaceType)

	return ok
}

func (t *TypeSpec) AsInterface() *InterfaceType {
	it, ok := t.TypeSpec.Type.(*ast.InterfaceType)

	if ok {
		return &InterfaceType{it}
	}

	return nil
}

func (t *TypeSpec) IsStruct() bool {
	_, ok := t.TypeSpec.Type.(*ast.StructType)

	return ok
}

func (t *TypeSpec) AsStruct() *StructType {
	st, ok := t.TypeSpec.Type.(*ast.StructType)

	if ok {
		return &StructType{st}
	}

	return nil
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

func (c *ChanType) Dir() reflect.ChanDir {
	switch c.ChanType.Dir {
	case ast.SEND:
		return reflect.SendDir
	case ast.RECV:
		return reflect.RecvDir
	case ast.SEND | ast.RECV:
		return reflect.BothDir
	default:
		panic(c.ChanType.Dir)
	}
}

func (c *ChanType) CanSend() bool {
	return (c.ChanType.Dir & ast.SEND) == ast.SEND
}

func (c *ChanType) CanRecv() bool {
	return (c.ChanType.Dir & ast.RECV) == ast.RECV
}

func (c *ChanType) Elem() Expr {
	return AsExpr(c.ChanType.Value)
}

func (c *ChanType) String() string {
	switch c.ChanType.Dir {
	case ast.SEND:
		return fmt.Sprintf("chan<- %s", c.Elem())
	case ast.RECV:
		return fmt.Sprintf("<-chan %s", c.Elem())
	default:
		return fmt.Sprintf("chan %s", c.Elem())
	}
}

func (c *ChanType) Dump() string {
	return AstDump(c.ChanType)
}

type InterfaceIter <-chan *InterfaceDef    // +iter
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

type StructIter <-chan *StructDef    // +iter
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

type ImportDeclIter <-chan *ImportDecl    // +iter
type ImportDeclMap map[string]*ImportDecl // +map

type ImportDecl struct {
	*File
	*GenDecl
	*ImportSpec
}

type ImportSpec struct {
	*ast.ImportSpec
}

func (i *ImportSpec) Name() string {
	if i.ImportSpec.Name != nil {
		return i.ImportSpec.Name.Name
	}

	return filepath.Base(i.Path())
}

func (i *ImportSpec) Path() string {
	return strings.Trim(i.ImportSpec.Path.Value, `"`)
}

func (i *ImportSpec) String() string {
	if i.ImportSpec.Name != nil {
		return fmt.Sprintf("import %s %v", i.ImportSpec.Name.Name, i.ImportSpec.Path.Value)
	}

	return fmt.Sprintf("import %v", i.ImportSpec.Path.Value)
}

type FuncDeclIter <-chan *FuncDecl    // +iter
type FuncDeclMap map[string]*FuncDecl // +map

type FuncDecl struct {
	*ast.FuncDecl
}

func (f *FuncDecl) Name() string {
	return f.FuncDecl.Name.Name
}

type ValueSpecMap map[string]*ValueSpec // +map

type ValueSpec struct {
	*ast.ValueSpec
}

func (v *ValueSpec) Names() (names []string) {
	for _, ident := range v.ValueSpec.Names {
		names = append(names, ident.Name)
	}

	return
}

func (v *ValueSpec) Type() Expr {
	return AsExpr(v.ValueSpec.Type)
}

func (v *ValueSpec) Values() (values []Expr) {
	for _, value := range v.ValueSpec.Values {
		values = append(values, AsExpr(value))
	}

	return
}

func (v *ValueSpec) String() string {
	buf := new(bytes.Buffer)

	if len(v.ValueSpec.Names) > 1 {
		io.WriteString(buf, fmt.Sprintf("(%s)", strings.Join(v.Names(), ", ")))
	} else {
		io.WriteString(buf, v.ValueSpec.Names[0].Name)
	}

	if ty := v.Type(); ty != nil {
		fmt.Fprintf(buf, " %s", ty)
	}
	switch len(v.ValueSpec.Values) {
	case 0:
	case 1:
		fmt.Fprintf(buf, " = %s", AsExpr(v.ValueSpec.Values[0]))
	default:
		var values []string

		for _, value := range v.ValueSpec.Values {
			values = append(values, AsExpr(value).String())
		}

		fmt.Fprintf(buf, " = (%s)", strings.Join(values, ", "))
	}

	return buf.String()
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
