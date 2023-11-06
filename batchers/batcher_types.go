package batchers

import (
	"sync"
)

type (
	// batch[T TypeConstraint] is a slice of elements used to define a processing batch.
	//
	// In particular, batch knows how to produce a clone of its elements.
	batch[T TypeConstraint] []T

	baseExecutor[T TypeConstraint] struct {
		batchSize int
		mx        sync.Mutex
		count     uint64
		*options
	}
)

// Clone the elements of the batch.
func (b batch[T]) Clone() batch[T] {
	clone := make(batch[T], len(b))
	copy(clone, b)

	return clone
}

// Len yields the number of elements in the batch.
func (b batch[T]) Len() int {
	return len(b)
}

// Empty returns an empty batch.
func (b batch[T]) Empty() batch[T] {
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
