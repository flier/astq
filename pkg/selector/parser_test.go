package selector

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiPathes(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with multi pathes", func() {
			parsed, err := ParseQuery("/foo, ./bar")

			So(err, ShouldBeNil)
			So(parsed,
				ShouldResemble,
				Query{
					Path{&Step{&Axis{Dir: "/"}, "foo", nil}},
					Path{&Step{&Axis{Dir: "./"}, "bar", nil}},
				})
		})
	})
}

func TestMultiSteps(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with multi steps", func() {
			parsed, err := ParseQuery("/foo ./bar")

			So(err, ShouldBeNil)
			So(parsed,
				ShouldResemble,
				Query{Path{
					&Step{&Axis{Dir: "/"}, "foo", nil},
					&Step{&Axis{Dir: "./"}, "bar", nil},
				}})
		})
	})
}

var step_with_axis = map[string]Query{
	"*":                 Query{Path{&Step{Match: "*"}}},
	"/*":                Query{Path{&Step{&Axis{Dir: "/"}, "*", nil}}},
	"//*":               Query{Path{&Step{&Axis{Dir: "//"}, "*", nil}}},
	"./*":               Query{Path{&Step{&Axis{Dir: "./"}, "*", nil}}},
	".//*":              Query{Path{&Step{&Axis{Dir: ".//"}, "*", nil}}},
	"-/*":               Query{Path{&Step{&Axis{Dir: "-/"}, "*", nil}}},
	"-//*":              Query{Path{&Step{&Axis{Dir: "-//"}, "*", nil}}},
	"+/*":               Query{Path{&Step{&Axis{Dir: "+/"}, "*", nil}}},
	"+//*":              Query{Path{&Step{&Axis{Dir: "+//"}, "*", nil}}},
	"~/*":               Query{Path{&Step{&Axis{Dir: "~/"}, "*", nil}}},
	"~//*":              Query{Path{&Step{&Axis{Dir: "~//"}, "*", nil}}},
	"../*":              Query{Path{&Step{&Axis{Dir: "../"}, "*", nil}}},
	"..//*":             Query{Path{&Step{&Axis{Dir: "..//"}, "*", nil}}},
	"<//*":              Query{Path{&Step{&Axis{Dir: "<//"}, "*", nil}}},
	">//*":              Query{Path{&Step{&Axis{Dir: ">//"}, "*", nil}}},
	"/:foo bar":         Query{Path{&Step{&Axis{Dir: "/", Type: "foo"}, "bar", nil}}},
	"/:\"foo\" bar":     Query{Path{&Step{&Axis{Dir: "/", Type: "foo"}, "bar", nil}}},
	"/:foo bar [@name]": Query{Path{&Step{&Axis{Dir: "/", Type: "foo"}, "bar", &WithAttr{"name"}}}},
}

func TestStepWithAxis(t *testing.T) {
	Convey("Given a parser", t, func() {
		Convey("When parse query with axis", func() {
			for s, expected := range step_with_axis {
				parsed, err := ParseQuery(s)

				So(err, ShouldBeNil)
				So(parsed, ShouldResemble, expected)
			}
		})
	})
}
