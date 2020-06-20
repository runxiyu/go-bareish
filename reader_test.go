package bare

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadU8(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x42})
	r := NewReader(b)
	v, err := r.ReadU8()
	assert.Nil(t, err, "Expected ReadU8 to return without error")
	assert.Equal(t, v, uint8(0x42), "Expected reader to return 0x42")
	_, err = r.ReadU8()
	assert.Equal(t, err, io.EOF, "Expected ReadU8 to return EOF")
}

func TestReadU16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xFE, 0xCA})
	r := NewReader(b)
	v, err := r.ReadU16()
	assert.Nil(t, err, "Expected ReadU16 to return without error")
	assert.Equal(t, v, uint16(0xCAFE), "Expected reader to return 0xCAFE")
	_, err = r.ReadU16()
	assert.Equal(t, err, io.EOF, "Expected ReadU16 to return EOF")
}

func TestReadU32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xEF, 0xBE, 0xAD, 0xDE})
	r := NewReader(b)
	v, err := r.ReadU32()
	assert.Nil(t, err, "Expected ReadU32 to return without error")
	assert.Equal(t, v, uint32(0xDEADBEEF), "Expected reader to return 0xDEADBEEF")
	_, err = r.ReadU32()
	assert.Equal(t, err, io.EOF, "Expected ReadU32 to return EOF")
}

func TestReadU64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA})
	r := NewReader(b)
	v, err := r.ReadU64()
	assert.Nil(t, err, "Expected ReadU64 to return without error")
	assert.Equal(t, v, uint64(0xCAFEBABEDEADBEEF),
		"Expected reader to return 0xCAFEBABEDEADBEEF")
	_, err = r.ReadU64()
	assert.Equal(t, err, io.EOF, "Expected ReadU64 to return EOF")
}

func TestReadI8(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xD6})
	r := NewReader(b)
	v, err := r.ReadI8()
	assert.Nil(t, err, "Expected ReadI8 to return without error")
	assert.Equal(t, v, int8(-42), "Expected reader to return -42")
	_, err = r.ReadI8()
	assert.Equal(t, err, io.EOF, "Expected ReadI8 to return EOF")
}

func TestReadI16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x2E, 0xFB})
	r := NewReader(b)
	v, err := r.ReadI16()
	assert.Nil(t, err, "Expected ReadI16 to return without error")
	assert.Equal(t, v, int16(-1234), "Expected reader to return -1234")
	_, err = r.ReadI16()
	assert.Equal(t, err, io.EOF, "Expected ReadI16 to return EOF")
}

func TestReadI32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xB2, 0x9E, 0x43, 0xFF})
	r := NewReader(b)
	v, err := r.ReadI32()
	assert.Nil(t, err, "Expected ReadI32 to return without error")
	assert.Equal(t, v, int32(-12345678), "Expected reader to return -12345678")
	_, err = r.ReadI32()
	assert.Equal(t, err, io.EOF, "Expected ReadI32 to return EOF")
}

func TestReadI64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF})
	r := NewReader(b)
	v, err := r.ReadI64()
	assert.Nil(t, err, "Expected ReadI64 to return without error")
	assert.Equal(t, v, int64(-12345678987654321),
		"Expected reader to return -12345678987654321")
	_, err = r.ReadI64()
	assert.Equal(t, err, io.EOF, "Expected ReadI64 to return EOF")
}

func TestReadF32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x71, 0x2D, 0xA7, 0x44})
	r := NewReader(b)
	v, err := r.ReadF32()
	assert.Nil(t, err, "Expected ReadF32 to return without error")
	assert.Equal(t, v, float32(1337.42), "Expected reader to return 1337.42")
	_, err = r.ReadF32()
	assert.Equal(t, err, io.EOF, "Expected ReadF32 to return EOF")
}

func TestReadF64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42})
	r := NewReader(b)
	v, err := r.ReadF64()
	assert.Nil(t, err, "Expected ReadF64 to return without error")
	assert.Equal(t, v, float64(133713371337.42424242),
		"Expected reader to return 133713371337.42424242")
	_, err = r.ReadF64()
	assert.Equal(t, err, io.EOF, "Expected ReadF64 to return EOF")
}

func TestReadBool(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x00, 0x01, 0x02})
	r := NewReader(b)
	v, err := r.ReadBool()
	assert.Nil(t, err, "Expected ReadBool to return without error")
	assert.Equal(t, v, false, "Expected reader to return false")

	v, err = r.ReadBool()
	assert.Nil(t, err, "Expected ReadBool to return without error")
	assert.Equal(t, v, true, "Expected reader to return true")

	v, err = r.ReadBool()
	assert.Nil(t, err, "Expected ReadBool to return without error")
	assert.Equal(t, v, true, "Expected reader to return true")

	_, err = r.ReadBool()
	assert.Equal(t, err, io.EOF, "Expected ReadBool to return EOF")
}

func TestReadString(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x1B, 0x00, 0x00, 0x00, 0xE3, 0x81, 0x93, 0xE3,
		0x82, 0x93, 0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3,
		0x80, 0x81, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81})
	r := NewReader(b)
	v, err := r.ReadString()
	assert.Nil(t, err, "Expected ReadString to return without error")
	assert.Equal(t, v, "こんにちは、世界！",
		"Expected reader to return 'こんにちは、世界！'")
	_, err = r.ReadString()
	assert.Equal(t, err, io.EOF, "Expected ReadString to return EOF")
}

func TestReadDataFixed(t *testing.T) {
	ref := []byte{0x13, 0x37, 0x42}
	b := bytes.NewBuffer(ref)
	r := NewReader(b)
	buf := make([]byte, 3)
	err := r.ReadDataFixed(buf)
	assert.Nil(t, err, "Expected ReadDataFixed to return without error")
	assert.Equal(t, buf, ref, "Expected reader to return 0x13, 0x37, 0x42")
	err = r.ReadDataFixed(buf)
	assert.Equal(t, err, io.EOF, "Expected ReadDataFixed to return EOF")
}

func TestReadData(t *testing.T) {
	ref := []byte{0x03, 0x00, 0x00, 0x00, 0x13, 0x37, 0x42}
	b := bytes.NewBuffer(ref)
	r := NewReader(b)
	v, err := r.ReadData()
	assert.Nil(t, err, "Expected ReadData to return without error")
	assert.Equal(t, v, ref[4:], "Expected reader to return 0x13, 0x37, 0x42")
	_, err = r.ReadData()
	assert.Equal(t, err, io.EOF, "Expected ReadData to return EOF")
}
