package bare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalValue(t *testing.T) {
	var (
		err error

		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
		u   uint
		vu  Uint
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		i   int
		vi  Int
		f32 float32
		f64 float64
		b   bool
		str string
	)

	payloads := [][]byte{
		[]byte{0x42},
		[]byte{0xFE, 0xCA},
		[]byte{0xEF, 0xBE, 0xAD, 0xDE},
		[]byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA},
		[]byte{0xEF, 0xFD, 0xB6, 0xF5, 0x0D},
		[]byte{0xD6},
		[]byte{0x2E, 0xFB},
		[]byte{0xB2, 0x9E, 0x43, 0xFF},
		[]byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF},
		[]byte{0x9B, 0x85, 0xE3, 0x0B},
		[]byte{0x71, 0x2D, 0xA7, 0x44},
		[]byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42},
		[]byte{0x00, 0x01, 0x02},
		[]byte{0x1B, 0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93, 0xE3, 0x81, 0xAB,
			0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3, 0x80, 0x81, 0xE4, 0xB8,
			0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81},
	}

	err = Unmarshal(payloads[0], &u8)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint8(0x42), u8, "Expected Unmarshal to read 0x42")

	err = Unmarshal(payloads[1], &u16)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint16(0xCAFE), u16, "Expected Unmarshal to read 0xCAFE")

	err = Unmarshal(payloads[2], &u32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint32(0xDEADBEEF), u32, "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal(payloads[3], &u64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint64(0xCAFEBABEDEADBEEF), u64, "Expected Unmarshal to read 0xCAFEBABEDEADBEEF")

	err = Unmarshal(payloads[4], &u)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint(0xDEADBEEF), u, "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal(payloads[4], &vu)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, Uint(0xDEADBEEF), vu, "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal(payloads[5], &i8)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, int8(-42), i8, "Expected Unmarshal to read -42")

	err = Unmarshal(payloads[6], &i16)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, int16(-1234), i16, "Expected Unmarshal to read -1234")

	err = Unmarshal(payloads[7], &i32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, int32(-12345678), i32, "Expected Unmarshal to read -12345678")

	err = Unmarshal(payloads[8], &i64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, int64(-12345678987654321), i64, "Expected Unmarshal to read -12345678987654321")

	err = Unmarshal(payloads[9], &i)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, int(-12345678), i, "Expected Unmarshal to read -12345678")

	err = Unmarshal(payloads[9], &vi)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, Int(-12345678), vi, "Expected Unmarshal to read -12345678")

	err = Unmarshal(payloads[10], &f32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, float32(1337.42), f32, "Expected Unmarshal to read 1337.42")

	err = Unmarshal(payloads[11], &f64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, float64(133713371337.42424242), f64, "Expected Unmarshal to read 133713371337.42424242")

	err = Unmarshal(payloads[12][0:], &b)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.False(t, b, "Expected Unmarshal to read false")

	err = Unmarshal(payloads[12][1:], &b)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.True(t, b, "Expected Unmarshal to read true")

	err = Unmarshal(payloads[12][2:], &b)
	assert.EqualError(t, err, "Invalid bool value: 0x2")

	err = Unmarshal(payloads[13], &str)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, "こんにちは、世界！", str, "Expected Unmarshal to read 'こんにちは、世界！'")
}

func TestUnmarshalOptional(t *testing.T) {
	var val *uint32
	err := Unmarshal([]byte{0x00}, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Nil(t, val, "Expected Unmarshal to read nil")

	err = Unmarshal([]byte{0x01, 0xEF, 0xBE, 0xAD, 0xDE}, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.NotNil(t, val, "Expected Unmarshal to read non-nil value")
	assert.Equal(t, uint32(0xDEADBEEF), *val, "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal([]byte{0x02}, &val)
	assert.EqualError(t, err, "Invalid optional value: 0x2")
}

func TestUnmarshalStruct(t *testing.T) {
	type Coordinates struct{ X, Y, Z uint }
	payload := []byte{0x01, 0x02, 0x03}
	var coords Coordinates
	err := Unmarshal(payload, &coords)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint(1), coords.X, "Expected Unmarshal to read {1, 2, 3}")
	assert.Equal(t, uint(2), coords.Y, "Expected Unmarshal to read {1, 2, 3}")
	assert.Equal(t, uint(3), coords.Z, "Expected Unmarshal to read {1, 2, 3}")
}

func TestUnmarshalOmittedFields(t *testing.T) {
	type Coordinates struct {
		X uint
		Y uint
		Z uint `bare:"-"`
	}
	payload := []byte{0x01, 0x02}
	var coords Coordinates
	err := Unmarshal(payload, &coords)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint(1), coords.X, "Expected Unmarshal to read {1, 2}")
	assert.Equal(t, uint(2), coords.Y, "Expected Unmarshal to read {1, 2}")
	assert.Equal(t, uint(0), coords.Z, "Expected Unmarshal to ignore field")
}

func TestUnmarshalArray(t *testing.T) {
	var val [4]uint8
	err := Unmarshal([]byte{0x11, 0x22, 0x33, 0x44}, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, uint8(0x11), val[0], "Expected Unmarshal to read 0x11")
	assert.Equal(t, uint8(0x22), val[1], "Expected Unmarshal to read 0x22")
	assert.Equal(t, uint8(0x33), val[2], "Expected Unmarshal to read 0x33")
	assert.Equal(t, uint8(0x44), val[3], "Expected Unmarshal to read 0x44")
}

func TestUnmarshalSlice(t *testing.T) {
	var val []uint8
	err := Unmarshal([]byte{0x04, 0x11, 0x22, 0x33, 0x44}, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, 4, len(val), "Expected Unmarshal to read slice of length 4")
	assert.Equal(t, uint8(0x11), val[0], "Expected Unmarshal to read 0x11")
	assert.Equal(t, uint8(0x22), val[1], "Expected Unmarshal to read 0x22")
	assert.Equal(t, uint8(0x33), val[2], "Expected Unmarshal to read 0x33")
	assert.Equal(t, uint8(0x44), val[3], "Expected Unmarshal to read 0x44")

	MaxArrayLength(64)
	err = Unmarshal([]byte{100}, &val)
	assert.EqualError(t, err, "Array length 100 exceeds configured limit of 64")
}

func TestUnmarshalMap(t *testing.T) {
	var val map[uint8]uint8
	payload := []byte{
		0x03,
		0x01, 0x11,
		0x02, 0x22,
		0x03, 0x33,
	}
	err := Unmarshal(payload, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, 3, len(val), "Expected Unmarshal to read map of length 3")
	assert.Equal(t, uint8(0x11), val[uint8(0x01)], "Expected Unmarshal to read 0x01 -> 0x11")
	assert.Equal(t, uint8(0x22), val[uint8(0x02)], "Expected Unmarshal to read 0x02 -> 0x22")
	assert.Equal(t, uint8(0x33), val[uint8(0x03)], "Expected Unmarshal to read 0x03 -> 0x33")

	t.Run("rejects duplicate keys", func(t *testing.T) {
		payload = []byte{
			0x02,
			0x01, 0x13,
			0x01, 0x37,
		}
		err = Unmarshal(payload, &val)
		assert.EqualError(t, err, "Encountered duplicate map key: 1")
	})

	t.Run("respects size limits", func(t *testing.T) {
		MaxMapSize(64)
		payload = []byte{100}
		err = Unmarshal(payload, &val)
		assert.EqualError(t, err, "Map size 100 exceeds configured limit of 64")
	})
}

func TestUnmarshalUnion(t *testing.T) {
	var val NameAge
	payload := []byte{0x00, 0x04, 0x4d, 0x61, 0x72, 0x79}
	err := Unmarshal(payload, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")

	name, ok := val.(*Name)
	assert.True(t, ok)
	assert.Equal(t, Name("Mary"), *name)

	payload = []byte{0x01, 0x30}
	err = Unmarshal(payload, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")

	age, ok := val.(*Age)
	assert.True(t, ok)
	assert.Equal(t, Age(24), *age)

	payload = []byte{0x13, 0x37}
	err = Unmarshal(payload, &val)
	assert.EqualError(t, err, "Invalid union tag 19 for type NameAge")
}

func TestUnmarshalCustom(t *testing.T) {
	var val Custom
	payload := []byte{0x0, 0x42}

	err := Unmarshal(payload, &val)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, Custom(0x42), val)
}
