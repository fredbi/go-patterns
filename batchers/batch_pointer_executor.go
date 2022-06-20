package batchers

type PointerExecutor[T TypeConstraint] struct {
	*baseExecutor[T]
	batch    BatchPointer[T]
	executor func(BatchPointer[T])
}

func NewPointerExecutor[T TypeConstraint](batchSize int, executor func(BatchPointer[T]), opts ...Option) *PointerExecutor[T] {
	return &PointerExecutor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     executor,
		batch:        make(BatchPointer[T], 0, batchSize),
	}
}

func (e *PointerExecutor[T]) Push(in *T) {
	if in == nil {
		return // skip nil values
	}

	e.mx.Lock()
	defer e.mx.Unlock()

	e.batch = append(e.batch, in)

	if e.batch.Len() < e.batchSize {
		return
	}

	e.executeClone()
}

func (e *PointerExecutor[T]) Flush() {
	e.executeClone()
}

func (e *PointerExecutor[T]) executeClone() {
	e.mx.Lock()
	defer e.mx.Unlock()

	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
