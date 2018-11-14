package query

import (
	"fmt"
	"go/parser"
	"go/token"
	"sort"
)

func Sorted(a []string) []string {
	sort.Strings(a)
	return a
}

func ExampleFile_Type() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Type("Bar").Name())
	// Output: Bar
}

func ExampleFile_Types() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Types().Keys()))
	// Output: [Bar Foo]
}

func ExampleTypeMap_WithTag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

// +map
type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors|parser.ParseComments)

	fmt.Println(FromFile(f).Types().WithTag("map").Keys())
	// Output: [Foo]
}

func ExampleFile_Interface() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo interface {
	Bar() error
}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Interface("Foo").Method("Bar").Name())
	// Output: Bar
}

func ExampleFile_Interfaces() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo interface {}
type Bar interface {}

`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Interfaces().Keys()))
	// Output: [Bar Foo]
}

func ExampleFile_Struct() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo struct {
	io.Closer
	Bar string
}

`, parser.AllErrors)

	fmt.Println(
		FromFile(f).Struct("Foo").Field("Closer").Path(),
		FromFile(f).Struct("Foo").Field("Bar").Type(),
	)
	// Output: io.Closer string
}

func ExampleFile_Structs() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo struct {}
type Bar struct {}

`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Structs().Keys()))
	// Output: [Bar Foo]
}

func ExampleField_Tag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	F string `+"`species:\"gopher\" color:\"blue\"`"+`
}

`, parser.AllErrors)

	tag := FromFile(f).Struct("S").Field("F").Tag()

	fmt.Println(tag, tag.Get("species"), tag.Get("color"))
	// Output: species:"gopher" color:"blue" gopher blue
}
