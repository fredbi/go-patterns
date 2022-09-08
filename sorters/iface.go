package sorters

type (
	// Ordered defines all ordered types (see  https://go.dev/ref/spec)
	Ordered interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float64 | ~float32 | ~string | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}

	Comparison[T any] func(a, b T) int
)
