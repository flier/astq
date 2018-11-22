package selector

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiPathes(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with multi pathes", func() {
			q := "/foo, ./bar"
			parsed, err := ParseQuery(q)

			So(err, ShouldBeNil)
			So(parsed,
				ShouldResemble,
				Query{
					Path{&Step{Axis: &Axis{Dir: "/"}, Match: "foo"}},
					Path{&Step{Axis: &Axis{Dir: "./"}, Match: "bar"}},
				})
			So(parsed.String(), ShouldEqual, q)
		})
	})
}

func TestMultiSteps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with multi steps", func() {
			q := "/foo ./bar"
			parsed, err := ParseQuery(q)

			So(err, ShouldBeNil)
			So(parsed,
				ShouldResemble,
				Query{Path{
					&Step{Axis: &Axis{Dir: "/"}, Match: "foo"},
					&Step{Axis: &Axis{Dir: "./"}, Match: "bar"},
				}})
			So(parsed.String(), ShouldEqual, q)
		})
	})
}

var stepWithAxis = map[string]*Step{
	"*":                     &Step{Match: "*"},
	"/*":                    &Step{Axis: &Axis{Dir: "/"}, Match: "*"},
	"//*":                   &Step{Axis: &Axis{Dir: "//"}, Match: "*"},
	"./*":                   &Step{Axis: &Axis{Dir: "./"}, Match: "*"},
	".//*":                  &Step{Axis: &Axis{Dir: ".//"}, Match: "*"},
	"-/*":                   &Step{Axis: &Axis{Dir: "-/"}, Match: "*"},
	"-//*":                  &Step{Axis: &Axis{Dir: "-//"}, Match: "*"},
	"+/*":                   &Step{Axis: &Axis{Dir: "+/"}, Match: "*"},
	"+//*":                  &Step{Axis: &Axis{Dir: "+//"}, Match: "*"},
	"~/*":                   &Step{Axis: &Axis{Dir: "~/"}, Match: "*"},
	"~//*":                  &Step{Axis: &Axis{Dir: "~//"}, Match: "*"},
	"../*":                  &Step{Axis: &Axis{Dir: "../"}, Match: "*"},
	"..//*":                 &Step{Axis: &Axis{Dir: "..//"}, Match: "*"},
	"<//*":                  &Step{Axis: &Axis{Dir: "<//"}, Match: "*"},
	">//*":                  &Step{Axis: &Axis{Dir: ">//"}, Match: "*"},
	"/:foo bar":             &Step{Axis: &Axis{Dir: "/", Type: "foo"}, Match: "bar"},
	"/:\"hello world\" bar": &Step{Axis: &Axis{Dir: "/", Type: "hello world"}, Match: "bar"},
	"/:foo bar [@name]":     &Step{Axis: &Axis{Dir: "/", Type: "foo"}, Match: "bar", Filter: &WithAttr{"name"}},
	"/:foo bar ![@name]":    &Step{Axis: &Axis{Dir: "/", Type: "foo"}, Match: "bar", Result: true, Filter: &WithAttr{"name"}},
}

func TestStepWithAxis(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with axis", func() {
			for s, expected := range stepWithAxis {
				parsed, err := ParseQuery(s)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{expected}})
				So(parsed.String(), ShouldEqual, s)
			}
		})
	})
}

func TestCond(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with if-else", func() {
			var exprs = map[string]Expr{
				"* [@value ?: 1]":    &Cond{Cond: &WithAttr{"value"}, Else: Num(1)},
				"* [@value ? 1 : 0]": &Cond{Cond: &WithAttr{"value"}, Then: Num(1), Else: Num(0)},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}

func TestLogicalOps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with logical operator", func() {
			var exprs = map[string]Expr{
				"* [@name && @value]": &Binary{&WithAttr{"name"}, "&&", &WithAttr{"value"}},
				"* [@name || @value]": &Binary{&WithAttr{"name"}, "||", &WithAttr{"value"}},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}

func TestBitwiseOps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with bitwise operator", func() {
			var exprs = map[string]Expr{
				"* [@name & @value]": &Binary{&WithAttr{"name"}, "&", &WithAttr{"value"}},
				"* [@name | @value]": &Binary{&WithAttr{"name"}, "|", &WithAttr{"value"}},
				"* [1 << 5]":         &Binary{Num(1), "<<", Num(5)},
				"* [8 >> 2]":         &Binary{Num(8), ">>", Num(2)},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}

func TestRelationalOps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with relational operator", func() {
			var exprs = map[string]Expr{
				"* [@name == @value]": &Binary{&WithAttr{"name"}, "==", &WithAttr{"value"}},
				"* [@name != @value]": &Binary{&WithAttr{"name"}, "!=", &WithAttr{"value"}},
				"* [1 > 5]":           &Binary{Num(1), ">", Num(5)},
				"* [1 >= 5]":          &Binary{Num(1), ">=", Num(5)},
				"* [8 < 2]":           &Binary{Num(8), "<", Num(2)},
				"* [8 <= 2]":          &Binary{Num(8), "<=", Num(2)},
				"* [@name =~ `\\w+`]": &Binary{&WithAttr{"name"}, "=~", &Regexp{regexp.MustCompile("\\w+")}},
				"* [@name !~ `\\w+`]": &Binary{&WithAttr{"name"}, "!~", &Regexp{regexp.MustCompile("\\w+")}},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}

func TestArithmethicalOps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with arithmethical operator", func() {
			var exprs = map[string]Expr{
				"* [1 + 5]": &Binary{Num(1), "+", Num(5)},
				"* [1 - 5]": &Binary{Num(1), "-", Num(5)},
				"* [8 * 2]": &Binary{Num(8), "*", Num(2)},
				"* [8 / 2]": &Binary{Num(8), "/", Num(2)},
				"* [8 % 2]": &Binary{Num(8), "%", Num(2)},
				"* [8 ^ 2]": &Binary{Num(8), "^", Num(2)},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}

func TestFuncCall(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with function call", func() {
			var exprs = map[string]Expr{
				"* [last()]":                               &FuncCall{"last", nil},
				"* [trim(\"hello world\")]":                &FuncCall{"trim", []Expr{Str("hello world")}},
				"* [index(\"hello world\", \"world\", 0)]": &FuncCall{"index", []Expr{Str("hello world"), Str("world"), Num(0)}},
			}

			for q, expected := range exprs {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{&Step{Match: "*", Filter: expected}}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}
