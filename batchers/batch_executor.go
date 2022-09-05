package batchers

// Executor runs an executor function over a batch of T.
//
// The Executor handles the slicing of Push-ed input into batches of fixed size.
//
// When the stream of input is complete, a call to Flush() executes the last (possibly incomplete) batch.
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
	e.mx.Lock()
	defer e.mx.Unlock()

	e.executeClone()
}

func (e *Executor[T]) executeClone() {
	if e.batch.Len() == 0 {
		return
	}

	e.executor(e.batch.Clone())
	e.count += uint64(e.batch.Len())
	e.batch = e.batch.Empty()
}
