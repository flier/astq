package query

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExampleNamedFieldMap_WithPrefix() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string
	Bar string
}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Struct("S").NamedFields().WithPrefix("F").Keys())
	// Output: [Foo]
}

func ExampleNamedFieldMap_WithSuffix() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string
	Bar string
}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Struct("S").NamedFields().WithSuffix("o").Keys())
	// Output: [Foo]
}

func ExampleNamedFieldMap_WithTag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string `+"`species:\"gopher\" color:\"red\"`"+`
	Bar string `+"`species:\"gopher\" color:\"\"`"+`
	Qux string `+"`species:\"gopher\"`"+`
}

`, parser.AllErrors)

	fmt.Println(
		FromFile(f).Struct("S").NamedFields().WithTagValue("color", "red").Keys(),
		Sorted(FromFile(f).Struct("S").NamedFields().WithTag("color").Keys()),
		FromFile(f).Struct("S").NamedFields().WithoutTag("color").Keys(),
	)
	// Output: [Foo] [Bar Foo] [Qux]
}
