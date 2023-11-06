package batchers

import (
	"log"
	"reflect"
	//"unsafe"
)

// isPointer determines if the generic type T is a pointer type.
func isPointer[T any]() bool {
	var x T

	return reflect.ValueOf(x).Kind() == reflect.Pointer
}

// find U such as *U = T

// clone operates a shallow clone of a pointer value, when type T is a pointer type.
func clone[T any]() func(T) (bool, T) {
	var x T
	typ := reflect.TypeOf(x)
	log.Printf("DEBUG: type: %v\n", typ) // e.g. *int

	return func(in T) (bool, T) {
		val := reflect.ValueOf(in)
		if val.IsNil() {
			return true, in // nil value to be skipped
		}

		indirectedVal := reflect.Indirect(val)
		underlyingType := indirectedVal.Type()
		clone := indirectedVal.Interface()
		var z struct {
			x T
		}
		log.Printf("DEBUG: underlying: %v\n", underlyingType)                                                // e.g. int
		log.Printf("DEBUG: clone: [%T][%p][%p]: %v|%v\n", clone, &clone, &shallowClone, clone, shallowClone) // e.g. int

		clonePtr := reflect.NewAt(underlyingType, reflect.ValueOf(shallowClone).Addr().UnsafePointer())

		cast := clonePtr.Interface()

		return false, cast.(T)
	}
}
