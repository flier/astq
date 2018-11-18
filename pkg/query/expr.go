package query

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

//go:generate astgen -t ../../template/dump.gogo -p $GOFILE -o expr_dump.go

type Expr interface {
	fmt.Stringer
}

func FromExpr(expr ast.Expr) Expr {
	return asExpr(expr)
}

func asExpr(e ast.Expr) Expr {
	if e == nil {
		return nil
	}

	switch expr := e.(type) {
	case *ast.BadExpr:
		return &BadExpr{&AstExpr{e}, expr}
	case *ast.Ident:
		return &Ident{&AstExpr{e}, expr}
	case *ast.Ellipsis:
		return &Ellipsis{&AstExpr{e}, expr}
	case *ast.BasicLit:
		return &BasicLit{&AstExpr{e}, expr}
	case *ast.FuncLit:
		return &FuncLit{&AstExpr{e}, expr}
	case *ast.CompositeLit:
		return &CompositeLit{&AstExpr{e}, expr}
	case *ast.ParenExpr:
		return &ParenExpr{&AstExpr{e}, expr}
	case *ast.SelectorExpr:
		return &SelectorExpr{&AstExpr{e}, expr}
	case *ast.IndexExpr:
		return &IndexExpr{&AstExpr{e}, expr}
	case *ast.SliceExpr:
		return &SliceExpr{&AstExpr{e}, expr}
	case *ast.TypeAssertExpr:
		return &TypeAssertExpr{&AstExpr{e}, expr}
	case *ast.CallExpr:
		return &CallExpr{&AstExpr{e}, expr}
	case *ast.StarExpr:
		return &StarExpr{&AstExpr{e}, expr}
	case *ast.UnaryExpr:
		return &UnaryExpr{&AstExpr{e}, expr}
	case *ast.BinaryExpr:
		return &BinaryExpr{&AstExpr{e}, expr}
	case *ast.KeyValueExpr:
		return &KeyValueExpr{&AstExpr{e}, expr}
	case *ast.ArrayType:
		return &ArrayExpr{&AstExpr{e}, &ArrayType{expr}}
	case *ast.StructType:
		return &StructExpr{&AstExpr{e}, &StructType{expr}}
	case *ast.FuncType:
		return &FuncExpr{&AstExpr{e}, &FuncType{expr}}
	case *ast.InterfaceType:
		return &InterfaceExpr{&AstExpr{e}, &InterfaceType{expr}}
	case *ast.MapType:
		return &MapExpr{&AstExpr{e}, &MapType{expr}}
	case *ast.ChanType:
		return &ChanExpr{&AstExpr{e}, &ChanType{expr}}
	default:
		panic(e)
	}
}

// +tag dump:""
type AstExpr struct {
	ast.Expr
}

func (e *AstExpr) Kind() reflect.Kind {
	switch expr := e.Expr.(type) {
	case *ast.Ident:
		switch expr.Name {
		case "bool":
			return reflect.Bool
		case "int":
			return reflect.Int
		case "int8":
			return reflect.Int8
		case "int16":
			return reflect.Int16
		case "int32":
			return reflect.Int32
		case "int64":
			return reflect.Int64
		case "uint":
			return reflect.Uint
		case "uint8":
			return reflect.Uint8
		case "uint16":
			return reflect.Uint16
		case "uint32":
			return reflect.Uint32
		case "uint64":
			return reflect.Uint64
		case "uintptr":
			return reflect.Uintptr
		case "float32":
			return reflect.Float32
		case "float64":
			return reflect.Float64
		case "complex64":
			return reflect.Complex64
		case "complex128":
			return reflect.Complex128
		case "string":
			return reflect.String
		default:
			return reflect.Invalid
		}
	case *ast.ArrayType:
		return reflect.Array
	case *ast.ChanType:
		return reflect.Chan
	case *ast.FuncType:
		return reflect.Func
	case *ast.InterfaceType:
		return reflect.Interface
	case *ast.MapType:
		return reflect.Map
	case *ast.StarExpr:
		return reflect.Ptr
	case *ast.SliceExpr:
		return reflect.Slice
	case *ast.StructType:
		return reflect.Struct
	case *ast.SelectorExpr:
		if x, ok := expr.X.(*ast.Ident); ok && x.Name == "unsafe" && expr.Sel.Name == "Pointer" {
			return reflect.UnsafePointer
		}
	}

	return reflect.Invalid
}

func (e *AstExpr) IsBool() bool          { return e.Kind() == reflect.Bool }
func (e *AstExpr) IsInt() bool           { return e.Kind() == reflect.Int }
func (e *AstExpr) IsInt8() bool          { return e.Kind() == reflect.Int8 }
func (e *AstExpr) IsInt16() bool         { return e.Kind() == reflect.Int16 }
func (e *AstExpr) IsInt32() bool         { return e.Kind() == reflect.Int32 }
func (e *AstExpr) IsInt64() bool         { return e.Kind() == reflect.Int64 }
func (e *AstExpr) IsUint() bool          { return e.Kind() == reflect.Uint }
func (e *AstExpr) IsUint8() bool         { return e.Kind() == reflect.Uint8 }
func (e *AstExpr) IsUint16() bool        { return e.Kind() == reflect.Uint16 }
func (e *AstExpr) IsUint32() bool        { return e.Kind() == reflect.Uint32 }
func (e *AstExpr) IsUint64() bool        { return e.Kind() == reflect.Uint64 }
func (e *AstExpr) IsUintptr() bool       { return e.Kind() == reflect.Uintptr }
func (e *AstExpr) IsFloat32() bool       { return e.Kind() == reflect.Float32 }
func (e *AstExpr) IsFloat64() bool       { return e.Kind() == reflect.Float64 }
func (e *AstExpr) IsComplex64() bool     { return e.Kind() == reflect.Complex64 }
func (e *AstExpr) IsComplex128() bool    { return e.Kind() == reflect.Complex128 }
func (e *AstExpr) IsArray() bool         { return e.Kind() == reflect.Array }
func (e *AstExpr) IsChan() bool          { return e.Kind() == reflect.Chan }
func (e *AstExpr) IsFunc() bool          { return e.Kind() == reflect.Func }
func (e *AstExpr) IsInterface() bool     { return e.Kind() == reflect.Interface }
func (e *AstExpr) IsMap() bool           { return e.Kind() == reflect.Map }
func (e *AstExpr) IsPtr() bool           { return e.Kind() == reflect.Ptr }
func (e *AstExpr) IsSlice() bool         { return e.Kind() == reflect.Slice }
func (e *AstExpr) IsString() bool        { return e.Kind() == reflect.String }
func (e *AstExpr) IsStruct() bool        { return e.Kind() == reflect.Struct }
func (e *AstExpr) IsUnsafePointer() bool { return e.Kind() == reflect.UnsafePointer }

func (e *AstExpr) IsBadExpr() bool {
	_, ok := e.Expr.(*ast.BadExpr)

	return ok
}

func (e *AstExpr) IsIdent() bool {
	_, ok := e.Expr.(*ast.Ident)

	return ok
}

func (e *AstExpr) IsEllipsis() bool {
	_, ok := e.Expr.(*ast.Ellipsis)

	return ok
}

func (e *AstExpr) IsBasicLit() bool {
	_, ok := e.Expr.(*ast.BasicLit)

	return ok
}

func (e *AstExpr) IsFuncLit() bool {
	_, ok := e.Expr.(*ast.FuncLit)

	return ok
}

func (e *AstExpr) IsCompositeLit() bool {
	_, ok := e.Expr.(*ast.CompositeLit)

	return ok
}

func (e *AstExpr) IsParenExpr() bool {
	_, ok := e.Expr.(*ast.ParenExpr)

	return ok
}

func (e *AstExpr) IsSelectorExpr() bool {
	_, ok := e.Expr.(*ast.SelectorExpr)

	return ok
}

func (e *AstExpr) IsIndexExpr() bool {
	_, ok := e.Expr.(*ast.IndexExpr)

	return ok
}

func (e *AstExpr) IsSliceExpr() bool {
	_, ok := e.Expr.(*ast.SliceExpr)

	return ok
}

func (e *AstExpr) IsTypeAssertExpr() bool {
	_, ok := e.Expr.(*ast.TypeAssertExpr)

	return ok
}

func (e *AstExpr) IsCallExpr() bool {
	_, ok := e.Expr.(*ast.CallExpr)

	return ok
}

func (e *AstExpr) IsStarExpr() bool {
	_, ok := e.Expr.(*ast.StarExpr)

	return ok
}

func (e *AstExpr) IsPointer() bool {
	_, ok := e.Expr.(*ast.StarExpr)

	return ok
}

func (e *AstExpr) IsUnaryExpr() bool {
	_, ok := e.Expr.(*ast.UnaryExpr)

	return ok
}

func (e *AstExpr) IsBinaryExpr() bool {
	_, ok := e.Expr.(*ast.BinaryExpr)

	return ok
}

func (e *AstExpr) IsKeyValueExpr() bool {
	_, ok := e.Expr.(*ast.KeyValueExpr)

	return ok
}

type BadExpr struct {
	*AstExpr
	*ast.BadExpr
}

func (e *BadExpr) String() string {
	return fmt.Sprintf("BAD[%d:%d]", e.From, e.To)
}

type Ident struct {
	*AstExpr
	*ast.Ident
}

func (i *Ident) Name() string {
	return i.Ident.Name
}

func (i *Ident) HasObject() bool {
	return i.Ident.Obj != nil
}

func (i *Ident) Object() *Object {
	return &Object{i.Ident.Obj}
}

func (i *Ident) Spec() Spec {
	if obj := i.Object(); obj != nil {
		return obj.Spec()
	}

	return nil
}

type Object struct {
	*ast.Object
}

func (o *Object) Dump() string {
	return astDump(o.Object)
}

func (obj *Object) IsPackage() bool  { return obj.Kind == ast.Pkg }
func (obj *Object) IsConstant() bool { return obj.Kind == ast.Con }
func (obj *Object) IsType() bool     { return obj.Kind == ast.Typ }
func (obj *Object) IsVar() bool      { return obj.Kind == ast.Var }
func (obj *Object) IsFunc() bool     { return obj.Kind == ast.Fun }
func (obj *Object) IsLabel() bool    { return obj.Kind == ast.Lbl }

type Spec interface{}

func (obj *Object) Spec() Spec {
	switch obj.Kind {
	case ast.Pkg: // package
		if spec, ok := obj.Decl.(*ast.ImportSpec); ok {
			return &ImportSpec{spec}
		}
	case ast.Con: // constant
	case ast.Typ: // type
		if spec, ok := obj.Decl.(*ast.TypeSpec); ok {
			return &TypeSpec{spec}
		}
	case ast.Var: // variable
	case ast.Fun: // function or method
		if decl, ok := obj.Decl.(*ast.FuncDecl); ok {
			return &FuncDecl{nil, decl, &FuncType{decl.Type}}
		}
	case ast.Lbl: // label
		if stmt, ok := obj.Decl.(*ast.LabeledStmt); ok {
			return &Labeled{stmt}
		}
	}

	return nil
}

type Ellipsis struct {
	*AstExpr
	*ast.Ellipsis
}

func (e *Ellipsis) Elem() Expr {
	return asExpr(e.Ellipsis.Elt)
}

func (e *Ellipsis) String() string {
	if elem := e.Elem(); elem != nil {
		return fmt.Sprintf("...%s", elem)
	}

	return "..."
}

type BasicLit struct {
	*AstExpr
	*ast.BasicLit
}

func (lit *BasicLit) IsInt() bool {
	return lit.BasicLit.Kind == token.INT
}

func (lit *BasicLit) IsFloat() bool {
	return lit.BasicLit.Kind == token.FLOAT
}

func (lit *BasicLit) IsImag() bool {
	return lit.BasicLit.Kind == token.IMAG
}

func (lit *BasicLit) IsChar() bool {
	return lit.BasicLit.Kind == token.CHAR
}

func (lit *BasicLit) IsString() bool {
	return lit.BasicLit.Kind == token.STRING
}

func (lit *BasicLit) Int() int64 {
	n, err := strconv.ParseInt(lit.BasicLit.Value, 10, 64)

	if err != nil {
		panic(err)
	}

	return n
}

func (lit *BasicLit) Float() float64 {
	n, err := strconv.ParseFloat(lit.BasicLit.Value, 64)

	if err != nil {
		panic(err)
	}

	return n
}

func (lit *BasicLit) Char() rune {
	c, _ := utf8.DecodeRuneInString(lit.BasicLit.Value)

	if c == utf8.RuneError {
		panic(c)
	}

	return c
}

func (lit *BasicLit) Value() string {
	return lit.BasicLit.Value
}

func (lit *BasicLit) String() string {
	return lit.Value()
}

type FuncLit struct {
	*AstExpr
	*ast.FuncLit
}

func (lit *FuncLit) Type() *FuncType {
	return &FuncType{lit.FuncLit.Type}
}

func (lit *FuncLit) String() string {
	return fmt.Sprintf("%s {}", lit.Type())
}

type CompositeLit struct {
	*AstExpr
	*ast.CompositeLit
}

func (lit *CompositeLit) Elems() (elems []Expr) {
	for _, elem := range lit.CompositeLit.Elts {
		elems = append(elems, asExpr(elem))
	}

	return
}

func (lit *CompositeLit) String() string {
	var elems []string

	for _, elem := range lit.CompositeLit.Elts {
		elems = append(elems, (asExpr(elem)).String())
	}

	return fmt.Sprintf("{ %s }", strings.Join(elems, ", "))
}

type ParenExpr struct {
	*AstExpr
	*ast.ParenExpr
}

func (e *ParenExpr) Elem() Expr {
	return asExpr(e.ParenExpr.X)
}

func (e *ParenExpr) String() string {
	return fmt.Sprintf("( %s )", e.Elem())
}

type SelectorExpr struct {
	*AstExpr
	*ast.SelectorExpr
}

func (e *SelectorExpr) Target() Expr {
	return asExpr(e.SelectorExpr.X)
}

func (e *SelectorExpr) Selector() *Ident {
	return &Ident{&AstExpr{e.SelectorExpr.Sel}, e.SelectorExpr.Sel}
}

func (e *SelectorExpr) String() string {
	return fmt.Sprintf("%s.%s", e.Target(), e.Selector())
}

type IndexExpr struct {
	*AstExpr
	*ast.IndexExpr
}

func (e *IndexExpr) Target() Expr {
	return asExpr(e.IndexExpr.X)
}

func (e *IndexExpr) Index() Expr {
	return asExpr(e.IndexExpr.Index)
}

func (e *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", e.Target(), e.Index())
}

type SliceExpr struct {
	*AstExpr
	*ast.SliceExpr
}

func (e *SliceExpr) Target() Expr {
	return asExpr(e.SliceExpr.X)
}

func (e *SliceExpr) Low() Expr {
	return asExpr(e.SliceExpr.Low)
}

func (e *SliceExpr) High() Expr {
	return asExpr(e.SliceExpr.High)
}

func (e *SliceExpr) Max() Expr {
	return asExpr(e.SliceExpr.Max)
}

func (e *SliceExpr) String() string {
	var low, high, max string

	if n := e.Low(); n != nil {
		low = n.String()
	}
	if n := e.High(); n != nil {
		high = n.String()
	}
	if n := e.Max(); n != nil {
		max = n.String()
	}

	if e.SliceExpr.Slice3 {
		return fmt.Sprintf("%s[%s:%s:%s]", e.Target(), low, high, max)
	}

	return fmt.Sprintf("%s[%s:%s]", e.Target(), low, high)
}

type TypeAssertExpr struct {
	*AstExpr
	*ast.TypeAssertExpr
}

func (e *TypeAssertExpr) Target() Expr {
	return asExpr(e.TypeAssertExpr.X)
}

func (e *TypeAssertExpr) Type() Expr {
	return asExpr(e.TypeAssertExpr.Type)
}

func (e *TypeAssertExpr) String() string {
	ty := e.Type()

	if ty == nil {
		return fmt.Sprintf("%s.(type)", e.Target())
	}

	return fmt.Sprintf("%s.(%s)", e.Target(), ty)
}

type CallExpr struct {
	*AstExpr
	*ast.CallExpr
}

func (e *CallExpr) Func() Expr {
	return asExpr(e.CallExpr.Fun)
}

func (e *CallExpr) String() string {
	return fmt.Sprintf("%s ()", e.Func())
}

type StarExpr struct {
	*AstExpr
	*ast.StarExpr
}

func (e *StarExpr) Target() Expr {
	return asExpr(e.StarExpr.X)
}

func (e *StarExpr) String() string {
	return fmt.Sprintf("*%s", e.Target())
}

type UnaryExpr struct {
	*AstExpr
	*ast.UnaryExpr
}

func (e *UnaryExpr) Op() token.Token {
	return e.UnaryExpr.Op
}

func (e *UnaryExpr) Elem() Expr {
	return asExpr(e.UnaryExpr.X)
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("%s %s", e.Op(), e.Elem())
}

type BinaryExpr struct {
	*AstExpr
	*ast.BinaryExpr
}

func (e *BinaryExpr) Op() token.Token {
	return e.BinaryExpr.Op
}

func (e *BinaryExpr) Left() Expr {
	return asExpr(e.BinaryExpr.X)
}

func (e *BinaryExpr) Right() Expr {
	return asExpr(e.BinaryExpr.Y)
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left(), e.Op(), e.Right())
}

type KeyValueExpr struct {
	*AstExpr
	*ast.KeyValueExpr
}

func (e *KeyValueExpr) Key() Expr {
	return asExpr(e.KeyValueExpr.Key)
}

func (e *KeyValueExpr) Value() Expr {
	return asExpr(e.KeyValueExpr.Value)
}

func (e *KeyValueExpr) String() string {
	return fmt.Sprintf("%s:%s", e.Key(), e.Value())
}

type ArrayExpr struct {
	*AstExpr
	*ArrayType
}

type StructExpr struct {
	*AstExpr
	*StructType
}

type FuncExpr struct {
	*AstExpr
	*FuncType
}

type MapExpr struct {
	*AstExpr
	*MapType
}

type ChanExpr struct {
	*AstExpr
	*ChanType
}

type InterfaceExpr struct {
	*AstExpr
	*InterfaceType
}

type Path struct {
	ast.Expr
}

func (p *Path) String() string {
	if p.Tail() == nil {
		return p.Head()
	}

	return fmt.Sprintf("%s.%s", p.Head(), p.Tail())
}

func (p *Path) Head() string {
	switch expr := p.Expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return "*"
	case *ast.SelectorExpr:
		p := &Path{expr.X}
		return p.Head()
	default:
		panic(fmt.Errorf("unexpect expr [%d:%d]: %v", expr.Pos(), expr.End(), expr))
	}
}

func (p *Path) Tail() *Path {
	switch expr := p.Expr.(type) {
	case *ast.Ident:
		return nil
	case *ast.StarExpr:
		return &Path{expr.X}
	case *ast.SelectorExpr:
		return &Path{expr.Sel}
	default:
		panic(fmt.Errorf("unexpect expr [%d:%d]: %v", expr.Pos(), expr.End(), expr))
	}
}

func (p *Path) Last() string {
	switch expr := p.Expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		p := &Path{expr.X}
		return p.Last()
	case *ast.SelectorExpr:
		p := &Path{expr.Sel}
		return p.Last()
	default:
		panic(fmt.Errorf("unexpect expr [%d:%d]: %v", expr.Pos(), expr.End(), expr))
	}
}
