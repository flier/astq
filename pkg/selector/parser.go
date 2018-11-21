package selector

import (
	"bytes"
	"io"
)

//go:generate goyacc -o query.go -p "query" query.y

func ParseQuery(s string) (Query, error) {
	return Parse(bytes.NewBuffer([]byte(s)))
}

func Parse(r io.Reader) (Query, error) {
	lexer := queryNewLexer(r)

	if queryParse(lexer) == 0 {
		return lexer.Result(), nil
	}

	return nil, lexer.Err()
}
