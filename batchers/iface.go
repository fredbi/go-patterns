package batchers

type (
	// TypeConstraint is any type.
	TypeConstraint interface {
		any
	}

	// Batcher knows how to push new elements to be processed.
	Batcher[T TypeConstraint] interface {
		Flush()
		Push(T)
	}

	// Executor knows how to apply some processing on on a slice of items of type T.
	Executor[T TypeConstraint] func([]T)

	// PointerExecutor knows how to apply some processing on a slice of pointers of type *T.
	PointerExecutor[T TypeConstraint] func([]*T)
)
