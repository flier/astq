package selector

import (
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

var exprWithIfElse = map[string]*Step{
	"* [@value ?: 1]":    &Step{Match: "*", Filter: &Cond{Cond: &WithAttr{"value"}, Else: Num(1)}},
	"* [@value ? 1 : 0]": &Step{Match: "*", Filter: &Cond{Cond: &WithAttr{"value"}, Then: Num(1), Else: Num(0)}},
}

func TestCond(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with if-else", func() {
			for q, expected := range exprWithIfElse {
				parsed, err := ParseQuery(q)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, Query{Path{expected}})
				So(parsed.String(), ShouldEqual, q)
			}
		})
	})
}
