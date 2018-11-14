package query

import (
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

func (f *File) Tags() Tags {
	return ExtractTags(f.File.Doc)
}

func (f *File) Type(name string) *Type {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok && spec.Name.Name == name {
					return &Type{
						f, decl, spec,
					}
				}
			}
		}
	}

	return nil
}

func (f *File) Types() TypeMap {
	items := make(TypeMap)

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					items[spec.Name.Name] = &Type{
						f, decl, spec,
					}
				}
			}
		}
	}

	return items
}

func (f *File) Interface(name string) *Interface {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if it, ok := spec.Type.(*ast.InterfaceType); ok && spec.Name.Name == name {
						return &Interface{&Type{f, decl, spec}, it}
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
						items[spec.Name.Name] = &Interface{&Type{f, decl, spec}, it}
					}
				}
			}
		}
	}

	return items
}

func (f *File) Struct(name string) *Struct {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := spec.Type.(*ast.StructType); ok && spec.Name.Name == name {
						return &Struct{&Type{f, decl, spec}, st}
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
						items[spec.Name.Name] = &Struct{&Type{f, decl, spec}, st}
					}
				}
			}
		}
	}

	return items
}

type TypeMap map[string]*Type // +map

type Type struct {
	*File
	*ast.GenDecl
	*ast.TypeSpec
}

func (t *Type) Name() string {
	return t.TypeSpec.Name.Name
}

func (t *Type) Doc() (doc []string) {
	if t.TypeSpec.Doc != nil {
		for _, comment := range t.TypeSpec.Doc.List {
			doc = append(doc, comment.Text)
		}
	}

	return
}

func (t *Type) Comment() (doc []string) {
	if t.TypeSpec.Comment != nil {
		for _, comment := range t.TypeSpec.Comment.List {
			doc = append(doc, comment.Text)
		}
	}

	return
}

func (t *Type) Tags() Tags {
	return ExtractTags(t.File.Doc, t.GenDecl.Doc, t.TypeSpec.Doc, t.TypeSpec.Comment)
}

func (t *Type) IsInterface() bool {
	_, ok := t.TypeSpec.Type.(*ast.InterfaceType)

	return ok
}

func (t *Type) AsInterface() *Interface {
	i, ok := t.TypeSpec.Type.(*ast.InterfaceType)

	if ok {
		return &Interface{t, i}
	}

	return nil
}

func (t *Type) IsStruct() bool {
	_, ok := t.TypeSpec.Type.(*ast.StructType)

	return ok
}

func (t *Type) AsStruct() *Struct {
	s, ok := t.TypeSpec.Type.(*ast.StructType)

	if ok {
		return &Struct{t, s}
	}

	return nil
}

type InterfaceMap map[string]*Interface // +map

type Interface struct {
	*Type
	*ast.InterfaceType
}

func (intf *Interface) Name() string {
	return intf.TypeSpec.Name.Name
}

func (intf *Interface) Method(name string) *Method {
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

func (intf *Interface) Methods() MethodMap {
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
	*Interface
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

type StructMap map[string]*Struct // +map

type Struct struct {
	*Type
	*ast.StructType
}

func (s *Struct) Name() string {
	return s.TypeSpec.Name.Name
}

func (s *Struct) Field(name string) *Field {
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

func (s *Struct) Fields() FieldMap {
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
	*Struct
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

func (f *Field) Type() ast.Expr {
	return f.Field.Type
}

func (f *Field) Tag() reflect.StructTag {
	if f.Field.Tag == nil {
		return ""
	}

	return reflect.StructTag(strings.Trim(f.Field.Tag.Value, "`"))
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
