package bare

import (
	"bytes"
	"errors"
	"reflect"
)

// Unmarshals a BARE message into val, which should be a pointer to a value of
// the message type.
//
// Go "int" and "uint" types are represented as BARE u32 and i32 types
// respectively, for message compatibility with both 32-bit and 64-bit systems.
func Unmarshal(data []byte, val interface{}) error {
	b := bytes.NewReader(data)
	r := NewReader(b)
	return UnmarshalReader(r, val)
}

// Unmarshals a BARE message into val, from a reader. val should be a pointer
// to a value of the message type. See Unmarshal for details.
func UnmarshalReader(r *Reader, val interface{}) error {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	} else {
		return errors.New("Expected val to be pointer type")
	}

	if t.Kind() == reflect.Ptr {
		// optional<type>
		s, err := r.ReadU8()
		if err != nil {
			return err
		}

		if s == 0 {
			v.Set(reflect.Zero(t))
			return nil
		} else {
			v.Set(reflect.New(t.Elem()))
			t = t.Elem()
			v = v.Elem()
		}
	}

	// TODO: 
	// - data, data<len>
	// - Decode maps, tagged unions

	var err error
	switch t.Kind() {
	case reflect.Uint8:
		var i uint8
		i, err = r.ReadU8()
		v.SetUint(uint64(i))
	case reflect.Uint16:
		var i uint16
		i, err = r.ReadU16()
		v.SetUint(uint64(i))
	case reflect.Uint32:
		var i uint32
		i, err = r.ReadU32()
		v.SetUint(uint64(i))
	case reflect.Uint64:
		var i uint64
		i, err = r.ReadU64()
		v.SetUint(uint64(i))
	case reflect.Uint:
		var i uint32
		i, err = r.ReadU32()
		v.SetUint(uint64(i))
	case reflect.Int8:
		var i int8
		i, err = r.ReadI8()
		v.SetInt(int64(i))
	case reflect.Int16:
		var i int16
		i, err = r.ReadI16()
		v.SetInt(int64(i))
	case reflect.Int32:
		var i int32
		i, err = r.ReadI32()
		v.SetInt(int64(i))
	case reflect.Int64:
		var i int64
		i, err = r.ReadI64()
		v.SetInt(int64(i))
	case reflect.Int:
		var i int32
		i, err = r.ReadI32()
		v.SetInt(int64(i))
	case reflect.Float32:
		var f float32
		f, err = r.ReadF32()
		v.SetFloat(float64(f))
	case reflect.Float64:
		var f float64
		f, err = r.ReadF64()
		v.SetFloat(float64(f))
	case reflect.Bool:
		var b bool
		b, err = r.ReadBool()
		v.SetBool(b)
	case reflect.String:
		var s string
		s, err = r.ReadString()
		v.SetString(s)
	case reflect.Array:
		return unmarshalArray(r, t, v)
	case reflect.Slice:
		return unmarshalSlice(r, t, v)
	case reflect.Struct:
		return unmarshalStruct(r, t, v)
	case reflect.Map:
		return unmarshalMap(r, t, v)
	default:
		return &UnsupportedTypeError{t}
	}
	return err
}

func unmarshalStruct(r *Reader, t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i)
		err := UnmarshalReader(r, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalArray(r *Reader, t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.Len(); i++ {
		value := v.Index(i)
		err := UnmarshalReader(r, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalSlice(r *Reader, t reflect.Type, v reflect.Value) error {
	l, err := r.ReadU32()
	if err != nil {
		return err
	}
	slice := reflect.MakeSlice(t, int(l), int(l))
	for i := 0; i < int(l); i++ {
		value := slice.Index(i)
		err := UnmarshalReader(r, value.Addr().Interface())
		if err != nil {
			return err
		}
	}
	v.Set(slice)
	return nil
}

func unmarshalMap(r *Reader, t reflect.Type, v reflect.Value) error {
	l, err := r.ReadU32()
	if err != nil {
		return err
	}
	m := reflect.MakeMapWithSize(t, int(l))
	for i := 0; i < int(l); i++ {
		key := reflect.New(t.Key())
		value := reflect.New(t.Elem())
		err := UnmarshalReader(r, key.Interface())
		if err != nil {
			return err
		}
		err = UnmarshalReader(r, value.Interface())
		if err != nil {
			return err
		}
		m.SetMapIndex(key.Elem(), value.Elem())
	}
	v.Set(m)
	return nil
}
