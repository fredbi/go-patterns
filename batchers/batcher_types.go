package batchers

import (
	"sync"
)

type (
	baseExecutor[T any] struct {
		batchSize int
		mx        sync.Mutex
		count     uint64
		*options
	}

	Executor[T any] struct {
		*baseExecutor[T]
		batch    []T
		executor func([]T)
	}

	PointerExecutor[T any] struct {
		*baseExecutor[T]
		batch    []*T
		executor func([]*T)
	}

	Option func(*options)

	options struct {
		// TODO
	}
)

func defaultOptions() *options {
	return &options{}
}

func newBaseExecutor[T any](batchSize int, opts ...Option) *baseExecutor[T] {
	e := &baseExecutor[T]{
		batchSize: batchSize,
		options:   defaultOptions(),
	}

	for _, apply := range opts {
		apply(e.options)
	}

	return e
}

func NewExecutor[T any | *any](batchSize int, executor func([]T), opts ...Option) *Executor[T] {
	return &Executor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     executor,
	}
}

func NewPointerExecutor[T any](batchSize int, executor func([]*T), opts ...Option) *PointerExecutor[T] {
	return &PointerExecutor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     executor,
	}
}

func (e *Executor[T]) cloneBatch(in []T) []T {
	clone := make([]T, len(in))
	copy(clone, e.batch)

	return clone
}

func (e *PointerExecutor[T]) cloneBatch(in []*T) []*T {
	// when adopting slices of pointers, we shallow clone individual elements
	clone := make([]*T, len(in))
	for i, element := range in {
		copied := *element
		clone[i] = &copied
	}

	return clone
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

func (e *PointerExecutor[T]) Push(in *T) {
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

func (e *PointerExecutor[T]) Flush() {
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

func (e *PointerExecutor[T]) executeClone() {
	e.mx.Lock()
	defer e.mx.Unlock()

	if len(e.batch) == 0 {
		return
	}

	e.executor(e.cloneBatch(e.batch))
	e.count += uint64(len(e.batch))
	e.batch = e.batch[:0]
}
