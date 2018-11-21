package selector

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Query []Path

func (q Query) String() string {
	var pathes []string

	for _, path := range q {
		pathes = append(pathes, path.String())
	}

	return strings.Join(pathes, ", ")
}

type Path []*Step

func (p Path) String() string {
	var steps []string

	for _, step := range p {
		steps = append(steps, step.String())
	}

	return strings.Join(steps, " ")
}

type Step struct {
	*Axis
	Match  string
	Result bool
	Filter Expr
}

func (s *Step) String() string {
	buf := new(bytes.Buffer)

	if s.Axis != nil {
		buf.WriteString(s.Axis.String())

		if len(s.Axis.Type) > 0 {
			buf.WriteString(" ")
		}
	}

	buf.WriteString(s.Match)

	if s.Result || s.Filter != nil {
		buf.WriteRune(' ')

		if s.Result {
			buf.WriteString("!")
		}

		if s.Filter != nil {
			buf.WriteString("[" + s.Filter.String() + "]")
		}
	}

	return buf.String()
}

type Axis struct {
	Dir  string
	Type string
}

func (a *Axis) String() string {
	if len(a.Type) > 0 {
		if !IsIdent(a.Type) {
			return fmt.Sprintf("%s:%s", a.Dir, strconv.Quote(a.Type))
		}

		return fmt.Sprintf("%s:%s", a.Dir, a.Type)
	}

	return a.Dir
}

type Expr interface {
	fmt.Stringer
}

type Cond struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (c *Cond) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(c.Cond.String() + " ")

	if c.Then != nil {
		buf.WriteString("? " + c.Then.String() + " ")
	} else {
		buf.WriteString("?")
	}

	buf.WriteString(": " + c.Else.String())

	return buf.String()
}

type Unary struct {
	Op string
	Expr
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Expr)
}

type Binary struct {
	Lhs Expr
	Op  string
	Rhs Expr
}

func (b *Binary) String() string {
	return fmt.Sprintf("%s %s %s", b.Lhs, b.Op, b.Rhs)
}

type FuncCall struct {
	ID   string
	Args []Expr
}

func (c *FuncCall) String() string {
	var args []string

	for _, arg := range c.Args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("%s(%s)", c.ID, strings.Join(args, ", "))
}

type Match struct {
	Expr
	*regexp.Regexp
}

func (m *Match) String() string {
	return fmt.Sprintf("%s =~ `%s`", m.Expr, m.Regexp)
}

type WithAttr struct {
	ID string
}

func (a *WithAttr) String() string {
	return "@" + a.ID
}

type Str string

func (s Str) String() string {
	return string(s)
}

type Num int64

func (n Num) String() string {
	return strconv.FormatInt(int64(n), 10)
}

type Keyword string

func (k Keyword) String() string {
	return string(k)
}

type QueryParam string

func (p QueryParam) String() string {
	return string(p)
}

func IsIdent(s string) bool {
	for _, c := range s {
		if unicode.IsControl(c) || unicode.IsSpace(c) {
			return false
		}
	}

	return true
}
