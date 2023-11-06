package sorters

import (
	"golang.org/x/exp/constraints"
)

type (
	// Ordered defines all ordered types (see  https://go.dev/ref/spec)
	Ordered = constraints.Ordered

	// Comparison knows how to compare 2 values of the same type.
	//
	// Return values:
	//
	// if a == b: 0
	// if a < b: -1
	// if a > b: 1
	Comparison[T any] func(a, b T) int
)
