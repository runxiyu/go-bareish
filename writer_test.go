package bare

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteUint(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)

	err := w.WriteUint(0x7F)
	assert.Nil(t, err)
	err = w.WriteUint(0x1337)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x7F, 0xB7, 0x26}, b.Bytes())
}

func TestWriteU8(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteU8(uint8(0x42))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x42}, b.Bytes())
}

func TestWriteU16(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteU16(uint16(0xCAFE))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xFE, 0xCA}, b.Bytes())
}

func TestWriteU32(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteU32(uint32(0xDEADBEEF))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xEF, 0xBE, 0xAD, 0xDE}, b.Bytes())
}

func TestWriteU64(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteU64(uint64(0xCAFEBABEDEADBEEF))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xEF, 0xBE, 0xAD, 0xDE, 0xBE, 0xBA, 0xFE, 0xCA}, b.Bytes())
}

func TestWriteInt(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)

	err := w.WriteInt(42)
	assert.Nil(t, err)
	err = w.WriteInt(-1337)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x54, 0xf1, 0x14}, b.Bytes())
}

func TestWriteI8(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteI8(int8(-42))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xD6}, b.Bytes())
}

func TestWriteI16(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteI16(int16(-1234))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x2E, 0xFB}, b.Bytes())
}

func TestWriteI32(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteI32(int32(-12345678))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xB2, 0x9E, 0x43, 0xFF}, b.Bytes())
}

func TestWriteI64(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteI64(int64(-12345678987654321))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x4F, 0x0B, 0x6E, 0x9D, 0xAB, 0x23, 0xD4, 0xFF}, b.Bytes())
}

func TestWriteF32(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteF32(float32(1337.42))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x71, 0x2D, 0xA7, 0x44}, b.Bytes())
}

func TestWriteF64(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteF64(float64(133713371337.42424242))
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x9B, 0x6C, 0xC9, 0x20, 0xF0, 0x21, 0x3F, 0x42}, b.Bytes())
}

func TestWriteBool(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteBool(false)
	assert.Nil(t, err)
	err = w.WriteBool(true)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x00, 0x01}, b.Bytes())
}

func TestWriteString(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteString("こんにちは、世界！")
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x1B, 0x00, 0x00, 0x00, 0xE3, 0x81, 0x93, 0xE3,
		0x82, 0x93, 0xE3, 0x81, 0xAB, 0xE3, 0x81, 0xA1, 0xE3, 0x81, 0xAF, 0xE3,
		0x80, 0x81, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x81},
		b.Bytes())
}

func TestWriteDataFixed(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteDataFixed([]byte{0x13, 0x37, 0x42})
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x13, 0x37, 0x42}, b.Bytes())
}

func TestWriteData(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := w.WriteData([]byte{0x13, 0x37, 0x42})
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x03, 0x00, 0x00, 0x00, 0x13, 0x37, 0x42}, b.Bytes())
}
