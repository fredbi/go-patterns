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
		withInnerClone bool
	}
)

func defaultOptions() *options {
	return &options{}
}

func NewExecutor[T any](batchSize int, executor func([]T), opts ...Option) *Executor[T] {
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

	e.executor(e.cloneBatch())
	e.count += uint64(len(e.batch))
	e.batch = e.batch[:0]
}

func (e *Executor[T]) cloneBatch() []T {
	clone := make([]T, len(e.batch))

	/*
		// TODO: do we have to clone elements? (option)
		if e.options.withInnerClone {
			for i, element := range e.batch {
				copied := *element
				clone[i] = &copied
			}

		} else {
			copy(clone, e.batch)
		}
	*/
	copy(clone, e.batch)

	return clone
}
