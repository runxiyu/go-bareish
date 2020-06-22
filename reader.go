package bare

import (
	"encoding/binary"
	"io"
)

// A Reader for BARE primitive types.
type Reader struct {
	base io.Reader
}

// Returns a new BARE primitive reader wrapping the given io.Reader.
func NewReader(base io.Reader) *Reader {
	return &Reader{base}
}

func (r *Reader) ReadU8() (uint8, error) {
	var i uint8
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadU16() (uint16, error) {
	var i uint16
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadU32() (uint32, error) {
	var i uint32
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadU64() (uint64, error) {
	var i uint64
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadI8() (int8, error) {
	var i int8
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadI16() (int16, error) {
	var i int16
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadI32() (int32, error) {
	var i int32
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadI64() (int64, error) {
	var i int64
	err := binary.Read(r.base, binary.LittleEndian, &i)
	return i, err
}

func (r *Reader) ReadF32() (float32, error) {
	var f float32
	err := binary.Read(r.base, binary.LittleEndian, &f)
	return f, err
}

func (r *Reader) ReadF64() (float64, error) {
	var f float64
	err := binary.Read(r.base, binary.LittleEndian, &f)
	return f, err
}

func (r *Reader) ReadBool() (bool, error) {
	var b bool
	err := binary.Read(r.base, binary.LittleEndian, &b)
	return b, err
}

func (r *Reader) ReadString() (string, error) {
	buf, err := r.ReadData()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Reads a fixed amount of arbitrary data, defined by the length of the slice.
func (r *Reader) ReadDataFixed(dest []byte) error {
	var amt int = 0
	for amt < len(dest) {
		n, err := r.base.Read(dest[amt:])
		if err != nil {
			return err
		}
		amt += n
	}
	return nil
}

// Reads arbitrary data whose length is read from the message.
func (r *Reader) ReadData() ([]byte, error) {
	l, err := r.ReadU32()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, l)
	var amt uint32 = 0
	for amt < l {
		n, err := r.base.Read(buf[amt:])
		if err != nil {
			return nil, err
		}
		amt += uint32(n)
	}
	return buf, nil
}
