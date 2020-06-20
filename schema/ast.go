package schema

type SchemaTypeKind int

type SchemaType interface {
	Name() string
}

type UserDefinedType struct {
	name string
	type_ Type
}

func (udt *UserDefinedType) Name() string {
	return udt.name
}

func (udt *UserDefinedType) Type() Type {
	return udt.type_
}

type UserDefinedEnum struct {
	name string
	// TODO: members & values
}

func (ude *UserDefinedEnum) Name() string {
	return ude.name
}

type TypeKind int

const (
	U8 TypeKind = iota
	U16
	U32
	U64
	I8
	I16
	I32
	I64
	F32
	F64
	Bool
	E8
	E16
	E32
	E64
	String
	// data
	Data
	// data<length>
	DataFixed
	// optional<type>
	Optional
	// [nmemb]type
	DataArray
	// []type
	DataSlice
	// map[type]type
	Map
	// (type | type | ...)
	TaggedUnion
	// { fields... }
	Struct
	// Named user type
	UserType
)

type Type interface {
	Kind() TypeKind
}

type PrimitiveType struct {
	kind TypeKind
}

func (pt *PrimitiveType) Kind() TypeKind {
	return pt.kind
}

type OptionalType struct {
	subtype Type
}

func (ot *OptionalType) Kind() TypeKind {
	return Optional
}

func (ot *OptionalType) Subtype() Type {
	return ot.subtype
}

type DataType struct {
	length uint
}

func (dt *DataType) Kind() TypeKind {
	if dt.length == 0 {
		return DataSlice
	}
	return DataArray
}

func (dt *DataType) Length() uint {
	return dt.length
}
