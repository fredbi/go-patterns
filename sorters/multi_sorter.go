package sorters

import (
	"sort"
)

var _ sort.Interface = &Multi[dummy]{}

type (
	Multi[T any] struct {
		collection       []T
		criteria         []Comparison[T]
		compoundCriteria Comparison[T]
	}

	dummy struct{}
)

// Len implements the sort.Interface
func (s Multi[T]) Len() int {
	return len(s.collection)
}

// Swap implements the sort.Interface
func (s *Multi[T]) Swap(i, j int) {
	s.collection[i], s.collection[j] = s.collection[j], s.collection[i]
}

// Less implements the sort.Interface
func (s *Multi[T]) Less(i, j int) bool {
	return s.compoundCriteria(s.collection[i], s.collection[j]) < 0
}

// Sort the inner collection. This is a shorthand for sort.Sort(s).
func (s *Multi[T]) Sort() {
	sort.Sort(s)
}

// Collection yields the inner collection.
func (s Multi[T]) Collection() []T {
	return s.collection
}

// NewMulti produces a sortable object, that supports multiple sorting criteria.
func NewMulti[T any](collection []T, criteria ...Comparison[T] /* sort options */) *Multi[T] {
	return &Multi[T]{
		collection:       collection,
		criteria:         criteria,
		compoundCriteria: compoundCriteria[T](criteria),
	}
}

func equal[T any](_, _ T) int { return 0 }

// CompoundCriteria compose multiple sorting criteria into a single Comparison[T].
func CompoundCriteria[T any](criteria ...Comparison[T]) Comparison[T] {
	return compoundCriteria(criteria)
}

func compoundCriteria[T any](criteria []Comparison[T]) Comparison[T] {
	if len(criteria) == 0 {
		return equal[T]
	}

	comparator := criteria[0]
	for _, extra := range criteria[1:] {
		extraComparator := extra
		pinnedComparator := comparator
		composed := func(a, b T) int {
			r := pinnedComparator(a, b)

			if r == 0 {
				return extraComparator(a, b)
			}

			return r
		}

		comparator = composed
	}

	return comparator
}
