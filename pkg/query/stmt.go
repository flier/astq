package query

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Stmt interface {
	fmt.Stringer
}

type AstStmt struct {
	ast.Stmt
}

func asStmt(stmt ast.Stmt) Stmt {
	if stmt == nil {
		return nil
	}

	switch s := stmt.(type) {
	case *ast.BadStmt:
		return &BadStmt{&AstStmt{s}, s}
	case *ast.DeclStmt:
		return &DeclStmt{&AstStmt{s}, s}
	case *ast.EmptyStmt:
		return &EmptyStmt{&AstStmt{s}, s}
	case *ast.LabeledStmt:
		return &LabeledStmt{&AstStmt{s}, s}
	case *ast.ExprStmt:
		return &ExprStmt{&AstStmt{s}, s}
	case *ast.SendStmt:
		return &SendStmt{&AstStmt{s}, s}
	case *ast.IncDecStmt:
		return &IncDecStmt{&AstStmt{s}, s}
	case *ast.AssignStmt:
		return &AssignStmt{&AstStmt{s}, s}
	case *ast.GoStmt:
		return &GoStmt{&AstStmt{s}, s}
	case *ast.DeferStmt:
		return &DeferStmt{&AstStmt{s}, s}
	case *ast.ReturnStmt:
		return &ReturnStmt{&AstStmt{s}, s}
	case *ast.BranchStmt:
		return &BranchStmt{&AstStmt{s}, s}
	case *ast.BlockStmt:
		return &BlockStmt{&AstStmt{s}, s}
	case *ast.IfStmt:
		return &IfStmt{&AstStmt{s}, s}
	case *ast.CaseClause:
		return &CaseClause{&AstStmt{s}, s}
	case *ast.SwitchStmt:
		return &SwitchStmt{&AstStmt{s}, s}
	case *ast.TypeSwitchStmt:
		return &TypeSwitchStmt{&AstStmt{s}, s}
	case *ast.CommClause:
		return &CommClause{&AstStmt{s}, s}
	case *ast.SelectStmt:
		return &SelectStmt{&AstStmt{s}, s}
	case *ast.ForStmt:
		return &ForStmt{&AstStmt{s}, s}
	case *ast.RangeStmt:
		return &RangeStmt{&AstStmt{s}, s}
	default:
		panic(s)
	}
}

type BadStmt struct {
	*AstStmt
	*ast.BadStmt
}

func (s *BadStmt) String() string {
	return fmt.Sprintf("BAD[%d:%d]", s.From, s.To)
}

type DeclStmt struct {
	*AstStmt
	*ast.DeclStmt
}

func (s *DeclStmt) Decl() *GenDecl {
	return &GenDecl{s.DeclStmt.Decl.(*ast.GenDecl)}
}

func (s *DeclStmt) String() string {
	return s.Decl().String()
}

type EmptyStmt struct {
	*AstStmt
	*ast.EmptyStmt
}

func (s *EmptyStmt) IsImplicit() bool {
	return s.EmptyStmt.Implicit
}

func (s *EmptyStmt) String() string {
	if s.EmptyStmt.Implicit {
		return ""
	}

	return ";"
}

type LabeledStmt struct {
	*AstStmt
	*ast.LabeledStmt
}

func (s *LabeledStmt) Name() string {
	return s.LabeledStmt.Label.Name
}

func (s *LabeledStmt) Stmt() Stmt {
	return asStmt(s.LabeledStmt.Stmt)
}

func (s *LabeledStmt) String() string {
	return s.Name() + ":\n" + s.Stmt().String()
}

type ExprStmt struct {
	*AstStmt
	*ast.ExprStmt
}

func (s *ExprStmt) Expr() Expr { return asExpr(s.ExprStmt.X) }

func (s *ExprStmt) String() string {
	return s.Expr().String()
}

type SendStmt struct {
	*AstStmt
	*ast.SendStmt
}

func (s *SendStmt) Chan() Expr  { return asExpr(s.SendStmt.Chan) }
func (s *SendStmt) Value() Expr { return asExpr(s.SendStmt.Value) }

func (s *SendStmt) String() string {
	return s.Chan().String() + " <- " + s.Value().String()
}

type IncDecStmt struct {
	*AstStmt
	*ast.IncDecStmt
}

func (s *IncDecStmt) Token() string { return s.IncDecStmt.Tok.String() }
func (s *IncDecStmt) IsInc() bool   { return s.IncDecStmt.Tok == token.INC }
func (s *IncDecStmt) IsDec() bool   { return s.IncDecStmt.Tok == token.DEC }
func (s *IncDecStmt) Expr() Expr    { return asExpr(s.IncDecStmt.X) }

func (s *IncDecStmt) String() string {
	return s.Expr().String() + s.Token()
}

type AssignStmt struct {
	*AstStmt
	*ast.AssignStmt
}

func (s *AssignStmt) Lhs() (exprs []Expr) {
	for _, expr := range s.AssignStmt.Lhs {
		exprs = append(exprs, asExpr(expr))
	}

	return
}

func (s *AssignStmt) Token() string  { return s.AssignStmt.Tok.String() }
func (s *AssignStmt) IsAssign() bool { return s.AssignStmt.Tok == token.ASSIGN }
func (s *AssignStmt) IsDefine() bool { return s.AssignStmt.Tok == token.DEFINE }

func (s *AssignStmt) Rhs() (exprs []Expr) {
	for _, expr := range s.AssignStmt.Rhs {
		exprs = append(exprs, asExpr(expr))
	}

	return
}

func (s *AssignStmt) String() string {
	var lhs, rhs []string

	for _, expr := range s.Lhs() {
		lhs = append(lhs, expr.String())
	}
	for _, expr := range s.Rhs() {
		rhs = append(rhs, expr.String())
	}

	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), s.Token(), strings.Join(rhs, ", "))
}

type GoStmt struct {
	*AstStmt
	*ast.GoStmt
}

func (s *GoStmt) Call() *CallExpr {
	return &CallExpr{&AstExpr{s.GoStmt.Call}, s.GoStmt.Call}
}

func (s *GoStmt) String() string {
	return "go " + s.Call().String()
}

type DeferStmt struct {
	*AstStmt
	*ast.DeferStmt
}

func (s *DeferStmt) Call() *CallExpr {
	return &CallExpr{&AstExpr{s.DeferStmt.Call}, s.DeferStmt.Call}
}

func (s *DeferStmt) String() string {
	return "defer " + s.Call().String()
}

type ReturnStmt struct {
	*AstStmt
	*ast.ReturnStmt
}

func (s *ReturnStmt) Results() (results []Expr) {
	if s.ReturnStmt.Results != nil {
		for _, result := range s.ReturnStmt.Results {
			results = append(results, asExpr(result))
		}
	}

	return
}

func (s *ReturnStmt) String() string {
	if s.ReturnStmt.Results == nil {
		return "return"
	}
	var exprs []string

	for _, expr := range s.Results() {
		exprs = append(exprs, expr.String())
	}

	return "return " + strings.Join(exprs, ", ")
}

type BranchStmt struct {
	*AstStmt
	*ast.BranchStmt
}

func (s *BranchStmt) Token() string {
	return s.Tok.String()
}

func (s *BranchStmt) Label() string {
	if s.BranchStmt.Label == nil {
		return ""
	}

	return s.BranchStmt.Label.Name
}

func (s *BranchStmt) String() string {
	if s.BranchStmt.Label != nil {
		return s.Token() + " " + s.Label()
	}

	return s.Token()
}

type BlockStmt struct {
	*AstStmt
	*ast.BlockStmt
}

func (s *BlockStmt) Stmts() (stmts []Stmt) {
	for _, stmt := range s.BlockStmt.List {
		stmts = append(stmts, asStmt(stmt))
	}

	return
}

func (s *BlockStmt) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("{")

	for _, stmt := range s.Stmts() {
		buf.WriteString("\t" + stmt.String() + "\n")
	}

	buf.WriteString("}")

	return buf.String()
}

type IfStmt struct {
	*AstStmt
	*ast.IfStmt
}

func (s *IfStmt) Init() Stmt       { return asStmt(s.IfStmt.Init) }
func (s *IfStmt) Cond() Expr       { return asExpr(s.IfStmt.Cond) }
func (s *IfStmt) Else() Stmt       { return asStmt(s.IfStmt.Else) }
func (s *IfStmt) Body() *BlockStmt { return &BlockStmt{&AstStmt{s.IfStmt.Body}, s.IfStmt.Body} }

func (s *IfStmt) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("if ")

	if stmt := s.Init(); stmt != nil {
		buf.WriteString(stmt.String() + "; ")
	}

	buf.WriteString(s.Cond().String() + " " + s.Body().String())

	if stmt := s.Else(); stmt != nil {
		buf.WriteString(" else " + stmt.String())
	}

	return buf.String()
}

type CaseClause struct {
	*AstStmt
	*ast.CaseClause
}

func (c *CaseClause) IsDefault() bool { return c.CaseClause.List == nil }

func (c *CaseClause) Exprs() (exprs []Expr) {
	if c.CaseClause.List != nil {
		for _, expr := range c.CaseClause.List {
			exprs = append(exprs, asExpr(expr))
		}
	}

	return
}

func (c *CaseClause) Body() (stmts []Stmt) {
	for _, stmt := range c.CaseClause.Body {
		stmts = append(stmts, asStmt(stmt))
	}

	return
}

func (c *CaseClause) String() string {
	buf := new(bytes.Buffer)

	if c.IsDefault() {
		buf.WriteString("default:\n")
	} else {
		var exprs []string

		for _, expr := range c.Exprs() {
			exprs = append(exprs, expr.String())
		}

		buf.WriteString("cast " + strings.Join(exprs, ", ") + ":\n")
	}

	for _, stmt := range c.Body() {
		buf.WriteString("\t" + stmt.String() + "\n")
	}

	return buf.String()
}

type SwitchStmt struct {
	*AstStmt
	*ast.SwitchStmt
}

func (s *SwitchStmt) Init() Stmt { return asStmt(s.SwitchStmt.Init) }
func (s *SwitchStmt) Tag() Stmt  { return asExpr(s.SwitchStmt.Tag) }
func (s *SwitchStmt) Body() *BlockStmt {
	return &BlockStmt{&AstStmt{s.SwitchStmt.Body}, s.SwitchStmt.Body}
}

func (s *SwitchStmt) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("switch ")

	if s := s.Init(); s != nil {
		buf.WriteString(s.String() + "; ")
	}

	if tag := s.Tag(); tag != nil {
		buf.WriteString(tag.String() + " ")
	}

	buf.WriteString(s.Body().String())

	return buf.String()
}

type TypeSwitchStmt struct {
	*AstStmt
	*ast.TypeSwitchStmt
}

func (s *TypeSwitchStmt) Init() Stmt   { return asStmt(s.TypeSwitchStmt.Init) }
func (s *TypeSwitchStmt) Assign() Stmt { return asStmt(s.TypeSwitchStmt.Assign) }
func (s *TypeSwitchStmt) Body() *BlockStmt {
	return &BlockStmt{&AstStmt{s.TypeSwitchStmt.Body}, s.TypeSwitchStmt.Body}
}

func (s *TypeSwitchStmt) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("switch ")

	if s := s.Init(); s != nil {
		buf.WriteString(s.String() + "; ")
	}

	buf.WriteString(s.Assign().String() + " " + s.Body().String())

	return buf.String()
}

type CommClause struct {
	*AstStmt
	*ast.CommClause
}

func (c *CommClause) IsDefault() bool { return c.CommClause.Comm == nil }
func (c *CommClause) Comm() Stmt      { return asStmt(c.CommClause.Comm) }

func (c *CommClause) Body() (stmts []Stmt) {
	for _, stmt := range c.CommClause.Body {
		stmts = append(stmts, asStmt(stmt))
	}

	return
}

func (c *CommClause) String() string {
	buf := new(bytes.Buffer)

	if c.IsDefault() {
		buf.WriteString("default:\n")
	} else {
		buf.WriteString("cast " + c.Comm().String() + ":\n")
	}

	for _, stmt := range c.Body() {
		buf.WriteString("\t" + stmt.String() + "\n")
	}

	return buf.String()
}

type SelectStmt struct {
	*AstStmt
	*ast.SelectStmt
}

func (s *SelectStmt) Body() *BlockStmt {
	return &BlockStmt{&AstStmt{s.SelectStmt.Body}, s.SelectStmt.Body}
}

func (s *SelectStmt) String() string {
	return "select " + s.Body().String()
}

type ForStmt struct {
	*AstStmt
	*ast.ForStmt
}

func (s *ForStmt) Init() Stmt       { return asStmt(s.ForStmt.Init) }
func (s *ForStmt) Cond() Expr       { return asExpr(s.ForStmt.Cond) }
func (s *ForStmt) Post() Stmt       { return asStmt(s.ForStmt.Post) }
func (s *ForStmt) Body() *BlockStmt { return &BlockStmt{&AstStmt{s.ForStmt.Body}, s.ForStmt.Body} }

func (s *ForStmt) String() string {
	var init, cond, post string

	if i := s.Init(); i != nil {
		init = i.String()
	}
	if e := s.Cond(); e != nil {
		cond = e.String()
	}
	if p := s.Post(); p != nil {
		post = p.String()
	}

	return fmt.Sprintf("for %s; %s; %s %s", init, cond, post, s.Body())
}

type RangeStmt struct {
	*AstStmt
	*ast.RangeStmt
}

func (s *RangeStmt) Key() Expr        { return asExpr(s.RangeStmt.Key) }
func (s *RangeStmt) Value() Expr      { return asExpr(s.RangeStmt.Value) }
func (s *RangeStmt) Receiver() Expr   { return asExpr(s.RangeStmt.X) }
func (s *RangeStmt) Body() *BlockStmt { return &BlockStmt{&AstStmt{s.RangeStmt.Body}, s.RangeStmt.Body} }

func (s *RangeStmt) String() string {
	key := s.Key()
	value := s.Value()

	var expr string

	if key != nil && value != nil {
		expr = fmt.Sprintf("%s, %s := ", key, value)
	} else if key != nil {
		expr = key.String() + " := "
	}

	return fmt.Sprintf("for %srange %s %s", expr, s.Receiver(), s.Body())
}
