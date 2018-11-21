package selector

import "regexp"

type Query []Path

type Path []*Step

type Step struct {
	*Axis
	Match  string
	Filter Expr
}

type Expr interface {
}

type Axis struct {
	Dir  string
	Type string
}

type Cond struct {
	Cond Expr
	Then Expr
	Else Expr
}

type Unary struct {
	Op string
	Expr
}

type Binary struct {
	Lhs Expr
	Op  string
	Rhs Expr
}

type FuncCall struct {
	ID   string
	Args []Expr
}

type Match struct {
	Expr
	*regexp.Regexp
}

type WithAttr struct {
	ID string
}

type SubQuery Path

type Str string

type Num int

type Keyword string

type QueryParam string
