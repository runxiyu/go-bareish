package schema

import (
	"errors"
	"fmt"
	"io"
)

// Returned when the lexer encounters an unexpected token
type ErrUnexpectedToken struct {
	token    Token
	expected string
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("Unexpected token '%s'; expected %s",
		e.token.String(), e.expected)
}

// Parses a BARE schema definition language document from the given reader and
// returns a list of the user-defined types it specifies.
func Parse(reader io.Reader) ([]SchemaType, error) {
	scanner := NewScanner(reader)
	var stypes []SchemaType
	for {
		st, err := parseSchemaType(scanner)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		stypes = append(stypes, st)
	}
	return stypes, nil
}

func parseSchemaType(scanner *Scanner) (SchemaType, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Token {
	case TTYPE:
		return parseUserType(scanner)
	case TENUM:
		return parseUserEnum(scanner)
	}

	return nil, &ErrUnexpectedToken{tok, "'type' or 'enum'"}
}

func parseUserType(scanner *Scanner) (SchemaType, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}
	if tok.Token != TNAME {
		return nil, &ErrUnexpectedToken{tok, "type name"}
	}

	udt := &UserDefinedType{name: tok.Value}
	udt.type_, err = parseType(scanner)
	if err != nil {
		return nil, err
	}

	return udt, nil
}

func parseUserEnum(scanner *Scanner) (SchemaType, error) {
	return nil, errors.New("TODO")
}

func parseType(scanner *Scanner) (Type, error) {
	tok, err := scanner.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Token {
	case TU8:
		return &PrimitiveType{U8}, nil
	case TU16:
		return &PrimitiveType{U16}, nil
	case TU32:
		return &PrimitiveType{U32}, nil
	case TU64:
		return &PrimitiveType{U64}, nil
	case TI8:
		return &PrimitiveType{I8}, nil
	case TI16:
		return &PrimitiveType{I16}, nil
	case TI32:
		return &PrimitiveType{I32}, nil
	case TI64:
		return &PrimitiveType{I64}, nil
	case TF32:
		return &PrimitiveType{F32}, nil
	case TF64:
		return &PrimitiveType{F64}, nil
	case TE8:
		return &PrimitiveType{E8}, nil
	case TE16:
		return &PrimitiveType{E16}, nil
	case TE32:
		return &PrimitiveType{E32}, nil
	case TE64:
		return &PrimitiveType{E64}, nil
	case TBOOL:
		return &PrimitiveType{Bool}, nil
	case TSTRING:
		return &PrimitiveType{String}, nil
	case TOPTIONAL:
		tok, err = scanner.Next()
		if tok.Token != TLANGLE {
			return nil, &ErrUnexpectedToken{tok, "<"}
		}

		st, err := parseType(scanner)
		if err != nil {
			return nil, err
		}

		tok, err = scanner.Next()
		if tok.Token != TRANGLE {
			return nil, &ErrUnexpectedToken{tok, ">"}
		}
		return &OptionalType{subtype: st}, nil
	case TDATA:
	case TMAP:
	case TNAME:
		return nil, errors.New("TODO")
	}

	return nil, &ErrUnexpectedToken{tok, "type"}
}
