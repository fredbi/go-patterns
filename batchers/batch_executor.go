package batchers

var _ Batcher[int] = &BatchExecutor[int]{}

// BatchExecutor runs an executor function over a slice []T.
//
// The Executor handles the slicing of Push-ed input into batches of fixed size.
//
// When the stream of input is complete, a call to Flush() executes the last (possibly incomplete) batch.
type BatchExecutor[T TypeConstraint] struct {
	*baseExecutor[T]
	batch         batch[T]
	executor      func(batch[T])
	wantsPointers bool
	shallowClone  func(T) (bool, T)
}

// NewBatchExecutor builds a batch executor that slices pushed items in groups of batchSize elements.
func NewBatchExecutor[T TypeConstraint](batchSize int, executor Executor[T], opts ...Option) *BatchExecutor[T] {
	e := &BatchExecutor[T]{
		baseExecutor:  newBaseExecutor[T](batchSize, opts...),
		executor:      asBatchSlice[T](executor),
		batch:         make(batch[T], 0, batchSize),
		wantsPointers: isPointer[T](),
	}

	if e.wantsPointers {
		// enables element cloning logics whenever T is a pointer type
		e.shallowClone = clone[T]()
	}

	return e
}

func asBatchSlice[T TypeConstraint](fn Executor[T]) func(batch[T]) {
	return func(in batch[T]) {
		fn([]T(in))
	}
}

func (e *BatchExecutor[T]) Push(in T) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.wantsPointers {
		// when adopting slices of pointers, we shallow-clone the elements that are pushed
		//
		// NOTE(fredbi): this indulge into a bit of reflection. The previously adopted alternative
		// was to expose a specific BatchPointerExecutor[T] operating explicitly on []*T slices.
		//
		// The adopted trade-off is eventually to simplify the use of this package, at the cost of
		// some reflection (there is no type assertion possible on parameterized types).
		skip, clone := e.shallowClone(in)
		// skip nil values
		if skip {
			return
		}

		in = clone
	}

	e.batch = append(e.batch, in)

	if e.batch.Len() < e.batchSize {
		return
	}

	e.executeClone()
}

func (e *BatchExecutor[T]) Flush() {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.executeClone()
}

func (e *BatchExecutor[T]) executeClone() {
	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
