package schema

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"git.sr.ht/~sircmpwn/go-bare"
)

func SchemaFor(val interface{}) (string, error) {
	t := reflect.TypeOf(val)
	return schemaFor(t)
}

func schemaFor(t reflect.Type) (string, error) {
	// TODO: Implement user-defined types for unparsing schemas from
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		return "", errors.New("Expected val to be pointer type")
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		schema, err := schemaFor(t)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("optional<%s>", schema), nil
	}

	switch t.Kind() {
	case reflect.Uint8:
		return "u8", nil
	case reflect.Uint16:
		return "u16", nil
	case reflect.Uint32:
		return "u32", nil
	case reflect.Uint64:
		return "u64", nil
	case reflect.Uint:
		return "u32", nil
	case reflect.Int8:
		return "i8", nil
	case reflect.Int16:
		return "i16", nil
	case reflect.Int32:
		return "i32", nil
	case reflect.Int64:
		return "i64", nil
	case reflect.Int:
		return "i32", nil
	case reflect.Float32:
		return "f32", nil
	case reflect.Float64:
		return "f64", nil
	case reflect.Bool:
		return "bool", nil
	case reflect.String:
		return "string", nil
	case reflect.Struct:
		return schemaForStruct(t)
	default:
		return "", &bare.UnsupportedTypeError{t}
	}
}

func schemaForStruct(t reflect.Type) (string, error) {
	buf := bytes.NewBufferString("{\n")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		schema, err := schemaFor(field.Type)
		if err != nil {
			return "", err
		}
		name := field.Name
		if tag, ok := field.Tag.Lookup("bare"); ok {
			name = tag
		}
		buf.WriteString(fmt.Sprintf("\t%s: %s\n", name, schema))
	}
	buf.WriteString("}")
	return buf.String(), nil
}
