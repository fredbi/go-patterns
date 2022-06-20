package batchers

import (
	"sync"
)

type (
	TypeConstraint interface {
		any
	}

	Batch[T TypeConstraint]        []T
	BatchPointer[T TypeConstraint] []*T

	baseExecutor[T TypeConstraint] struct {
		batchSize int
		mx        sync.Mutex
		count     uint64
		*options
	}

	Option func(*options)

	options struct {
		// TODO
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

func (b BatchPointer[T]) Clone() BatchPointer[T] {
	// when adopting slices of pointers, we shallow clone individual elements
	clone := make(BatchPointer[T], len(b))
	for i, element := range b {
		copied := *element
		clone[i] = &copied
	}

	return clone
}

func (b BatchPointer[T]) Len() int {
	return len(b)
}

func (b BatchPointer[T]) Empty() BatchPointer[T] {
	return b[:0]
}

func defaultOptions() *options {
	return &options{}
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
