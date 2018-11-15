package query

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

func ExampleFile_Type() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors)

	fmt.Println(FromFile(f).TypeDecl("Bar").Name())
	// Output: Bar
}

func ExampleFile_Types() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

type Foo map[string]interface{}
type Bar interface {}

`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).TypeDecls().Keys()))
	// Output: [Bar Foo]
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

func ExampleFile_Imports() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

import  (
	"lib/a"
	m "lib/b"
	. "lib/c"
)
`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Imports().Keys()))
	// Output: [lib/a lib/b lib/c]
}

func ExampleFile_Import() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

import  (
	"lib/a"
	m "lib/b"
	. "lib/c"
)
`, parser.AllErrors)

	fmt.Println(
		FromFile(f).Import("lib/a").Name(),
		FromFile(f).Import("lib/b").Name(),
		FromFile(f).Import("lib/c").Name(),
	)
	// Output: a m .
}

func ExampleFile_Var() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	i       int
	u, v, s = 2.0, 3.0, "bar"
)
var re, im = complexSqrt(-1)
var _, found = entries[name]  // map lookup; only interested in "found"
`, parser.AllErrors)

	fmt.Println(
		FromFile(f).Var("i").Type(),
		FromFile(f).Var("U").Type(),
		FromFile(f).Var("V").Type(),
		FromFile(f).Var("W").Type(),
		FromFile(f).Var("k").Type(), FromFile(f).Var("k").Value(),
		FromFile(f).Var("x").Type(), FromFile(f).Var("x").Value(),
		FromFile(f).Var("y").Type(), FromFile(f).Var("y").Value(),
		FromFile(f).Var("u").Type(), FromFile(f).Var("u").Value(),
		FromFile(f).Var("v").Type(), FromFile(f).Var("v").Value(),
		FromFile(f).Var("s").Type(), FromFile(f).Var("s").Value(),
	)
	// Output: int float64 float64 float64 <nil> 0 float32 - 1 float32 - 2 <nil> 2.0 <nil> 3.0 <nil> "bar"
}

func ExampleFile_Vars() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	i       int
	u, v, s = 2.0, 3.0, "bar"
)
var re, im = complexSqrt(-1)
var _, found = entries[name]  // map lookup; only interested in "found"
`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Vars().Keys()))

	vars := make(map[string]bool)
	for _, v := range FromFile(f).Vars() {
		vars[v.String()] = true
	}

	var strs []string
	for s := range vars {
		strs = append(strs, s)
	}

	fmt.Println(strings.Join(Sorted(strs), "\n"))

	// Output: [U V W _ found i im k re s u v x y]
	// var (U, V, W) float64
	// var (_, found) = entries[name]
	// var (re, im) = complexSqrt ()
	// var (u, v, s) = (2.0, 3.0, "bar")
	// var (x, y) float32 = (- 1, - 2)
	// var i int
	// var k = 0
}

func ExampleFile_Consts() {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", `

package test

const Pi float64 = 3.14159265358979323846
const zero = 0.0         // untyped floating-point constant
const (
	size int64 = 1024
	eof        = -1  // untyped integer constant
)
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants
const u, v float32 = 0, 3    // u = 0.0, v = 3.0
`, parser.AllErrors)

	fmt.Println(Sorted(FromFile(f).Consts().Keys()))

	vars := make(map[string]bool)
	for _, v := range FromFile(f).Consts() {
		vars[v.String()] = true
	}

	var strs []string
	for s := range vars {
		strs = append(strs, s)
	}

	fmt.Println(strings.Join(Sorted(strs), "\n"))

	// Output: [Pi a b c eof size u v zero]
	// const (a, b, c) = (3, 4, "foo")
	// const (u, v) float32 = (0, 3)
	// const Pi float64 = 3.14159265358979323846
	// const eof = - 1
	// const size int64 = 1024
	// const zero = 0.0
}
