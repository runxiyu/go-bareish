package schema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrimitives(t *testing.T) {
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
		type MyF32 f32
		type MyF64 f64
		type MyBool bool
		type MyString string
	`))
	assert.NoError(t, err)
	assert.Len(t, types, len(reference))

	for i, ty := range types {
		ref := reference[i]
		assert.IsType(t, ty, new(UserDefinedType))

		udt := ty.(*UserDefinedType)
		assert.Equal(t, ref.name, udt.Name())
		assert.Equal(t, ref.kind, udt.Type().Kind())
	}
}

func TestParseOptional(t *testing.T) {
	types, err := Parse(strings.NewReader("type MyOptional optional<u32>"))
	assert.NoError(t, err)
	assert.Len(t, types, 1)

	assert.IsType(t, new(UserDefinedType), types[0])
	udt := types[0].(*UserDefinedType)
	assert.Equal(t, "MyOptional", udt.Name())

	assert.IsType(t, new(OptionalType), udt.Type())
	ot := udt.Type().(*OptionalType)
	assert.Equal(t, U32, ot.Subtype().Kind())
}
