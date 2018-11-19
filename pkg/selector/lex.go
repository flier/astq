package selector

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
	"unicode"
)

const EOF = 0

var tokens map[int]string

func init() {
	tokens := make(map[int]string)
	tokens[LAST] = "last"
	tokens[POSITION] = "position"
}

type queryLex struct {
	*bufio.Reader
}

func (lex *queryLex) Lex(lval *querySymType) int {
	for {
		c := lex.next()

		switch c {
		case '/', '[', ']', '(', ')', '*', '=', ':', '@', '.', '"':
			return int(c)

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			lex.UnreadRune()
			lval.num = lex.num()
			return NUM

		case ' ', '\t', '\n', '\r':
			// ignore

		case EOF:
			return int(EOF)

		default:
			if unicode.IsLetter(c) {
				lex.UnreadRune()
				lval.str = lex.str()
				return STR
			}

			log.Printf("unrecognized character %q", c)
		}
	}
}

func (lex *queryLex) Error(s string) {}

func (lex *queryLex) next() rune {
	c, _, err := lex.ReadRune()

	if err != nil {
		return EOF
	}

	if c == unicode.ReplacementChar {
		return lex.next()
	}

	return c
}

func (lex *queryLex) num() int {
	buf := new(bytes.Buffer)

	var c rune

	for {
		c = lex.next()

		if unicode.IsNumber(c) {
			buf.WriteRune(c)
		} else {
			break
		}
	}

	if c != EOF {
		lex.UnreadRune()
	}

	n, err := strconv.ParseInt(buf.String(), 10, 64)

	if err != nil {
		panic(err)
	}

	return int(n)
}

func (lex *queryLex) str() string {
	buf := new(bytes.Buffer)

	var c rune

	for {
		c = lex.next()

		if unicode.IsLetter(c) {
			buf.WriteRune(c)
		} else {
			break
		}
	}

	return buf.String()
}
