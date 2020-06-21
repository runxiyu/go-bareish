package bare

import (
	"encoding/binary"
	"io"
)

// A Writer for BARE primitive types.
type Writer struct {
	base io.Writer
}

// Returns a new BARE primitive writer wrapping the given io.Writer.
func NewWriter(base io.Writer) *Writer {
	return &Writer{base}
}

func (w *Writer) WriteU8(i uint8) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteU16(i uint16) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteU32(i uint32) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteU64(i uint64) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteI8(i int8) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteI16(i int16) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteI32(i int32) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteI64(i int64) error {
	return binary.Write(w.base, binary.LittleEndian, i)
}

func (w *Writer) WriteF32(f float32) error {
	return binary.Write(w.base, binary.LittleEndian, f)
}

func (w *Writer) WriteF64(f float64) error {
	return binary.Write(w.base, binary.LittleEndian, f)
}

func (w *Writer) WriteBool(b bool) error {
	return binary.Write(w.base, binary.LittleEndian, b)
}

func (w *Writer) WriteE8(e uint8) error {
	return w.WriteU8(e)
}

func (w *Writer) WriteE16(e uint16) error {
	return w.WriteU16(e)
}

func (w *Writer) WriteE32(e uint32) error {
	return w.WriteU32(e)
}

func (w *Writer) WriteE64(e uint64) error {
	return w.WriteU64(e)
}

func (w *Writer) WriteString(str string) error {
	return w.WriteData([]byte(str))
}

// Writes a fixed amount of arbitrary data, defined by the length of the slice.
func (w *Writer) WriteDataFixed(data []byte) error {
	var amt int = 0
	for amt < len(data) {
		n, err := w.base.Write(data[amt:])
		if err != nil {
			return err
		}
		amt += n
	}
	return nil
}

// Writes arbitrary data whose length is encoded into the message.
func (w *Writer) WriteData(data []byte) error {
	err := w.WriteU32(uint32(len(data)))
	if err != nil {
		return err
	}
	var amt int = 0
	for amt < len(data) {
		n, err := w.base.Write(data[amt:])
		if err != nil {
			return err
		}
		amt += n
	}
	return nil
}
