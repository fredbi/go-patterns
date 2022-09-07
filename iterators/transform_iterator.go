package iterators

import (
	"context"
)

const (
	ctxKeyIteration ctxTransformIterator = iota + 1
)

var _ StructIterator[dummy] = &TransformIterator[SqlxIterator[dummy], dummy]{}

type (
	ctxTransformIterator uint8

	// TransformerCtx converts a struct of type S into one of type T, with a context about the state of the iterator.
	TransformerCtx[S, T any] func(context.Context, S) (T, error)

	// TransformIterator transforms any iterator into an iterator that transforms the input of type S into type T
	// at every call to Item().
	TransformIterator[S, T any] struct {
		StructIterator[S]
		iterated    int
		ctx         context.Context
		transformer TransformerCtx[S, T]

		*rowsIteratorOptions
	}

	IteratorContext struct {
		Iterated int
	}
)

// GetIteratorContext allows the retrieval of the context of the iterator from within a transformer.
func GetIteratorContext(ctx context.Context) *IteratorContext {
	val, ok := ctx.Value(ctxKeyIteration).(*IteratorContext)
	if !ok {
		return nil
	}

	return val
}

// NewTransformIterator makes a StructIterator[T] from a ScannableIterator.
//
// The parent context provided allows the transformer to know about the current context of the iterator.
//
// This is useful if the transformation depends on the currently iterated step.
//
// Notice that the transformer may also perform some other things, e.g. logging, collecting some stats or traces.
func NewTransformIterator[S, T any](ctx context.Context, iterator StructIterator[S], transformer TransformerCtx[S, T], opts ...RowsIteratorOption) *TransformIterator[S, T] {
	return &TransformIterator[S, T]{
		StructIterator:      iterator,
		ctx:                 ctx,
		transformer:         transformer,
		rowsIteratorOptions: rowsIteratorOptionsWithDefault(opts),
	}
}

func (rt *TransformIterator[S, T]) iteratorContext() context.Context {
	return context.WithValue(rt.ctx, ctxKeyIteration, &IteratorContext{Iterated: rt.iterated})
}

func (rt *TransformIterator[S, T]) Next() bool {
	isNext := rt.StructIterator.Next()
	if isNext {
		rt.iterated++
	}

	return isNext
}

func (rt *TransformIterator[S, T]) Item() (T, error) {
	input, err := rt.StructIterator.Item()
	if err != nil {
		var empty T

		return empty, err
	}

	output, err := rt.transformer(rt.iteratorContext(), input)
	if err != nil {
		var empty T

		return empty, err
	}

	return output, nil
}

func (rt *TransformIterator[S, T]) Collect() ([]T, error) {
	return collectAndClose[T](rt, rt.preallocatedItems)
}

func (rt *TransformIterator[S, T]) CollectPtr() ([]*T, error) {
	return collectPtrAndClose[T](rt, rt.preallocatedItems)
}
