package bare

import (
	"bytes"
	"errors"
	"reflect"
)

// Marshals a value (val, which must be a pointer) into a BARE message.
//
// Go "int" and "uint" types are represented as BARE u32 and i32 types
// respectively, for message compatibility with both 32-bit and 64-bit systems.
func Marshal(val interface{}) ([]byte, error) {
	ctx := NewContext()
	return ctx.Marshal(val)
}

// Marshals a value (val, which must be a pointer) into a BARE message.
//
// Go "int" and "uint" types are represented as BARE u32 and i32 types
// respectively, for message compatibility with both 32-bit and 64-bit systems.
func (ctx *Context) Marshal(val interface{}) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b)
	err := ctx.MarshalWriter(w, val)
	return b.Bytes(), err
}

// Marshals a value (val, which must be a pointer) into a BARE message and
// writes it to a Writer. See Marshal for details.
func MarshalWriter(w *Writer, val interface{}) error {
	ctx := NewContext()
	return ctx.MarshalWriter(w, val)
}

// Marshals a value (val, which must be a pointer) into a BARE message and
// writes it to a Writer. See Marshal for details.
func (ctx *Context) MarshalWriter(w *Writer, val interface{}) error {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	} else {
		return errors.New("Expected val to be pointer type, was " +
			t.Kind().String())
	}

	return ctx.marshalWriter(w, t, v)
}

func (ctx *Context) marshalWriter(w *Writer,
	t reflect.Type, v reflect.Value) error {
	if t.Kind() == reflect.Ptr {
		// optional<type>
		var set uint8
		if !v.IsZero() {
			t = t.Elem()
			v = v.Elem()
			set = 1
		}
		err := w.WriteU8(set)
		if err != nil {
			return err
		}
		if set == 0 {
			return nil
		}
	}

	if union, ok := v.Interface().(Union); ok {
		err := w.WriteU8(union.UnionTag())
		if err != nil {
			return err
		}
		v = reflect.ValueOf(union)
		t = v.Type()
	}

	// TODO: custom encoders
	switch t.Kind() {
	case reflect.Uint8:
		return w.WriteU8(uint8(v.Uint()))
	case reflect.Uint16:
		return w.WriteU16(uint16(v.Uint()))
	case reflect.Uint32:
		return w.WriteU32(uint32(v.Uint()))
	case reflect.Uint64:
		return w.WriteU64(uint64(v.Uint()))
	case reflect.Uint:
		return w.WriteU32(uint32(v.Uint()))
	case reflect.Int8:
		return w.WriteI8(int8(v.Int()))
	case reflect.Int16:
		return w.WriteI16(int16(v.Int()))
	case reflect.Int32:
		return w.WriteI32(int32(v.Int()))
	case reflect.Int64:
		return w.WriteI64(int64(v.Int()))
	case reflect.Int:
		return w.WriteI32(int32(v.Int()))
	case reflect.Float32:
		return w.WriteF32(float32(v.Float()))
	case reflect.Float64:
		return w.WriteF64(float64(v.Float()))
	case reflect.Bool:
		return w.WriteBool(v.Bool())
	case reflect.String:
		return w.WriteString(v.String())
	case reflect.Array:
		return ctx.marshalArray(w, t, v)
	case reflect.Slice:
		return ctx.marshalSlice(w, t, v)
	case reflect.Struct:
		return ctx.marshalStruct(w, t, v)
	case reflect.Map:
		return ctx.marshalMap(w, t, v)
	}

	return &UnsupportedTypeError{t}
}

func (ctx *Context) marshalStruct(w *Writer,
	t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i)
		err := ctx.MarshalWriter(w, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) marshalArray(w *Writer,
	t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.Len(); i++ {
		value := v.Index(i)
		err := ctx.MarshalWriter(w, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) marshalSlice(w *Writer,
	t reflect.Type, v reflect.Value) error {
	err := w.WriteU32(uint32(v.Len()))
	if err != nil {
		return err
	}
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i)
		err := ctx.MarshalWriter(w, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) marshalMap(w *Writer,
	t reflect.Type, v reflect.Value) error {
	err := w.WriteU32(uint32(v.Len()))
	if err != nil {
		return err
	}
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		err := ctx.marshalWriter(w, key.Type(), key)
		if err != nil {
			return err
		}
		err = ctx.marshalWriter(w, value.Type(), value)
		if err != nil {
			return err
		}
	}
	return nil
}
