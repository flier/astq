package query

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExampleTypeMap_WithTag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

// +map
type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors|parser.ParseComments)

	fmt.Println(FromFile(f).TypeDecls().WithTag("map").Keys())
	// Output: [Foo]
}

func ExampleField_Tag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	F string `+"`species:\"gopher\" color:\"blue\"`"+`
}

`, parser.AllErrors)

	tag := FromFile(f).Struct("S").NamedField("F").Tag()

	fmt.Println(tag, tag.Get("species"), tag.Get("color"))
	// Output: species:"gopher" color:"blue" gopher blue
}

func ExampleInterfaceType_Methods() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Buf interface {
	Read(b Buffer) bool
	Write(b Buffer) bool
	Close()
}

`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Interface("Buf").Methods().Keys()))
	// Output: [Close Read Write]
}

func ExampleInterfaceType_Method() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Buf interface {
	Read(b Buffer) bool
	Write(b Buffer) bool
	Close()
}

`, parser.AllErrors)

	intf := FromFile(f).Interface("Buf")

	fmt.Println(intf.Method("Read").String())
	fmt.Println(intf.Method("Write").String())
	fmt.Println(intf.Method("Close").String())
	// Output: Read(b Buffer) bool
	// Write(b Buffer) bool
	// Close()
}
