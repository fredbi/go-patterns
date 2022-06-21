package batchers

type PointerExecutor[T TypeConstraint] struct {
	*baseExecutor[T]
	batch    Batch[*T]
	executor func(Batch[*T])
}

func NewPointerExecutor[T TypeConstraint](batchSize int, executor func(Batch[*T]), opts ...Option) *PointerExecutor[T] {
	return &PointerExecutor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     executor,
		batch:        make(Batch[*T], 0, batchSize),
	}
}

func (e *PointerExecutor[T]) Push(in *T) {
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

func (e *PointerExecutor[T]) Flush() {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.executeClone()
}

func (e *PointerExecutor[T]) executeClone() {
	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
