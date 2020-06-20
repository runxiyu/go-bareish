package schema

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

// Returned when the lexer encounters an unexpected character
type ErrUnknownToken struct {
	token rune
}

func (e *ErrUnknownToken) Error() string {
	return fmt.Sprintf("Unknown token '%c'", e.token)
}

// A single lexographic token from a schema language token stream
type Token int

const (
	TYPE Token = iota
	ENUM

	// NAME is used for name, user-type-name, and enum-value-name.
	// Distinguishing between these requires context.
	NAME
	INTEGER

	U8
	U16
	U32
	U64
	I8
	I16
	I32
	I64
	F32
	F64
	E8
	E16
	E32
	E64
	BOOL
	STRING
	DATA
	MAP
	OPTIONAL

	// <
	LANGLE
	// >
	RANGLE
	// {
	LBRACE
	// }
	RBRACE
	// [
	LBRACKET
	// ]
	RBRACKET
	// (
	LPAREN
	// )
	RPAREN
	// ,
	COMMA
	// |
	PIPE
	// =
	EQUAL
	// :
	COLON
)

type Scanner struct {
	br *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(reader)}
}

// Returns the next token from the reader. If the token has a string associated
// with it (e.g. UserTypeName, Name, and Integer), the second return value is
// set to that string.
func (sc *Scanner) Next() (Token, string, error) {
	var (
		err error
		r   rune
	)

	for {
		r, _, err = sc.br.ReadRune()
		if err != nil {
			break
		}

		if unicode.IsSpace(r) {
			continue
		}
		if unicode.IsLetter(r) {
			sc.br.UnreadRune()
			return sc.scanWord()
		}
		if unicode.IsDigit(r) {
			sc.br.UnreadRune()
			return sc.scanInteger()
		}

		switch r {
		case '#':
			sc.br.ReadString('\n')
			continue
		case '<':
			return LANGLE, "", nil
		case '>':
			return RANGLE, "", nil
		case '{':
			return LBRACE, "", nil
		case '}':
			return RBRACE, "", nil
		case '[':
			return LBRACKET, "", nil
		case ']':
			return RBRACKET, "", nil
		case '(':
			return LPAREN, "", nil
		case ')':
			return RPAREN, "", nil
		case ',':
			return COMMA, "", nil
		case '|':
			return PIPE, "", nil
		case '=':
			return EQUAL, "", nil
		case ':':
			return COLON, "", nil
		}

		return 0, "", &ErrUnknownToken{r}
	}

	return 0, "", err
}

func (sc *Scanner) scanWord() (Token, string, error) {
	var buf bytes.Buffer

	for {
		r, _, err := sc.br.ReadRune()
		if err != nil  {
			if err == io.EOF {
				break
			}
			return 0, "", err
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			buf.WriteRune(r)
		} else {
			sc.br.UnreadRune()
			break
		}
	}

	tok := buf.String()
	switch tok {
	case "type":
		return TYPE, "", nil
	case "enum":
		return ENUM, "", nil
	case "u8":
		return U8, "", nil
	case "u16":
		return U16, "", nil
	case "u32":
		return U32, "", nil
	case "u64":
		return U64, "", nil
	case "i8":
		return I8, "", nil
	case "i16":
		return I16, "", nil
	case "i32":
		return I32, "", nil
	case "i64":
		return I64, "", nil
	case "f32":
		return F32, "", nil
	case "f64":
		return F64, "", nil
	case "bool":
		return BOOL, "", nil
	case "e8":
		return E8, "", nil
	case "e16":
		return E16, "", nil
	case "e32":
		return E32, "", nil
	case "e64":
		return E64, "", nil
	case "string":
		return STRING, "", nil
	case "data":
		return DATA, "", nil
	case "optional":
		return OPTIONAL, "", nil
	case "map":
		return MAP, "", nil
	}

	return NAME, tok, nil
}

func (sc *Scanner) scanInteger() (Token, string, error) {
	var buf bytes.Buffer

	for {
		r, _, err := sc.br.ReadRune()
		if err != nil  {
			if err == io.EOF {
				break
			}
			return 0, "", err
		}

		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			sc.br.UnreadRune()
			break
		}
	}

	return INTEGER, buf.String(), nil
}
