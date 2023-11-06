package batchers

var _ Batcher[*int] = &BatchPointerExecutor[int]{}

// BatchPointerExecutor runs an executor function over a batch of pointers *T.
//
// A shallow copy of *T is performed when assembling into a new batch.
//
// Apart from the pointer logic, the BatchPointerExecutor behaves like the Executor.
type BatchPointerExecutor[T TypeConstraint] struct {
	*baseExecutor[T]
	batch    batch[*T]
	executor func(batch[*T])
}

func NewBatchPointerExecutor[T TypeConstraint](batchSize int, executor PointerExecutor[T], opts ...Option) *BatchPointerExecutor[T] {
	return &BatchPointerExecutor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     asBatchSlicePtr(executor),
		batch:        make(batch[*T], 0, batchSize),
	}
}

func asBatchSlicePtr[T TypeConstraint](fn PointerExecutor[T]) func(batch[*T]) {
	return func(in batch[*T]) {
		fn([]*T(in))
	}
}

func (e *BatchPointerExecutor[T]) Push(in *T) {
	if in == nil {
		return // skip nil values
	}

	e.mx.Lock()
	defer e.mx.Unlock()

	// when adopting slices of pointers, we shallow-clone individual elements
	clone := *in

	e.batch = append(e.batch, &clone)

	if e.batch.Len() < e.batchSize {
		return
	}

	e.executeClone()
}

func (e *BatchPointerExecutor[T]) Flush() {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.executeClone()
}

func (e *BatchPointerExecutor[T]) executeClone() {
	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
