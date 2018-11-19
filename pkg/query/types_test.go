package query

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

func ExampleField_Tag() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type S struct {
	// +tag line:"dot"
	F string `+"`species:\"gopher\" color:\"blue\"`"+`
}

`, parser.AllErrors)

	tags := FromFile(f).Struct("S").NamedField("F").Tags()

	fmt.Println(tags, tags.Get("species"), tags.Get("color"))
	// Output: [species:"gopher" color:"blue"] gopher blue
}

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

func ExampleStructType_NamedFields() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo struct {
	x, y int
	u float32
	_ float32  // padding
	A *[]int
	F func()
}

type Bar struct {
	T1        // field name is T1
	*T2       // field name is T2
	P.T3      // field name is T3
	*P.T4     // field name is T4
	x, y int  // field names are x and y
}
`, parser.AllErrors)

	fmt.Println(strings.Replace(FromFile(f).Struct("Foo").String(), "\t", strings.Repeat(" ", 8), -1))
	fmt.Println(Sorted(FromFile(f).Struct("Foo").NamedFields().Keys()))
	fmt.Println(strings.Replace(FromFile(f).Struct("Bar").String(), "\t", strings.Repeat(" ", 8), -1))
	fmt.Println(Sorted(FromFile(f).Struct("Bar").NamedFields().Keys()))
	// Output: type Foo struct {
	//         x, y int
	//         u float32
	//         _ float32
	//         A *[]int
	//         F func ()
	// }
	// [A F _ u x y]
	// type Bar struct {
	//         T1
	//         *T2
	//         P.T3
	//         *P.T4
	//         x, y int
	// }
	// [T1 T2 T3 T4 x y]
}
