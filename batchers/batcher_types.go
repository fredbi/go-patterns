package batchers

import (
	"sync"
)

type (
	Executor[T any] struct {
		batchSize int
		executor  func([]T)
		batch     []T
		mx        sync.Mutex
		count     uint64
		*options
	}

	Option func(*options)

	options struct {
	}
)

func defaultOptions() *options {
	return &options{}
}

func NewExecutor[T any | *any](batchSize int, executor func([]T), opts ...Option) *Executor[T] {
	e := &Executor[T]{
		batchSize: batchSize,
		executor:  executor,
		options:   defaultOptions(),
	}

	for _, apply := range opts {
		apply(e.options)
	}

	return e
}

func (e *Executor[T]) cloneBatch(in []T) []T {
	clone := make([]T, len(in))
	copy(clone, e.batch)

	return clone
}

/*
func (e *Executor[T]) cloneElementsBatch(in []T) []T {
	clone := make([]T, len(in))
	for i, element := range in {
		copied := *element
		clone[i] = &copied
	}

	return clone
}
*/

func (e *Executor[T]) Push(in T) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.batch = append(e.batch, in)

	if len(e.batch) < e.batchSize {
		return
	}

	e.executeClone()
}

func (e *Executor[T]) Flush() {
	e.executeClone()
}

func (e *Executor[T]) executeClone() {
	e.mx.Lock()
	defer e.mx.Unlock()

	if len(e.batch) == 0 {
		return
	}

	e.executor(e.cloneBatch(e.batch))
	e.count += uint64(len(e.batch))
	e.batch = e.batch[:0]
}
