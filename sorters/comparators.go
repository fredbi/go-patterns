package sorters

import (
	"bytes"
	"strings"

	"golang.org/x/text/collate"
)

// Reverse the comparison order.
func Reverse[T any](in Comparison[T]) Comparison[T] {
	return func(a, b T) int {
		return -1 * in(a, b)
	}
}

// ReversePtr reverses the nil comparison logic for pointer operands: !nil < nil instead of nil < !nil.
//
// The comparison logic for non-nil values is not altered.
func ReversePtr[T any](in Comparison[*T]) Comparison[*T] {
	return func(a, b *T) int {
		if r, ok := comparePointers(a, b); ok {
			return -1 * r
		}

		return in(a, b)
	}
}

// Ptr makes a comparison over pointers from a comparison over types.
func Ptr[T any](in Comparison[T]) Comparison[*T] {
	return func(a, b *T) int {
		if r, ok := comparePointers(a, b); ok {
			return r
		}

		return in(*a, *b)
	}
}

// OrderedComparator builds a comparator for any ordered golang type (i.e. numerical types).
func OrderedComparator[T Ordered]() Comparison[T] {
	return func(a, b T) int {
		if a == b {
			return 0
		}

		if a < b {
			return -1
		}

		return 1
	}
}

// OrderedPtrComparator builds a comparator for pointers to any ordered golang type (i.e. numerical types).
//
// Pointer logic is that nil < !nil , nil == nil. Usual comparison occurs whenever
// both arguments are not nil.
func OrderedPtrComparator[T Ordered]() Comparison[*T] {
	return Ptr(OrderedComparator[T]())
}

// BytesComparator builds a comparator for []byte slices, with some collating sequence option.
func BytesComparator(opts ...CompareOption) Comparison[[]byte] {
	options := compareOptionsWithDefault(opts)

	if options.Locale == nil {
		return bytes.Compare
	}

	collator := collate.New(*options.Locale)

	return func(a, b []byte) int {
		return collator.Compare(a, b)
	}
}

// StringsComparator builds a comparator for strings, with some collating sequence option.
func StringsComparator(opts ...CompareOption) Comparison[string] {
	options := compareOptionsWithDefault(opts)

	if options.Locale == nil {
		return strings.Compare
	}

	collator := collate.New(*options.Locale)

	return func(a, b string) int {
		return collator.CompareString(a, b)
	}
}

// StringsPtrComparator builds a comparator for pointers to strings, with some collating sequence option.
//
// Pointer logic is that nil < !nil , nil == nil. Usual comparison occurs whenever
// both arguments are not nil.
func StringsPtrComparator(opts ...CompareOption) Comparison[*string] {
	return Ptr(StringsComparator(opts...))
}

// BoolsComparator compare booleans: false < true.
func BoolsComparator(_ ...CompareOption) Comparison[bool] {
	return func(a, b bool) int {
		if a == b {
			return 0
		}

		if !a {
			return -1
		}

		return 1
	}
}

// BoolsPtrComparator compare pointers to booleans: false < true.
//
// Pointer logic is that nil < !nil , nil == nil. Usual comparison occurs whenever
// both arguments are not nil.
func BoolsPtrComparator(opts ...CompareOption) Comparison[*bool] {
	return Ptr(BoolsComparator(opts...))
}

// comparePointers applies the pointer comparison logic for nil values.
//
// It returns true if some decision is made regarding nil.
// When false, the values are both non-nil and the integer result shall be discarded.
func comparePointers[T any](a, b *T) (int, bool) {
	if a == nil && b == nil {
		return 0, true
	}
	if a == nil && b != nil {
		return -1, true
	}
	if a != nil && b == nil {
		return 1, true
	}

	return 0, false
}
