package selector

type Query interface {
}

type DocElem struct {
	Query
}

type AllElems struct {
	Query
}

type ChildElems struct {
	Query
}

type WithAttr struct {
	Query
	Attr string
}

type WithAttrValue struct {
	Query
	Attr  string
	Value string
}

type WithIndex struct {
	Query
	Index int
}

type WithPosition struct {
	Query
	BinOp string
	Value int
}

type WithName struct {
	Name string
}
