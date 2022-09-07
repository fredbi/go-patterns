package iterators

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type outStruct struct {
	X int
	Y string
}

func TestTransformIterator(t *testing.T) {
	transformer := func(ctx context.Context, in dummyStruct) (outStruct, error) {
		ictx := GetIteratorContext(ctx)
		var index int

		// this retrieves the current iterated count from the context
		if ictx != nil {
			index = ictx.Iterated
		}

		return outStruct{
			X: in.A + index,
			Y: strings.Repeat(in.B, 2),
		}, nil
	}

	t.Run("should Collect and transform 2 items", func(t *testing.T) {
		baseIterator := NewSliceIterator[dummyStruct](dummySlice())
		iterator := NewTransformIterator[dummyStruct, outStruct](context.Background(), baseIterator, transformer)

		items, err := iterator.Collect()
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("should CollectPtr and transform 2 items", func(t *testing.T) {
		baseIterator := NewSliceIterator[dummyStruct](dummySlice())
		iterator := NewTransformIterator[dummyStruct, outStruct](context.Background(), baseIterator, transformer)

		items, err := iterator.CollectPtr()
		require.NoError(t, err)
		require.Len(t, items, 2)
	})

	t.Run("should error after 1 iteration", func(t *testing.T) {
		errTest := errors.New("test error")
		baseIterator := NewSliceIterator[dummyStruct](dummySlice())
		errorer := func(ctx context.Context, in dummyStruct) (dummyStruct, error) {
			ictx := GetIteratorContext(ctx)
			var index int

			// this retrieves the current iterated count from the context
			if ictx != nil {
				index = ictx.Iterated
			}
			if index > 1 {
				return dummyStruct{}, errTest
			}
			return in, nil
		}

		errorIterator := NewTransformIterator[dummyStruct, dummyStruct](context.Background(), baseIterator, errorer)

		iterator := NewTransformIterator[dummyStruct, outStruct](context.Background(), errorIterator, transformer)

		items, err := iterator.Collect()
		require.ErrorIs(t, err, errTest)
		require.Len(t, items, 1)
	})
}

func TestGetIteratorContext(t *testing.T) {
	ctx := context.Background()

	// does not panic, just return nil
	require.Nil(t, GetIteratorContext(ctx))
}
