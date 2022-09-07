package iterators

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChanIterator(t *testing.T) {
	t.Run("should Collect 4 items", func(t *testing.T) {
		baseIterators := []StructIterator[dummyStruct]{
			NewSliceIterator[dummyStruct](dummySlice()),
			NewSliceIterator[dummyStruct](dummySlice()),
		}

		iterator := NewChanIterator[dummyStruct](context.Background(), baseIterators, WithChanFanInBuffers(4))
		items, err := iterator.Collect()
		require.NoError(t, err)
		require.Len(t, items, 4)

		t.Run("should Collect again, with 0 items", func(t *testing.T) {
			iterator := NewChanIterator[dummyStruct](context.Background(), baseIterators)
			items, err := iterator.Collect()
			require.NoError(t, err)
			require.Len(t, items, 0)
		})
	})

	t.Run("should Collect 4 items (unbuffered)", func(t *testing.T) {
		baseIterators := []StructIterator[dummyStruct]{
			NewSliceIterator[dummyStruct](dummySlice()),
			NewSliceIterator[dummyStruct](dummySlice()),
		}

		iterator := NewChanIterator[dummyStruct](context.Background(), baseIterators,
			WithChanFanInBuffers(0),
			WithChanPreallocatedItems(10),
		)
		items, err := iterator.Collect()
		require.NoError(t, err)
		require.Len(t, items, 4)
		require.Equal(t, 10, cap(items))
	})

	t.Run("should CollectPtr 4 items", func(t *testing.T) {
		baseIterators := []StructIterator[dummyStruct]{
			NewSliceIterator[dummyStruct](dummySlice()),
			NewSliceIterator[dummyStruct](dummySlice()),
		}
		iterator := NewChanIterator[dummyStruct](context.Background(), baseIterators)
		itemsPtr, err := iterator.CollectPtr()
		require.NoError(t, err)
		require.Len(t, itemsPtr, 4)
	})

	t.Run("Items() should stop on cancelled context", func(t *testing.T) {
		baseIterators := []StructIterator[dummyStruct]{
			NewSliceIterator[dummyStruct](dummySlice()),
			NewSliceIterator[dummyStruct](dummySlice()),
		}
		ctx, cancel := context.WithCancel(context.Background())
		iterator := NewChanIterator[dummyStruct](ctx, baseIterators, WithChanFanInBuffers(0))

		require.True(t, iterator.Next())
		cancel()
		_, err := iterator.Item()
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Next() should stop on cancelled context", func(t *testing.T) {
		baseIterators := []StructIterator[dummyStruct]{
			NewSliceIterator[dummyStruct](dummySlice()),
			NewSliceIterator[dummyStruct](dummySlice()),
		}
		ctx, cancel := context.WithCancel(context.Background())
		iterator := NewChanIterator[dummyStruct](ctx, baseIterators, WithChanFanInBuffers(0))

		cancel()
		require.False(t, iterator.Next())
	})

	errTest := errors.New("test error")
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

	t.Run("Collect should error after 1 iteration", func(t *testing.T) {
		baseIterator := NewSliceIterator[dummyStruct](dummySlice())
		errorIterator := NewTransformIterator[dummyStruct, dummyStruct](context.Background(), baseIterator, errorer)

		iterator := NewChanIterator[dummyStruct](context.Background(),
			[]StructIterator[dummyStruct]{errorIterator},
			WithChanFanInBuffers(0),
		)

		_, err := iterator.Collect()
		require.ErrorIs(t, err, errTest)
	})

	t.Run("CollectPtr should error after 1 iteration", func(t *testing.T) {
		baseIterator := NewSliceIterator[dummyStruct](dummySlice())
		errorIterator := NewTransformIterator[dummyStruct, dummyStruct](context.Background(), baseIterator, errorer)

		iterator := NewChanIterator[dummyStruct](context.Background(),
			[]StructIterator[dummyStruct]{errorIterator},
			WithChanFanInBuffers(0),
		)

		_, err := iterator.CollectPtr()
		require.ErrorIs(t, err, errTest)
	})
}

func TestPreferErrorOverContext(t *testing.T) {
	errTest := errors.New("test error")
	require.ErrorIs(t, preferErrorOverContext(context.Canceled, errTest), errTest)
	require.ErrorIs(t, preferErrorOverContext(errTest, context.Canceled), errTest)
	require.ErrorIs(t, preferErrorOverContext(errTest, errors.New("another error")), errTest)
}
