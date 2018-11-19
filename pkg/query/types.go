package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"path/filepath"
	"reflect"
	"strings"
)

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o types_dump.go
//go:generate astgen -t ../../template/iter.gogo -p $GOFILE -o types_iter.go
//go:generate astgen -t ../../template/map.gogo -p $GOFILE -o types_map.go
//go:generate astgen -t ../../template/tag.gogo -p $GOFILE -o types_tag.go

type Named interface {
	Name() string
}

// +tag dump:""
type TypeSpec struct {
	*ast.TypeSpec
}

func (t *TypeSpec) Name() string {
	return t.TypeSpec.Name.Name
}

func (t *TypeSpec) Type() Expr {
	return asExpr(t.TypeSpec.Type)
}

func (t *TypeSpec) String() string {
	return fmt.Sprintf("%s %s", t.Name(), t.Type())
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

// +tag dump:""
type ArrayType struct {
	*ast.ArrayType
}

func (a *ArrayType) Len() Expr {
	return asExpr(a.ArrayType.Len)
}

func (a *ArrayType) Elem() Expr {
	return asExpr(a.ArrayType.Elt)
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[]%s", a.Elem())
}

// +tag dump:""
type MapType struct {
	*ast.MapType
}

func (m *MapType) Key() Expr {
	return asExpr(m.MapType.Key)
}

func (m *MapType) Value() Expr {
	return asExpr(m.MapType.Value)
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.Key(), m.Value())
}

// +tag dump:""
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
	return asExpr(c.ChanType.Value)
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

// +tag dump:""
type InterfaceType struct {
	*ast.InterfaceType
}

func (intf *InterfaceType) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("interface {\n")

	for method := range intf.MethodIter() {
		buf.WriteString("\t" + method.String() + "\n")
	}

	buf.WriteString("}")

	return buf.String()
}

func (intf *InterfaceType) MethodIter() MethodIter {
	c := make(chan *Method)

	go func() {
		defer close(c)

		for _, field := range intf.InterfaceType.Methods.List {
			if ty, ok := field.Type.(*ast.FuncType); ok {
				for _, ident := range field.Names {
					c <- &Method{&FuncType{ty}, ident}
				}
			}
		}
	}()

	return c
}

func (intf *InterfaceType) Method(name string) *Method {
	for method := range intf.MethodIter() {
		if method.Name() == name {
			return method
		}
	}

	return nil
}

func (intf *InterfaceType) Methods() MethodMap {
	methods := make(MethodMap)

	for method := range intf.MethodIter() {
		methods[method.Name()] = method
	}

	return methods
}

type MethodIter <-chan *Method    // +tag iter:""
type MethodMap map[string]*Method // +tag map:""

// +tag dump:""
type Method struct {
	*FuncType
	*ast.Ident
}

func (m *Method) Name() string {
	return m.Ident.Name
}

func (m *Method) String() string {
	return m.Name() + m.Signature().String()
}

// +tag dump:""
type StructType struct {
	*ast.StructType
}

func (s *StructType) String() string {
	var fields []string

	for _, field := range s.Fields() {
		fields = append(fields, "\t"+field.String())
	}

	return fmt.Sprintf("struct {\n%s\n}", strings.Join(fields, "\n"))
}

func (s *StructType) Fields() FieldList {
	return asFieldList(s.StructType.Fields)
}

func (s *StructType) HasField(name string) bool {
	return s.NamedField(name) != nil
}

func (s *StructType) NamedField(name string) *NamedField {
	for _, field := range s.Fields() {
		if len(field.Names) > 0 {
			for _, ident := range field.Names {
				if ident.Name == name {
					return &NamedField{field, ident}
				}
			}
		} else {
			f := &NamedField{field, nil}

			if f.Name() == name {
				return f
			}
		}
	}

	return nil
}

func (s *StructType) NamedFields() NamedFieldMap {
	return asNamedFieldMap(s.StructType.Fields)
}

type FieldList []*Field // +list

func (fields FieldList) String() string {
	var strs []string

	for _, field := range fields {
		strs = append(strs, field.String())
	}

	return strings.Join(strs, ", ")
}

func asFieldList(fields *ast.FieldList) (items FieldList) {
	if fields != nil && fields.List != nil {
		for _, field := range fields.List {
			items = append(items, &Field{field})
		}
	}

	return
}

type NamedFieldMap map[string]*NamedField // +tag map:"" tag:""

func asNamedFieldMap(fields *ast.FieldList) NamedFieldMap {
	items := make(NamedFieldMap)

	if fields != nil && fields.List != nil {
		for _, field := range fields.List {
			if len(field.Names) > 0 {
				for _, ident := range field.Names {
					items[ident.Name] = &NamedField{&Field{field}, ident}
				}
			} else {
				f := &NamedField{&Field{field}, nil}
				items[f.Name()] = f
			}
		}
	}

	return items
}

// +tag dump:""
type Field struct {
	*ast.Field
}

func (f *Field) Path() *Path {
	return &Path{f.Field.Type}
}

func (f *Field) Type() Expr {
	return asExpr(f.Field.Type)
}

func (f *Field) Tags() Tags {
	tags := extractTags(f.Field.Doc, f.Field.Comment)

	if f.Field.Tag != nil {
		tags = append(tags, reflect.StructTag(strings.Trim(f.Field.Tag.Value, "`")))
	}

	return tags
}

func (f *Field) String() string {
	ty := f.Type()

	if f.Names == nil {
		return ty.String()
	}

	var names []string

	for _, ident := range f.Names {
		names = append(names, ident.Name)
	}

	return fmt.Sprintf("%s %s", strings.Join(names, ", "), ty)
}

// +tag dump:""
type NamedField struct {
	*Field
	*ast.Ident
}

func (f *NamedField) Name() string {
	if f.Ident != nil {
		return f.Ident.Name
	}

	return f.Path().Last()
}

func (f *NamedField) String() string {
	if f.Ident != nil {
		return fmt.Sprintf("%s %s", f.Ident.Name, f.Type())
	}

	return f.Type().String()
}

// +tag dump:""
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
		return fmt.Sprintf("%s %v", i.ImportSpec.Name.Name, i.ImportSpec.Path.Value)
	}

	return fmt.Sprintf("%v", i.ImportSpec.Path.Value)
}

// +tag dump:""
type Signature struct {
	*FuncType
}

func (s *Signature) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("(%s)", s.Params().String()))

	results := s.Results()

	switch len(results) {
	case 0:
	case 1:
		buf.WriteString(" " + results.String())
	default:
		buf.WriteString(fmt.Sprintf(" (%s)", results.String()))
	}

	return buf.String()
}

// +tag dump:""
type FuncType struct {
	*ast.FuncType
}

func (f *FuncType) Params() FieldList {
	return asFieldList(f.FuncType.Params)
}

func (f *FuncType) NamedParams() NamedFieldMap {
	return asNamedFieldMap(f.FuncType.Params)
}

func (f *FuncType) Results() FieldList {
	return asFieldList(f.FuncType.Results)
}

func (f *FuncType) NamedResults() NamedFieldMap {
	return asNamedFieldMap(f.FuncType.Results)
}

func (f *FuncType) Signature() *Signature {
	return &Signature{f}
}

func (f *FuncType) String() string {
	return "func " + f.Signature().String()
}

type ValueSpecMap map[string]*ValueSpec // +tag map:""

// +tag dump:""
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
	return asExpr(v.ValueSpec.Type)
}

func (v *ValueSpec) Values() (values []Expr) {
	for _, value := range v.ValueSpec.Values {
		values = append(values, asExpr(value))
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
		fmt.Fprintf(buf, " = %s", asExpr(v.ValueSpec.Values[0]))
	default:
		var values []string

		for _, value := range v.ValueSpec.Values {
			values = append(values, asExpr(value).String())
		}

		fmt.Fprintf(buf, " = (%s)", strings.Join(values, ", "))
	}

	return buf.String()
}

// +tag dump:""
type Labeled struct {
	*ast.LabeledStmt
}

func (l *Labeled) Stmt() Stmt {
	return asStmt(l.LabeledStmt.Stmt)
}

func (l *Labeled) String() string {
	return l.LabeledStmt.Label.Name + ":\n" + l.Stmt().String()
}
