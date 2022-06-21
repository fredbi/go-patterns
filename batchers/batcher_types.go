package batchers

import (
	"sync"
)

type (
	TypeConstraint interface {
		any
	}

	// Batch[T TypeConstraint] is a slice of elements used to define a processing batch.
	// A Batch knows how to produce a clone.
	Batch[T TypeConstraint] []T

	baseExecutor[T TypeConstraint] struct {
		batchSize int
		mx        sync.Mutex
		count     uint64
		*options
	}
)

func (b Batch[T]) Clone() Batch[T] {
	clone := make(Batch[T], len(b))
	copy(clone, b)

	return clone
}

func (b Batch[T]) Len() int {
	return len(b)
}

func (b Batch[T]) Empty() Batch[T] {
	return b[:0]
}

//nolint:revive // bug in current version of revive linter
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
