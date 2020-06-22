package main

import (
	"fmt"
	"io"
	"strings"

	"git.sr.ht/~sircmpwn/go-bare/schema"
)

type Context struct {
	unions       []*schema.UserDefinedType
	unionMembers map[schema.Type]interface{}
}

func genTypes(w io.Writer, types []schema.SchemaType) {
	fmt.Fprintf(w, `
// THIS FILE WAS GENERATED BY A TOOL, DO NOT EDIT

import (
	"errors"
	"git.sr.ht/~sircmpwn/go-bare"
)
`)

	ctx := Context{
		unions:       nil,
		unionMembers: make(map[schema.Type]interface{}),
	}
	for _, ty := range types {
		switch ty := ty.(type) {
		case *schema.UserDefinedType:
			ctx.genUserType(w, ty)
		case *schema.UserDefinedEnum:
			ctx.genUserEnum(w, ty)
		}
	}

	if len(ctx.unions) > 0 {
		fmt.Fprintf(w, "\nfunc init() {\n")
		for _, udt := range ctx.unions {
			fmt.Fprintf(w, "\tbare.RegisterUnion((*%s)(nil)).\n", udt.Name())
			ut, _ := udt.Type().(*schema.UnionType)
			for i, ty := range ut.Types() {
				tag := ty.Tag()
				switch ty := ty.Type().(type) {
				case *schema.NamedUserType:
					fmt.Fprintf(w, "\t\tMember(*new(%s), %d)", ty.Name(), tag)
				default:
					panic(fmt.Errorf("TODO: Implement unions with primitive types"))
				}
				if i < len(ut.Types()) - 1 {
					fmt.Fprintf(w, ".\n")
				}
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "}\n")
	}
}

func (ctx *Context) genUserType(w io.Writer, udt *schema.UserDefinedType) {
	if udt.Type().Kind() == schema.Union {
		ctx.genUserUnion(w, udt)
		return
	}

	fmt.Fprintf(w, "\ntype %s ", udt.Name())
	genType(w, udt.Type(), 0)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\nfunc (t *%s) Decode(data []byte) error {", udt.Name())
	fmt.Fprintf(w, "\n\treturn bare.Unmarshal(data, t)")
	fmt.Fprintf(w, "\n}\n")

	fmt.Fprintf(w, "\nfunc (t *%s) Encode() ([]byte, error) {", udt.Name())
	fmt.Fprintf(w, "\n\treturn bare.Marshal(t)")
	fmt.Fprintf(w, "\n}\n")
}

func (ctx *Context) genUserEnum(w io.Writer, ude *schema.UserDefinedEnum) {
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

	fmt.Fprintf(w, "\nfunc (t %s) String() string {", ude.Name())
	fmt.Fprintf(w, "\n\tswitch (t) {")
	for _, val := range ude.Values() {
		fmt.Fprintf(w, "\n\tcase %s:", val.Name())
		fmt.Fprintf(w, "\n\t\treturn \"%s\"", val.Name())
	}
	fmt.Fprintf(w, "\n\t}")
	fmt.Fprintf(w, "\n\tpanic(errors.New(\"Invalid %s value\"))", ude.Name())
	fmt.Fprintf(w, "\n}\n")
}

func (ctx *Context) genUserUnion(w io.Writer, udt *schema.UserDefinedType) {
	fmt.Fprintf(w, "\ntype %s interface {", udt.Name())
	fmt.Fprintf(w, "\n\tbare.Union")
	fmt.Fprintf(w, "\n}\n")

	ut, _ := udt.Type().(*schema.UnionType)
	for _, ty := range ut.Types() {
		// XXX: This doesn't actually work the way it looks like it ought to
		if _, ok := ctx.unionMembers[ty.Type()]; ok {
			continue
		}

		ctx.unionMembers[ty.Type()] = nil

		switch ty := ty.Type().(type) {
		case *schema.NamedUserType:
			fmt.Fprintf(w, "\nfunc (_ %s) IsUnion() { }\n", ty.Name())
		default:
			panic(fmt.Errorf("TODO: Implement unions with primitive types"))
		}
	}

	ctx.unions = append(ctx.unions, udt)
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
		panic(fmt.Errorf("Unimplemented schema type: %T", ty))
	}
}

func genUnion(w io.Writer, ut *schema.UnionType, indent int) {
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
