package main

import (
	"fmt"
	"io"
	"strings"

	"git.sr.ht/~sircmpwn/go-bare/schema"
)

func genTypes(w io.Writer, types []schema.SchemaType) {
	fmt.Fprintf(w, "\nimport \"git.sr.ht/~sircmpwn/go-bare\"\n")

	for _, ty := range types {
		switch ty := ty.(type) {
		case *schema.UserDefinedType:
			genUserType(w, ty)
		case *schema.UserDefinedEnum:
			genUserEnum(w, ty)
		}
	}
}

func genUserType(w io.Writer, udt *schema.UserDefinedType) {
	fmt.Fprintf(w, "\ntype %s ", udt.Name())
	genType(w, udt.Type(), 0)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\nfunc (t *%s) Decode(data []byte) {", udt.Name())
	fmt.Fprintf(w, "\n\tbare.Unmarshal(data, t)")
	fmt.Fprintf(w, "\n}\n")
}

func genUserEnum(w io.Writer, ude *schema.UserDefinedEnum) {
	// TODO: Disambiguate between enums with conflicting value names
	fmt.Fprintf(w, "\ntype %s %s\n", ude.Name(), primitiveType(ude.Kind()))
	fmt.Fprintf(w, "\nconst (")
	for i, val := range ude.Values() {
		if i == 0 {
			fmt.Fprintf(w, "\n\t%s %s = %d", val.Name(), ude.Name(), val.Value())
		} else {
			fmt.Fprintf(w, "\n\t%s = %d", val.Name(), val.Value())
		}
	}
	fmt.Fprintf(w, "\n)\n")
}

func genType(w io.Writer, ty schema.Type, indent int) {
	switch ty := ty.(type) {
	case *schema.PrimitiveType:
		fmt.Fprintf(w, "%s", primitiveType(ty.Kind()))
	case *schema.DataType:
		if ty.Kind() == schema.DataArray {
			fmt.Fprintf(w, "[%d]byte", ty.Length())
		} else {
			fmt.Fprintf(w, "[]byte")
		}
	case *schema.StructType:
		maxName := 0
		for _, field := range ty.Fields() {
			if len(field.Name()) > maxName {
				maxName = len(field.Name())
			}
		}

		fmt.Fprintf(w, "struct {\n")
		for _, field := range ty.Fields() {
			genIndent(w, indent + 1)
			n := fieldName(field.Name())
			fmt.Fprintf(w, "%s ", n)
			for i := len(n); i < maxName; i++ {
				fmt.Fprintf(w, " ")
			}
			genType(w, field.Type(), indent + 1)
			fmt.Fprintf(w, " `bare:\"%s\"`", field.Name())
			fmt.Fprintf(w, "\n")
		}
		genIndent(w, indent)
		fmt.Fprintf(w, "}")
	case *schema.NamedUserType:
		fmt.Fprintf(w, "%s", ty.Name())
	case *schema.MapType:
		fmt.Fprintf(w, "map[")
		genType(w, ty.Key(), indent)
		fmt.Fprintf(w, "]")
		genType(w, ty.Value(), indent)
	case *schema.ArrayType:
		if ty.Kind() == schema.Array {
			fmt.Fprintf(w, "[%d]", ty.Length())
		} else {
			fmt.Fprintf(w, "[]")
		}
		genType(w, ty.Member(), indent)
	case *schema.OptionalType:
		fmt.Fprintf(w, "*")
		genType(w, ty.Subtype(), indent)
	default:
		panic(fmt.Errorf("TODO: %T", ty))
	}
}

func primitiveType(kind schema.TypeKind) string {
	switch kind {
	case schema.U8:
		return "uint8"
	case schema.U16:
		return "uint16"
	case schema.U32:
		return "uint32"
	case schema.U64:
		return "uint64"
	case schema.I8:
		return "int8"
	case schema.I16:
		return "int16"
	case schema.I32:
		return "int32"
	case schema.I64:
		return "int64"
	case schema.E8:
		return "uint8"
	case schema.E16:
		return "uint16"
	case schema.E32:
		return "uint32"
	case schema.E64:
		return "uint64"
	case schema.F32:
		return "float32"
	case schema.F64:
		return "float64"
	case schema.Bool:
		return "bool"
	case schema.String:
		return "string"
	}
	panic(fmt.Errorf("Invalid primitive type %d", kind))
}

func genIndent(w io.Writer, indent int) {
	for ; indent > 0; indent-- {
		fmt.Fprintf(w, "\t")
	}
}

func fieldName(n string) string {
	// TODO: Correct initialisms
	return strings.ToUpper(n[:1]) + n[1:]
}
