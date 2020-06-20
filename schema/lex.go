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
	TTYPE Token = iota
	TENUM

	// NAME is used for name, user-type-name, and enum-value-name.
	// Distinguishing between these requires context.
	TNAME
	TINTEGER

	TU8
	TU16
	TU32
	TU64
	TI8
	TI16
	TI32
	TI64
	TF32
	TF64
	TE8
	TE16
	TE32
	TE64
	TBOOL
	TSTRING
	TDATA
	TMAP
	TOPTIONAL

	// <
	TLANGLE
	// >
	TRANGLE
	// {
	TLBRACE
	// }
	TRBRACE
	// [
	TLBRACKET
	// ]
	TRBRACKET
	// (
	TLPAREN
	// )
	TRPAREN
	// ,
	TCOMMA
	// |
	TPIPE
	// =
	TEQUAL
	// :
	TCOLON
)

type Scanner struct {
	// TODO: track lineno/colno information and attach it to the tokens
	// returned, for better error reporting
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
			return TLANGLE, "", nil
		case '>':
			return TRANGLE, "", nil
		case '{':
			return TLBRACE, "", nil
		case '}':
			return TRBRACE, "", nil
		case '[':
			return TLBRACKET, "", nil
		case ']':
			return TRBRACKET, "", nil
		case '(':
			return TLPAREN, "", nil
		case ')':
			return TRPAREN, "", nil
		case ',':
			return TCOMMA, "", nil
		case '|':
			return TPIPE, "", nil
		case '=':
			return TEQUAL, "", nil
		case ':':
			return TCOLON, "", nil
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
		return TTYPE, "", nil
	case "enum":
		return TENUM, "", nil
	case "u8":
		return TU8, "", nil
	case "u16":
		return TU16, "", nil
	case "u32":
		return TU32, "", nil
	case "u64":
		return TU64, "", nil
	case "i8":
		return TI8, "", nil
	case "i16":
		return TI16, "", nil
	case "i32":
		return TI32, "", nil
	case "i64":
		return TI64, "", nil
	case "f32":
		return TF32, "", nil
	case "f64":
		return TF64, "", nil
	case "bool":
		return TBOOL, "", nil
	case "e8":
		return TE8, "", nil
	case "e16":
		return TE16, "", nil
	case "e32":
		return TE32, "", nil
	case "e64":
		return TE64, "", nil
	case "string":
		return TSTRING, "", nil
	case "data":
		return TDATA, "", nil
	case "optional":
		return TOPTIONAL, "", nil
	case "map":
		return TMAP, "", nil
	}

	return TNAME, tok, nil
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

	return TINTEGER, buf.String(), nil
}
