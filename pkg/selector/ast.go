package selector

type Query []Path

type Path []*Step

type Step struct {
	*Axis
	Match  string
	Filter Expr
}

type Expr interface {
}

type AxisDirection int

const (
	// DirectChild for direct child nodes
	DirectChild AxisDirection = iota
	// AnyDescendant for any descendant nodes
	AnyDescendant
	// CurrentDirectChild for current node plus direct child nodes
	CurrentDirectChild
	// CurrentAnyDescendant for current node plus any descendant nodes
	CurrentAnyDescendant
	// DirectLeftSibling for direct left sibling node, or
	DirectLeftSibling
	// AnyLeftSibling for any left sibling nodes
	AnyLeftSibling
	// DirectRightSibling for direct right sibling node
	DirectRightSibling
	// AnyRightSibling for any right sibling nodes
	AnyRightSibling
	// DirectLeftAndRightSibling  for direct left and right sibling nodes
	DirectLeftAndRightSibling
	// AnyLeftAndRightSibling for all left and right sibling nodes
	AnyLeftAndRightSibling
	// DirectParent for direct parent node
	DirectParent
	// AnyParent for any parent nodes
	AnyParent
	// AnyPreceding for any preceding nodes
	AnyPreceding
	// AnyFollowing for any following nodes
	AnyFollowing
)

type Axis struct {
	Direction AxisDirection
	Type      string
}

type Condition struct {
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

type Attr struct {
	ID string
}

type QueryParam struct {
	ID string
}

type Str string

type Num int

type Keyword string
