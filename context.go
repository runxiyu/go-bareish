package bare

import (
	"fmt"
	"reflect"
)

// Any type which is a union member must implement this interface
type Union interface {
	UnionTag() uint8
}

// A marshal context, which stores an internal register of union types.
type Context struct {
	unions map[reflect.Type][]reflect.Type
}

// Creates a new marshal context. The use of a context is only required if you
// are working with a schema that uses union types.
func NewContext() *Context {
	return &Context{
		unions: make(map[reflect.Type][]reflect.Type),
	}
}

// Registers a union type in this context. Pass the union interface and the
// list of types associated with it, sorted ascending by their union tag.
func (ctx *Context) RegisterUnion(iface interface{}, types ...interface{}) {
	ity := reflect.TypeOf(iface).Elem()
	if !ity.Implements(reflect.TypeOf((*Union)(nil)).Elem()) {
		panic(fmt.Errorf("Type %s does not implement bare.Union", ity.Name()))
	}

	var utypes []reflect.Type
	for _, t := range types {
		ty := reflect.TypeOf(t)
		if !ty.AssignableTo(ity) {
			panic(fmt.Errorf("Type %s does not implement interface %s",
				ty.Name(), ity.Name()))
		}
		utypes = append(utypes, ty)
	}
	ctx.unions[ity] = utypes
}
