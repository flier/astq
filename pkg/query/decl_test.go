package query

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExampleTypeDeclMap_WithTag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

// +tag map:""
type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors|parser.ParseComments)

	fmt.Println(FromFile(f).TypeDecls().WithTag("map").Keys())
	// Output: [Foo]
}

func ExampleStructDef_Methods() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Point struct {
	x, y float64
}

func Foo()

func Bar(name string)

func (p Point) Length() float64 {
	return math.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
	p.x *= factor
	p.y *= factor
}
`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Struct("Point").Methods().Keys()))

	// Output: [Length Scale]
}
