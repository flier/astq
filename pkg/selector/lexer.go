package selector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-multierror"
)

const eof = 0

type queryLexerImpl struct {
	*bufio.Reader
	result Query
	errs   *multierror.Error
}

func queryNewLexer(r io.Reader) *queryLexerImpl {
	return &queryLexerImpl{Reader: bufio.NewReader(r)}
}

func (l *queryLexerImpl) Result() Query {
	return l.result
}

func (l *queryLexerImpl) Err() error {
	return l.errs.ErrorOrNil()
}

func (l *queryLexerImpl) Error(s string) {
	l.errs = multierror.Append(l.errs, errors.New(s))
}

func (l *queryLexerImpl) Lex(lval *querySymType) int {
	for {
		c := l.next()

		switch c {
		case '[', ']', '(', ')', ':', '@', '.', '~', ',', '+', '-', '*', '/', '^', '%':
			break

		case '<':
			switch l.next() {
			case '<':
				lval.str = "<<"
				return LSHIFT
			case '=':
				lval.str = "<="
				return LTE
			case eof:
			default:
				l.UnreadRune()
			}

		case '>':
			switch l.next() {
			case '>':
				lval.str = ">>"
				return RSHIFT
			case '=':
				lval.str = ">="
				return GTE
			case eof:
			default:
				l.UnreadRune()
			}

		case '&':
			switch l.next() {
			case '&':
				lval.str = "&&"
				return AND
			case eof:
			default:
				l.UnreadRune()
			}

		case '|':
			switch l.next() {
			case '|':
				lval.str = "||"
				return OR
			case eof:

			default:
				l.UnreadRune()
			}

		case '=':
			switch l.next() {
			case '=':
				lval.str = "=="
				return EQ
			case '~':
				lval.str = "=~"
				return MATCH
			case eof:
			default:
				l.UnreadRune()
			}

		case '!':
			switch l.next() {
			case '=':
				lval.str = "!="
				return NE
			case '~':
				lval.str = "!~"
				return NONMATCH
			case eof:
			default:
				l.UnreadRune()
			}

		case '?':
			switch l.next() {
			case ':':
				lval.str = "?:"
				return ELSE_OR
			case eof:
			default:
				l.UnreadRune()
			}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.UnreadRune()
			lval.num, lval.err = l.num()

			if lval.err != nil {
				return ERR
			}
			return NUM

		case ' ', '\t', '\r', '\n':
			continue

		case '"':
			lval.str, lval.err = l.str()

			if lval.err != nil {
				return ERR
			}
			return STR

		case '`':
			lval.regexp, lval.err = l.regexp()

			if lval.err != nil {
				return ERR
			}
			return REGEXP

		case eof:
			return int(eof)

		default:
			if unicode.IsLetter(c) || c == '_' {
				l.UnreadRune()
				lval.str = l.id()
				return ID
			}

			lval.err = fmt.Errorf("unrecognized character %q", c)
			return ERR
		}

		lval.str = string(c)
		return int(c)
	}
}

func (l *queryLexerImpl) next() rune {
	c, _, err := l.ReadRune()

	if err != nil {
		return eof
	}

	if c == unicode.ReplacementChar {
		return l.next()
	}

	return c
}

func (l *queryLexerImpl) num() (int64, error) {
	buf := new(bytes.Buffer)

	var c rune

	for {
		c = l.next()

		if unicode.IsNumber(c) {
			buf.WriteRune(c)
		} else {
			break
		}
	}

	if c != eof {
		l.UnreadRune()
	}

	return strconv.ParseInt(buf.String(), 10, 64)
}

func (l *queryLexerImpl) id() string {
	buf := new(bytes.Buffer)

	var c rune

L:
	for {
		c = l.next()

		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
			buf.WriteRune(c)
		} else {
			break L
		}
	}

	if c != eof {
		l.UnreadRune()
	}

	return buf.String()
}

func (l *queryLexerImpl) str() (s string, err error) {
	buf := new(bytes.Buffer)

L:
	for {
		switch c := l.next(); {
		case c == eof:
			err = fmt.Errorf("incomplete string: \"%s\"", buf.String())
			break L

		case c == '"':
			s = buf.String()
			break L

		case c == '\\':
			switch c = l.next(); c {
			case 'a':
				buf.WriteRune('\a')
			case 'b':
				buf.WriteRune('\b')
			case 'f':
				buf.WriteRune('\f')
			case 'n':
				buf.WriteRune('\n')
			case 'r':
				buf.WriteRune('\r')
			case 't':
				buf.WriteRune('\t')
			case 'v':
				buf.WriteRune('\v')
			case '\\':
				buf.WriteRune('\\')
			case '\'':
				buf.WriteRune('\'')
			case '"':
				buf.WriteRune('"')
			case 'x', 'u', 'U':
				n := 0

				switch c {
				case 'x':
					n = 2
				case 'u':
					n = 4
				case 'U':
					n = 8
				}

				var v int

				for i := 0; i < n; i++ {
					c = l.next()
					n, ok := unhex(c)

					if !ok {
						err = fmt.Errorf("invalid hex char: '%c'", c)
						break L
					}

					v = (v << 4) | n
				}

				if c == 'x' {
					buf.WriteByte(byte(v))
				} else {
					if v > utf8.MaxRune {
						err = fmt.Errorf("invalid UNICODE rune: '\\U%08x'", v)
						break L
					}

					buf.WriteRune(rune(v))
				}

			case '0', '1', '2', '3', '4', '5', '6', '7':
				v := c - '0'

				for i := 0; i < 2; i++ {
					c = l.next()
					n := c - '0'

					if n < 0 || n > 7 {
						err = fmt.Errorf("invalid octal digit: '%c'", c)
						break L
					}

					v = (v << 3) | n
				}

				if v > 255 {
					err = fmt.Errorf("invalid octal value: %d", v)
					break L
				}

				buf.WriteRune(v)
			default:
				err = fmt.Errorf("unexpected escaped char: '%c'", c)
				break L
			}
		default:
			buf.WriteRune(c)
		}
	}

	return
}

func unhex(c rune) (v int, ok bool) {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0'), true
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10, true
	case 'A' <= c && c <= 'F':
		return int(c-'A') + 10, true
	}
	return
}

func (l *queryLexerImpl) regexp() (re *Regexp, err error) {
	buf := new(bytes.Buffer)

L:
	for {
		switch c := l.next(); c {
		case eof:
			err = fmt.Errorf("incomplete regex: `%s`", buf.String())
			break L

		case '`':
			r, err := regexp.Compile(buf.String())

			if err != nil {
				return nil, err
			}

			return &Regexp{r}, nil

		default:
			buf.WriteRune(c)
		}
	}

	return
}
