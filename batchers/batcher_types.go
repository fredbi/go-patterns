package batchers

import (
	"sync"
)

type (
	// TypeConstraint is any type.
	TypeConstraint interface {
		any
	}

	// Batch[T TypeConstraint] is a slice of elements used to define a processing batch.
	// A Batch knows how to produce a clone of its elements.
	Batch[T TypeConstraint] []T

	baseExecutor[T TypeConstraint] struct {
		batchSize int
		mx        sync.Mutex
		count     uint64
		*options
	}
)

// Clone the elements of the batch.
func (b Batch[T]) Clone() Batch[T] {
	clone := make(Batch[T], len(b))
	copy(clone, b)

	return clone
}

// Len yields the number of elements in the batch.
func (b Batch[T]) Len() int {
	return len(b)
}

// Empty returns an empty batch.
func (b Batch[T]) Empty() Batch[T] {
	return b[:0]
}

// Executed yields the count of batch executions.
func (e *baseExecutor[T]) Executed() uint64 {
	e.mx.Lock()
	defer e.mx.Unlock()

	return e.count
}

func newBaseExecutor[T TypeConstraint](batchSize int, opts ...Option) *baseExecutor[T] {
	e := &baseExecutor[T]{
		batchSize: batchSize,
		options:   defaultOptions(),
	}

	for _, apply := range opts {
		apply(e.options)
	}

	return e
}
