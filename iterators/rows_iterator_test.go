package iterators

import (
	"testing"

	"github.com/fredbi/go-patterns/iterators/internal/testdb"
	"github.com/stretchr/testify/require"
)

func TestRowsIterator(t *testing.T) {
	t.Run("with happy data path", func(t *testing.T) {
		dbName := testdb.UniqueDBName()
		db, err := testdb.CreateDBAndData(dbName)
		require.NoError(t, err)

		t.Run("with happy path cursor", func(t *testing.T) {
			t.Run("should Collect 2 items", func(t *testing.T) {
				rows, err := testdb.OpenDBCursor(db)
				require.NoError(t, err)

				iterator := NewSqlxIterator[testdb.DummyRow](rows)

				items, err := iterator.Collect()
				require.NoError(t, err)

				require.Len(t, items, 2)
			})

			t.Run("should CollectPtr 2 items", func(t *testing.T) {
				rows, err := testdb.OpenDBCursor(db)
				require.NoError(t, err)

				iterator := NewSqlxIterator[testdb.DummyRow](rows, WithRowsPreallocatedItems(10))
				items, err := iterator.CollectPtr()
				require.NoError(t, err)

				require.Len(t, items, 2)
				require.Equal(t, 10, cap(items))
			})
		})

		t.Run("with empty iterator", func(t *testing.T) {
			rows, err := testdb.EmptyDBCursor(db)
			require.NoError(t, err)

			iterator := NewSqlxIterator[testdb.DummyRow](rows)

			require.False(t, iterator.Next())

			item, err := iterator.Item()
			require.ErrorContains(t, err, "sql: Rows are closed")
			require.Empty(t, item)

			items, err := iterator.Collect()
			require.NoError(t, err)
			require.Empty(t, items)

			itemsPtr, err := iterator.CollectPtr()
			require.NoError(t, err)
			require.Empty(t, itemsPtr)

			require.NoError(t, iterator.Close())

			t.Run("should not error if closed twice", func(t *testing.T) {
				require.NoError(t, iterator.Close())
			})
		})
	})

	t.Run("with sql cursor error", func(t *testing.T) {
		dbName := testdb.UniqueDBName()
		db, err := testdb.CreateDBWithWrongData(dbName)
		require.NoError(t, err)

		t.Run("should iterate over 1 item then error", func(t *testing.T) {
			rows, err := testdb.OpenDBCursor(db)
			require.NoError(t, err)

			iterator := NewSqlxIterator[testdb.DummyRow](rows)

			count := 0
			for iterator.Next() {
				item, e := iterator.Item()
				if e != nil {
					break
				}
				require.NotEmpty(t, item)
				count++
			}
			require.Equal(t, 1, count)
			t.Logf("received expected error: %v", err)
		})

		t.Run("should collect 1 item then error", func(t *testing.T) {
			rows, err := testdb.OpenDBCursor(db)
			require.NoError(t, err)

			iterator := NewSqlxIterator[testdb.DummyRow](rows)

			items, err := iterator.Collect()
			require.Error(t, err)
			require.Len(t, items, 1)
		})

		t.Run("should collect pointers over 1 item then error", func(t *testing.T) {
			rows, err := testdb.OpenDBCursor(db)
			require.NoError(t, err)

			iterator := NewSqlxIterator[testdb.DummyRow](rows)

			items, err := iterator.CollectPtr()
			require.Error(t, err)
			require.Len(t, items, 1)
		})
	})

	t.Run("with close error", func(t *testing.T) {
		dbName := testdb.UniqueDBName()
		db, err := testdb.CreateDBAndData(dbName)
		require.NoError(t, err)

		t.Run("should error on close", func(t *testing.T) {
			rows, err := testdb.OpenDBCursor(db)
			require.NoError(t, err)

			iterator := NewSqlxIterator[testdb.DummyRow](rows)

			require.True(t, iterator.Next())
			require.NoError(t, rows.Close())
			_, err = iterator.Item()
			require.Error(t, err)
		})
	})
}
