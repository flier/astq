package selector

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var tokens = []rune{
	'[', ']', '(', ')', ':', '@', '.', '~', ',', '+', '-', '*', '/', '^', '%',
}

var operators = map[int]string{
	int('<'): "<",
	LSHIFT:   "<<",
	LTE:      "<=",
	int('>'): ">",
	RSHIFT:   ">>",
	GTE:      ">=",
	int('&'): "&",
	AND:      "&&",
	int('|'): "|",
	OR:       "||",
	int('='): "=",
	EQ:       "==",
	MATCH:    "=~",
	int('!'): "!",
	NE:       "!=",
	NONMATCH: "!~",
	int('?'): "?",
	ELSE_OR:  "?:",
}

var strs = map[string]string{
	"hello world":                     "hello world",
	"\a\b\f\n\r\t\v\\'\"\x20\xFF\123": `\a\b\f\n\r\t\v\\\'\"\x20\xfF\123`,
	"\u0020\u0020":                    `\u0020\U00000020`,
	"汉字":                              `汉字`,
}

func TestLexer(t *testing.T) {
	Convey("Given a lexer", t, func() {
		buf := new(bytes.Buffer)
		lexer := queryNewLexer(buf)

		So(lexer.errs, ShouldBeNil)

		lval := new(querySymType)

		Convey("When parse a token", func() {
			for _, tok := range tokens {
				buf.WriteRune(tok)

				So(lexer.Lex(lval), ShouldEqual, tok)
			}
		})

		Convey("When parse an operator", func() {
			for op, s := range operators {
				buf.WriteString(s)

				So(lexer.Lex(lval), ShouldEqual, op)
			}
		})

		Convey("When parse a number", func() {
			buf.WriteString("123")

			So(lexer.Lex(lval), ShouldEqual, NUM)
			So(lval.num, ShouldEqual, 123)
		})

		Convey("When parse a string", func() {
			for s, escaped := range strs {
				buf.WriteString(`"` + escaped + `"`)

				So(lexer.Lex(lval), ShouldEqual, STR)
				So(lval.str, ShouldEqual, s)
			}

			Convey("When parse an incomplete string", func() {
				buf.WriteString(`"hello`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})

			Convey("When parse a wrong escaped char", func() {
				buf.WriteString(`"\kkk"`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})

			Convey("When parse a wrong escaped hex char", func() {
				buf.WriteString(`"\xKK"`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})

			Convey("When parse a wrong escaped oct char", func() {
				buf.WriteString(`"\788"`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})

			Convey("When parse a wrong too large oct char", func() {
				buf.WriteString(`"\777"`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})

			Convey("When parse a wrong escaped rune", func() {
				buf.WriteString(`"\U00110000"`)

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})
		})

		Convey("When parse a regexp", func() {
			buf.WriteString("`([a-z]+)\\[([0-9]+)\\]`")

			So(lexer.Lex(lval), ShouldEqual, REGEXP)
			So(lval.regexp.FindAllStringSubmatch("adam[23] snakey eve[7]", -1),
				ShouldResemble,
				[][]string{{"adam[23]", "adam", "23"}, {"eve[7]", "eve", "7"}})

			Convey("When parse an incomplete regexp", func() {
				buf.WriteString("`hello")

				So(lexer.Lex(lval), ShouldEqual, ERR)
				So(lval.err, ShouldEqual, errSyntax)
			})
		})

		Convey("When parse an ID", func() {
			buf.WriteString("hello _world 测试")

			So(lexer.Lex(lval), ShouldEqual, ID)
			So(lval.str, ShouldEqual, "hello")

			So(lexer.Lex(lval), ShouldEqual, ID)
			So(lval.str, ShouldEqual, "_world")

			So(lexer.Lex(lval), ShouldEqual, ID)
			So(lval.str, ShouldEqual, "测试")
		})

		Convey("When parse a unexpected rune", func() {
			buf.WriteString("#")

			So(lexer.Lex(lval), ShouldEqual, ERR)
		})
	})
}
