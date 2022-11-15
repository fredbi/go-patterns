package sorters

import (
	"golang.org/x/exp/constraints"
)

type (
	// Ordered defines all ordered types (see  https://go.dev/ref/spec)
	Ordered = constraints.Ordered

	Comparison[T any] func(a, b T) int
)
