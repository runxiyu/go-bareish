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

func TestParseData(t *testing.T) {
	types, err := Parse(strings.NewReader(`
		type MyData data
		type MyData128 data<128>`))
	assert.NoError(t, err)
	assert.Len(t, types, 2)

	assert.IsType(t, new(UserDefinedType), types[0])
	udt := types[0].(*UserDefinedType)
	assert.Equal(t, "MyData", udt.Name())

	assert.IsType(t, new(DataType), udt.Type())
	dt := udt.Type().(*DataType)
	assert.Equal(t, DataSlice, dt.Kind())

	assert.IsType(t, new(UserDefinedType), types[1])
	udt = types[1].(*UserDefinedType)
	assert.Equal(t, "MyData128", udt.Name())

	assert.IsType(t, new(DataType), udt.Type())
	dt = udt.Type().(*DataType)
	assert.Equal(t, DataArray, dt.Kind())
	assert.Equal(t, uint(128), dt.Length())
}

func TestParseMap(t *testing.T) {
	types, err := Parse(strings.NewReader(`type MyMap map[u8]string`))
	assert.NoError(t, err)
	assert.Len(t, types, 1)

	assert.IsType(t, new(UserDefinedType), types[0])
	udt := types[0].(*UserDefinedType)
	assert.Equal(t, "MyMap", udt.Name())

	assert.IsType(t, new(MapType), udt.Type())
	mt := udt.Type().(*MapType)
	assert.Equal(t, Map, mt.Kind())

	assert.Equal(t, U8, mt.Key().Kind())
	assert.Equal(t, String, mt.Value().Kind())
}

func TestParseArrays(t *testing.T) {
	types, err := Parse(strings.NewReader(`
		type MyArray [128]string
		type MySlice []string`))
	assert.NoError(t, err)
	assert.Len(t, types, 2)

	assert.IsType(t, new(UserDefinedType), types[0])
	udt := types[0].(*UserDefinedType)
	assert.Equal(t, "MyArray", udt.Name())

	assert.IsType(t, new(ArrayType), udt.Type())
	at := udt.Type().(*ArrayType)
	assert.Equal(t, Array, at.Kind())
	assert.Equal(t, uint(128), at.Length())
	assert.Equal(t, String, at.Member().Kind())

	assert.IsType(t, new(UserDefinedType), types[1])
	udt = types[1].(*UserDefinedType)
	assert.Equal(t, "MySlice", udt.Name())

	assert.IsType(t, new(ArrayType), udt.Type())
	at = udt.Type().(*ArrayType)
	assert.Equal(t, Slice, at.Kind())
	assert.Equal(t, String, at.Member().Kind())
}

func TestParseStruct(t *testing.T) {
	types, err := Parse(strings.NewReader(`
		type MyStruct {
			x: i32
			y: i32
			z: i32
		}`))
	assert.NoError(t, err)
	assert.Len(t, types, 1)

	ty := types[0]
	assert.IsType(t, new(UserDefinedType), ty)
	udt := ty.(*UserDefinedType)
	assert.Equal(t, "MyStruct", udt.Name())

	assert.IsType(t, new(StructType), udt.Type())
	st := udt.Type().(*StructType)
	assert.Equal(t, Struct, st.Kind())
	assert.Len(t, st.Fields(), 3)

	f := st.Fields()[0]
	assert.Equal(t, "x", f.Name())
	assert.Equal(t, I32, f.Type().Kind())

	f = st.Fields()[1]
	assert.Equal(t, "y", f.Name())
	assert.Equal(t, I32, f.Type().Kind())

	f = st.Fields()[2]
	assert.Equal(t, "z", f.Name())
	assert.Equal(t, I32, f.Type().Kind())
}
