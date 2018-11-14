package query

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExampleFieldMap_WithPrefix() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string
	Bar string
}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Struct("S").Fields().WithPrefix("F").Keys())
	// Output: [Foo]
}

func ExampleFieldMap_WithSuffix() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string
	Bar string
}

`, parser.AllErrors)

	fmt.Println(FromFile(f).Struct("S").Fields().WithSuffix("o").Keys())
	// Output: [Foo]
}

func ExampleFieldMap_WithTag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	Foo string `+"`species:\"gopher\" color:\"red\"`"+`
	Bar string `+"`species:\"gopher\" color:\"\"`"+`
	Qux string `+"`species:\"gopher\"`"+`
}

`, parser.AllErrors)

	fmt.Println(
		FromFile(f).Struct("S").Fields().WithTagValue("color", "red").Keys(),
		Sorted(FromFile(f).Struct("S").Fields().WithTag("color").Keys()),
		FromFile(f).Struct("S").Fields().WithoutTag("color").Keys(),
	)
	// Output: [Foo] [Bar Foo] [Qux]
}
