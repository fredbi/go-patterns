package iterators

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type dummyStruct struct {
	A int
	B string
}

func TestSliceIterator(t *testing.T) {
	t.Run("with happy path slice", func(t *testing.T) {
		slice := dummySlice()

		t.Run(fmt.Sprintf("should Collect %d items", len(slice)), func(t *testing.T) {
			iterator := NewSliceIterator(slice)
			items, err := iterator.Collect()
			require.NoError(t, err)

			require.Len(t, items, len(slice))
		})

		t.Run(fmt.Sprintf("should CollectPtr %d items", len(slice)), func(t *testing.T) {
			iterator := NewSliceIterator(slice)
			items, err := iterator.CollectPtr()
			require.NoError(t, err)

			require.Len(t, items, len(slice))
		})
	})

	t.Run("with empty iterator", func(t *testing.T) {
		iterator := NewSliceIterator([]dummyStruct{})

		require.False(t, iterator.Next())

		item, err := iterator.Item()
		require.ErrorIs(t, io.EOF, err)
		require.Empty(t, item)

		items, err := iterator.Collect()
		require.NoError(t, err)
		require.Empty(t, items)

		itemsPtr, err := iterator.CollectPtr()
		require.NoError(t, err)
		require.Empty(t, itemsPtr)

		require.NoError(t, iterator.Close())
	})

	t.Run("with nil iterator", func(t *testing.T) {
		iterator := NewSliceIterator[dummyStruct](nil)

		require.False(t, iterator.Next())
		require.NoError(t, iterator.Close())

		empty, err := iterator.Item()
		require.ErrorIs(t, io.EOF, err)
		require.Empty(t, empty)

		items, err := iterator.Collect()
		require.NoError(t, err)
		require.Empty(t, items)

		itemsPtr, err := iterator.CollectPtr()
		require.NoError(t, err)
		require.Empty(t, itemsPtr)

		require.NoError(t, iterator.Close())
	})

	t.Run("with out-of-sync call to Item()", func(t *testing.T) {
		t.Run("should error if Next() has never been called", func(t *testing.T) {
			iterator := NewSliceIterator[dummyStruct](dummySlice())
			empty, err := iterator.Item()
			require.ErrorIs(t, io.EOF, err)
			require.Empty(t, empty)

			require.NoError(t, iterator.Close())
		})

		t.Run("should error if Item() is called after iteration is complete", func(t *testing.T) {
			iterator := NewSliceIterator[dummyStruct](dummySlice())
			for iterator.Next() {
				_, err := iterator.Item()
				require.NoError(t, err)
			}

			empty, err := iterator.Item()
			require.ErrorIs(t, io.EOF, err)
			require.Empty(t, empty)

			require.NoError(t, iterator.Close())
		})
	})
}

func dummySlice() []dummyStruct {
	return []dummyStruct{
		{
			A: 1,
			B: "x",
		},
		{
			A: 2,
			B: "y",
		},
	}
}
