package schema

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanWords(t *testing.T) {
	cases := map[string]Token{
		"u8": U8,
		"u16": U16,
		"u32": U32,
		"u64": U64,
		"i8": I8,
		"i16": I16,
		"i32": I32,
		"i64": I64,
		"f32": F32,
		"f64": F64,
		"e8": E8,
		"e16": E16,
		"e32": E32,
		"e64": E64,
		"bool": BOOL,
		"string": STRING,
		"data": DATA,
		"map": MAP,
		"optional": OPTIONAL,
	}

	for input, reference := range cases {
		scanner := NewScanner(strings.NewReader(input))
		tok, val, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error")
		assert.Empty(t, val, "Expected Scan to return no value")
		assert.Equal(t, reference, tok,
			"Expected Scan to return reference value for %s", input)
		_, _, err = scanner.Next()
		assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
	}

	scanner := NewScanner(strings.NewReader("hello"))
	tok, val, err := scanner.Next()
	assert.NoError(t, err, "Expected Scan to return without error")
	assert.Equal(t, val, "hello", "Expected Scan to return value 'hello'")
	assert.Equal(t, NAME, tok, "Expected Scan to return NAME")
	_, _, err = scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}

func TestScanInteger(t *testing.T) {
	scanner := NewScanner(strings.NewReader("12345"))
	tok, val, err := scanner.Next()
	assert.NoError(t, err, "Expected Scan to return without error")
	assert.Equal(t, val, "12345", "Expected Scan to return value '12345'")
	assert.Equal(t, INTEGER, tok, "Expected Scan to return INTEGER")
	_, _, err = scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}

func TestScanSymbols(t *testing.T) {
	cases := map[string]Token{
		"<": LANGLE,
		">": RANGLE,
		"{": LBRACE,
		"}": RBRACE,
		"[": LBRACKET,
		"]": RBRACKET,
		"(": LPAREN,
		")": RPAREN,
	}

	for input, reference := range cases {
		scanner := NewScanner(strings.NewReader(input))
		tok, val, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error")
		assert.Empty(t, val, "Expected Scan to return no value")
		assert.Equal(t, reference, tok,
			"Expected Scan to return reference value for %s", input)
		_, _, err = scanner.Next()
		assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
	}
}

func TestScanSample(t *testing.T) {
	sample := `
	type PublicKey data<128>
	type Time string # ISO 8601

	enum Department e8 {
		ACCOUNTING,
		ADMINISTRATION,
		CUSTOMER_SERVICE,
		DEVELOPMENT,

		# Reserved for the CEO
		JSMITH = 99,
	}

	type Customer {
		name: string,
		email: string,
		address: Address,
		orders: []{
			orderId: i64,
			quantity: i32,
		},
		metadata: map[string]data,
	}

	type Person (Customer | Employee)`
	type Reference struct {
		tok Token
		val string
	}
	reference := []Reference{
		{TYPE, ""}, {NAME, "PublicKey"}, {DATA, ""},
			{LANGLE, ""}, {INTEGER, "128"}, {RANGLE, ""},
		{TYPE, ""}, {NAME, "Time"}, {STRING, ""},
		{ENUM, ""}, {NAME, "Department"}, {E8, ""}, {LBRACE, ""},
			{NAME, "ACCOUNTING"}, {COMMA, ""},
			{NAME, "ADMINISTRATION"}, {COMMA, ""},
			{NAME, "CUSTOMER_SERVICE"}, {COMMA, ""},
			{NAME, "DEVELOPMENT"}, {COMMA, ""},
			{NAME, "JSMITH"}, {EQUAL, ""}, {INTEGER, "99"}, {COMMA, ""},
		{RBRACE, ""},
		{TYPE, ""}, {NAME, "Customer"}, {LBRACE, ""},
		{NAME, "name"}, {COLON, ""}, {STRING, ""}, {COMMA, ""},
		{NAME, "email"}, {COLON, ""}, {STRING, ""}, {COMMA, ""},
		{NAME, "address"}, {COLON, ""}, {NAME, "Address"}, {COMMA, ""},
		{NAME, "orders"}, {COLON, ""}, {LBRACKET, ""}, {RBRACKET, ""}, {LBRACE, ""},
			{NAME, "orderId"}, {COLON, ""}, {I64, ""}, {COMMA, ""},
			{NAME, "quantity"}, {COLON, ""}, {I32, ""}, {COMMA, ""},
		{RBRACE, ""}, {COMMA, ""},
		{NAME, "metadata"}, {COLON, ""},
			{MAP, ""}, {LBRACKET, ""}, {STRING, ""}, {RBRACKET, ""},
				{DATA, ""}, {COMMA, ""},
		{RBRACE, ""},
		{TYPE, ""}, {NAME, "Person"},
			{LPAREN, ""}, {NAME, "Customer"},
			{PIPE, ""}, {NAME, "Employee"},
			{RPAREN, ""},
	}
	scanner := NewScanner(strings.NewReader(sample))
	for i, ref := range reference {
		tok, val, err := scanner.Next()
		assert.NoError(t, err, "Expected Scan to return without error for reference %d", i)
		assert.Equal(t, ref.tok, tok, "Expected Scan to return correct token for reference %d", i)
		assert.Equal(t, ref.val, val, "Expected Scan to return correct value for reference %d", i)
	}

	_, _, err := scanner.Next()
	assert.Equal(t, io.EOF, err, "Expected Scan to return EOF")
}
