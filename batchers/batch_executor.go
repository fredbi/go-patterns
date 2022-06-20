package batchers

type Executor[T TypeConstraint] struct {
	*baseExecutor[T]
	batch    Batch[T]
	executor func(Batch[T])
}

func NewExecutor[T TypeConstraint](batchSize int, executor func(Batch[T]), opts ...Option) *Executor[T] {
	return &Executor[T]{
		baseExecutor: newBaseExecutor[T](batchSize, opts...),
		executor:     executor,
		batch:        make(Batch[T], 0, batchSize),
	}
}

func (e *Executor[T]) Push(in T) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.batch = append(e.batch, in)

	if e.batch.Len() < e.batchSize {
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

	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
