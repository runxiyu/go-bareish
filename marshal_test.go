package bare

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalValue(t *testing.T) {
	var (
		data []byte
		err  error

		u8  uint8   = 0x42
		u16 uint16  = 0xCAFE
		u32 uint32  = 0xDEADBEEF
		u64 uint64  = 0xCAFEBABEDEADBEEF
		u   uint    = 0xDEADBEEF
		i8  int8    = -42
		i16 int16   = -1234
		i32 int32   = -12345678
		i64 int64   = -12345678987654321
		i   int     = -12345678
		f32 float32 = 1337.42
		f64 float64 = 133713371337.42424242
		b   bool    = true
		str string  = "こんにちは、世界！"
	)

	payloads := [][]byte{
		[]byte{0x42},
		[]byte{0xFE, 0xCA},
		[]byte{0xEF, 0xBE, 0xAD, 0xDE},
		[]byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA},
		[]byte{0xEF, 0xBE, 0xAD, 0xDE},
		[]byte{0xD6},
		[]byte{0x2E, 0xFB},
		[]byte{0xB2, 0x9E, 0x43, 0xFF},
		[]byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF},
		[]byte{0xB2, 0x9E, 0x43, 0xFF},
		[]byte{0x71, 0x2D, 0xA7, 0x44},
		[]byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42},
		[]byte{0x01},
		[]byte{0x00},
		[]byte{0x1B, 0x00, 0x00, 0x00, 0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93,
			0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3, 0x80,
			0x81, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81},
	}

	data, err = Marshal(&u8)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[0], data)

	data, err = Marshal(&u16)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[1], data)

	data, err = Marshal(&u32)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[2], data)

	data, err = Marshal(&u64)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[3], data)

	data, err = Marshal(&u)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[4], data)

	data, err = Marshal(&i8)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[5], data)

	data, err = Marshal(&i16)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[6], data)

	data, err = Marshal(&i32)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[7], data)

	data, err = Marshal(&i64)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[8], data)

	data, err = Marshal(&i)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[9], data)

	data, err = Marshal(&f32)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[10], data)

	data, err = Marshal(&f64)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[11], data)

	data, err = Marshal(&b)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[12], data)

	b = false
	data, err = Marshal(&b)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[13], data)

	data, err = Marshal(&str)
	assert.Nil(t, err, "Expected Marshal to return without error")
	assert.Equal(t, payloads[14], data)
}

func TestMarshalOptional(t *testing.T) {
	var val *uint32
	data, err := Marshal(&val)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x00}, data)

	val = new(uint32)
	*val = 0xDEADBEEF
	data, err = Marshal(&val)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x01, 0xEF, 0xBE, 0xAD, 0xDE}, data)
}

func TestMarshalStruct(t *testing.T) {
	type Coordinates struct { X, Y, Z int }
	coords := Coordinates{1, 2, 3}
	reference := []byte{
		0x01, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x00, 0x00,
	}
	data, err := Marshal(&coords)
	assert.Nil(t, err)
	assert.Equal(t, reference, data)
}

func TestMarshalArray(t *testing.T) {
	val := [4]uint8{0x11, 0x22, 0x33, 0x44}
	reference := []byte{0x11, 0x22, 0x33, 0x44}
	data, err := Marshal(&val)
	assert.Nil(t, err)
	assert.Equal(t, reference, data)
}

func TestMarshalSlice(t *testing.T) {
	val := []uint8{0x11, 0x22, 0x33, 0x44}
	reference := []byte{0x04, 0x00, 0x00, 0x00, 0x11, 0x22, 0x33, 0x44}
	data, err := Marshal(&val)
	assert.Nil(t, err)
	assert.Equal(t, reference, data)
}

func TestMarshalMap(t *testing.T) {
	val := map[uint8]uint8{
		0x01: 0x11,
		0x02: 0x22,
		0x03: 0x33,
	}
	// Go maps are unordered, so any of these are considered valid
	references := [][]byte{
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x01, 0x11,
			0x02, 0x22,
			0x03, 0x33,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x01, 0x11,
			0x03, 0x33,
			0x02, 0x22,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x03, 0x33,
			0x02, 0x22,
			0x01, 0x11,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x01, 0x11,
			0x03, 0x33,
			0x02, 0x22,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x03, 0x33,
			0x01, 0x11,
			0x02, 0x22,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x02, 0x22,
			0x03, 0x33,
			0x01, 0x11,
		},
		[]byte{
			0x03, 0x00, 0x00, 0x00,
			0x02, 0x22,
			0x01, 0x11,
			0x03, 0x33,
		},
	}
	data, err := Marshal(&val)
	assert.Nil(t, err)
	var found bool
	for _, ref := range references {
		if reflect.DeepEqual(ref, data) {
			found = true
			break
		}
	}
	assert.True(t, found, "%x does not match reference", data)
}

type NameAge interface {
	Union
}

type Name string

func (n Name) UnionTag() uint8 {
	return 0
}

type Age int

func (a Age) UnionTag() uint8 {
	return 1
}

func TestMarshalUnion(t *testing.T) {
	var na NameAge = Name("Mary")
	ctx := NewContext()
	ctx.RegisterUnion(&na, *new(Name), *new(Age))

	data, err := ctx.Marshal(&na)
	assert.Nil(t, err)
	reference := []byte{0x00, 0x04, 0x00, 0x00, 0x00, 0x4d, 0x61, 0x72, 0x79}
	assert.Equal(t, reference, data)

	na = Age(24)
	data, err = ctx.Marshal(&na)
	assert.Nil(t, err)
	reference = []byte{0x01, 0x18, 0x00, 0x00, 0x00}
	assert.Equal(t, reference, data)
}
