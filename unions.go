package bare

import (
	"fmt"
	"reflect"
)

// Any type which is a union member must implement this interface. You must
// also call RegisterUnion for go-bare to marshal or unmarshal messages which
// utilize your union type.
type Union interface {
	IsUnion()
}

var unionRegistry map[reflect.Type][]reflect.Type

func init() {
	unionRegistry = make(map[reflect.Type][]reflect.Type)
}

// Registers a union type in this context. Pass the union interface and the
// list of types associated with it, sorted ascending by their union tag.
func RegisterUnion(iface interface{}, types ...interface{}) {
	ity := reflect.TypeOf(iface).Elem()
	if _, ok := unionRegistry[ity]; ok {
		return // Already registered
	}

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
	unionRegistry[ity] = utypes
}
