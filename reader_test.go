package bare

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadUint(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x7F, 0xB7, 0x26})
	r := NewReader(b)

	v, err := r.ReadUint()
	assert.Nil(t, err)
	assert.Equal(t, uint64(0x7F), v)

	v, err = r.ReadUint()
	assert.Nil(t, err)
	assert.Equal(t, uint64(0x1337), v)

	_, err = r.ReadUint()
	assert.Equal(t, err, io.EOF)
}

func TestReadU8(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x42})
	r := NewReader(b)
	v, err := r.ReadU8()
	assert.Nil(t, err)
	assert.Equal(t, v, uint8(0x42))
	_, err = r.ReadU8()
	assert.Equal(t, err, io.EOF)
}

func TestReadU16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xFE, 0xCA})
	r := NewReader(b)
	v, err := r.ReadU16()
	assert.Nil(t, err)
	assert.Equal(t, v, uint16(0xCAFE))
	_, err = r.ReadU16()
	assert.Equal(t, err, io.EOF)
}

func TestReadU32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xEF, 0xBE, 0xAD, 0xDE})
	r := NewReader(b)
	v, err := r.ReadU32()
	assert.Nil(t, err)
	assert.Equal(t, v, uint32(0xDEADBEEF))
	_, err = r.ReadU32()
	assert.Equal(t, err, io.EOF)
}

func TestReadU64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA})
	r := NewReader(b)
	v, err := r.ReadU64()
	assert.Nil(t, err)
	assert.Equal(t, v, uint64(0xCAFEBABEDEADBEEF))
	_, err = r.ReadU64()
	assert.Equal(t, err, io.EOF)
}

func TestReadInt(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x54, 0xf1, 0x14})
	r := NewReader(b)

	v, err := r.ReadInt()
	assert.Nil(t, err)
	assert.Equal(t, int64(42), v)

	v, err = r.ReadInt()
	assert.Nil(t, err)
	assert.Equal(t, v, int64(-1337))

	_, err = r.ReadUint()
	assert.Equal(t, err, io.EOF)
}

func TestReadI8(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xD6})
	r := NewReader(b)
	v, err := r.ReadI8()
	assert.Nil(t, err)
	assert.Equal(t, v, int8(-42))
	_, err = r.ReadI8()
	assert.Equal(t, err, io.EOF)
}

func TestReadI16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x2E, 0xFB})
	r := NewReader(b)
	v, err := r.ReadI16()
	assert.Nil(t, err)
	assert.Equal(t, v, int16(-1234))
	_, err = r.ReadI16()
	assert.Equal(t, err, io.EOF)
}

func TestReadI32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xB2, 0x9E, 0x43, 0xFF})
	r := NewReader(b)
	v, err := r.ReadI32()
	assert.Nil(t, err)
	assert.Equal(t, v, int32(-12345678))
	_, err = r.ReadI32()
	assert.Equal(t, err, io.EOF)
}

func TestReadI64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF})
	r := NewReader(b)
	v, err := r.ReadI64()
	assert.Nil(t, err)
	assert.Equal(t, v, int64(-12345678987654321))
	_, err = r.ReadI64()
	assert.Equal(t, err, io.EOF)
}

func TestReadF32(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x71, 0x2D, 0xA7, 0x44})
	r := NewReader(b)
	v, err := r.ReadF32()
	assert.Nil(t, err)
	assert.Equal(t, v, float32(1337.42))
	_, err = r.ReadF32()
	assert.Equal(t, err, io.EOF)
}

func TestReadF64(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42})
	r := NewReader(b)
	v, err := r.ReadF64()
	assert.Nil(t, err)
	assert.Equal(t, v, float64(133713371337.42424242))
	_, err = r.ReadF64()
	assert.Equal(t, err, io.EOF)
}

func TestReadBool(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x00, 0x01, 0x02})
	r := NewReader(b)
	v, err := r.ReadBool()
	assert.Nil(t, err)
	assert.Equal(t, v, false)

	v, err = r.ReadBool()
	assert.Nil(t, err)
	assert.Equal(t, v, true)

	_, err = r.ReadBool()
	assert.EqualError(t, err, "Invalid bool value: 0x2")

	_, err = r.ReadBool()
	assert.Equal(t, err, io.EOF)
}

func TestReadString(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x1B, 0xE3, 0x81, 0x93, 0xE3, 0x82, 0x93, 0xE3,
		0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3, 0x80, 0x81, 0xE4,
		0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81})

	b.Write([]byte{0x3, 66, 250, 67}) // Invalid sequence

	r := NewReader(b)
	v, err := r.ReadString()
	assert.Nil(t, err)
	assert.Equal(t, v, "こんにちは、世界！")

	_, err = r.ReadString()
	assert.EqualError(t, err, ErrInvalidStr.Error())

	_, err = r.ReadString()
	assert.Equal(t, err, io.EOF)
}

func TestReadDataFixed(t *testing.T) {
	ref := []byte{0x13, 0x37, 0x42}
	b := bytes.NewBuffer(ref)
	r := NewReader(b)
	buf := make([]byte, 3)
	err := r.ReadDataFixed(buf)
	assert.Nil(t, err)
	assert.Equal(t, buf, ref)
	err = r.ReadDataFixed(buf)
	assert.Equal(t, err, io.EOF)
}

func TestReadData(t *testing.T) {
	ref := []byte{0x03, 0x13, 0x37, 0x42}
	b := bytes.NewBuffer(ref)
	r := NewReader(b)
	v, err := r.ReadData()
	assert.Nil(t, err)
	assert.Equal(t, v, ref[1:])
	_, err = r.ReadData()
	assert.Equal(t, err, io.EOF)
}
