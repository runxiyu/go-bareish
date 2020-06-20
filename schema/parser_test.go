package schema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ParsePrimitives(t *testing.T) {
	type Reference struct {
		name string
		kind TypeKind
	}

	reference := []Reference{
		{ "MyU8", U8 },
		{ "MyU16", U16 },
		{ "MyU32", U32 },
		{ "MyU64", U64 },
		{ "MyI8", I8 },
		{ "MyI16", I16 },
		{ "MyI32", I32 },
		{ "MyI64", I64 },
		{ "MyF32", F32 },
		{ "MyF64", F64 },
		{ "MyBool", Bool },
		{ "MyString", String },
	}

	types, err := Parse(strings.NewReader(`
		type MyU8 u8
		type MyU16 u16
		type MyU32 u32
		type MyU64 u64
		type MyI8 i8
		type MyI16 i16
		type MyI32 i32
		type MyI64 i64
		type MyF32 i32
		type MyF64 i64
		type MyBool bool
		type MyString string
	`))
	assert.NoError(t, err, "Expected Parse to return without error")

	for i, ty := range types {
		ref := reference[i]
		assert.IsType(t, ty, new(*UserDefinedType), "Expected primitive type")

		udt, _ := ty.(*UserDefinedType)
		assert.Equal(t, udt.Name(), ref.name, "Incorrect type name")
		assert.Equal(t, udt.Type().Kind(), ref.kind, "Incorrect type kind")
	}
}
