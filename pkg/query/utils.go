package query

import (
	"bytes"
	"go/ast"
	"sort"
)

func AstDump(x interface{}) string {
	var buf bytes.Buffer

	if err := ast.Fprint(&buf, nil, x, ast.NotNilFilter); err != nil {
		panic(err)
	}

	return buf.String()
}

func Sorted(a []string) []string {
	sort.Strings(a)
	return a
}
