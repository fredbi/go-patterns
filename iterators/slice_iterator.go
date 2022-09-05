package iterators

import (
	"io"
	"sync"
)

var _ StructIterator[dummy] = &SliceIterator[dummy]{}

// SliceIterator constructs an iterator based on a slice of items.
//
// This very simple iterator is essentially used for testing.
type SliceIterator[T any] struct {
	rows  []T
	index int
	mx    sync.RWMutex
}

// NewSliceIterator constructs a SliceIterator from a slice of items (rows).
func NewSliceIterator[T any](rows []T) *SliceIterator[T] {
	return &SliceIterator[T]{
		index: -1,
		rows:  rows,
	}
}

func (si *SliceIterator[T]) Close() error {
	return nil
}

func (si *SliceIterator[T]) Next() bool {
	si.mx.Lock()
	defer si.mx.Unlock()

	si.index++

	return si.index < len(si.rows)
}

func (si *SliceIterator[T]) Item() (T, error) {
	si.mx.RLock()
	defer si.mx.RUnlock()

	if si.index < 0 || si.index >= len(si.rows) {
		var empty T
		return empty, io.EOF
	}

	return si.rows[si.index], nil
}

func (si *SliceIterator[T]) Collect() ([]T, error) {
	return si.rows, nil
}

func (si *SliceIterator[T]) CollectPtr() ([]*T, error) {
	ptrs := make([]*T, 0, len(si.rows))

	for _, toPin := range si.rows {
		val := toPin
		ptrs = append(ptrs, &val)
	}

	return ptrs, nil
}
