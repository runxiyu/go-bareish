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
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		i   int
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
		[]byte{0xEF, 0xBE, 0xAD, 0xDE},
		[]byte{0xD6},
		[]byte{0x2E, 0xFB},
		[]byte{0xB2, 0x9E, 0x43, 0xFF},
		[]byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF},
		[]byte{0xB2, 0x9E, 0x43, 0xFF},
		[]byte{0x71, 0x2D, 0xA7, 0x44},
		[]byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42},
		[]byte{0x00, 0x01, 0x02},
		[]byte{0x1B, 0x00, 0x00, 0x00, 0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93,
			0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3, 0x80,
			0x81, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81},
	}

	err = Unmarshal(payloads[0], &u8)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, u8, uint8(0x42), "Expected Unmarshal to read 0x42")

	err = Unmarshal(payloads[1], &u16)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, u16, uint16(0xCAFE), "Expected Unmarshal to read 0xCAFE")

	err = Unmarshal(payloads[2], &u32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, u32, uint32(0xDEADBEEF), "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal(payloads[3], &u64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, u64, uint64(0xCAFEBABEDEADBEEF), "Expected Unmarshal to read 0xCAFEBABEDEADBEEF")

	err = Unmarshal(payloads[4], &u)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, u, uint(0xDEADBEEF), "Expected Unmarshal to read 0xDEADBEEF")

	err = Unmarshal(payloads[5], &i8)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, i8, int8(-42), "Expected Unmarshal to read -42")

	err = Unmarshal(payloads[6], &i16)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, i16, int16(-1234), "Expected Unmarshal to read -1234")

	err = Unmarshal(payloads[7], &i32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, i32, int32(-12345678), "Expected Unmarshal to read -12345678")

	err = Unmarshal(payloads[8], &i64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, i64, int64(-12345678987654321), "Expected Unmarshal to read -12345678987654321")

	err = Unmarshal(payloads[9], &i)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, i, int(-12345678), "Expected Unmarshal to read -12345678")

	err = Unmarshal(payloads[10], &f32)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, f32, float32(1337.42), "Expected Unmarshal to read 1337.42")

	err = Unmarshal(payloads[11], &f64)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, f64, float64(133713371337.42424242), "Expected Unmarshal to read 133713371337.42424242")

	err = Unmarshal(payloads[12][0:], &b)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, b, false, "Expected Unmarshal to read false")

	err = Unmarshal(payloads[12][1:], &b)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, b, true, "Expected Unmarshal to read true")

	err = Unmarshal(payloads[12][2:], &b)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, b, true, "Expected Unmarshal to read true")

	err = Unmarshal(payloads[13], &str)
	assert.Nil(t, err, "Expected Unmarshal to return without error")
	assert.Equal(t, str, "こんにちは、世界！", "Expected Unmarshal to read 'こんにちは、世界！'")
}
