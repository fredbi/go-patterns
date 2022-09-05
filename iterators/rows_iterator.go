package iterators

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

var _ StructIterator[dummy] = &RowsIterator[*sqlx.Rows, dummy]{}

type (
	// RowsIterator transforms a ScannableIterator of type R (e.g. a DB cursor such as sqlx.Rows)
	// into a StructIterator with target type T.
	//
	// Rows iteratated over R are scanned into structs of type T.
	//
	// Notice that the rows iterator is not goroutine-safe and should not be iterated concurrently.
	RowsIterator[R ScannableIterator, T any] struct {
		rows     R
		mx       sync.Mutex
		isClosed bool

		*rowsIteratorOptions
	}

	// SqlxIterator is a shorthand for RowsIterator[*sqlx.Rows, T].
	SqlxIterator[T any] struct {
		*RowsIterator[*sqlx.Rows, T]
	}
)

// NewSqlxIterator makes a SqlxIterator[T] producing items of type T from a github.com/jmoiron/sqlx.Rows cursor.
func NewSqlxIterator[T any](rows *sqlx.Rows, opts ...RowsIteratorOption) *SqlxIterator[T] {
	return &SqlxIterator[T]{
		RowsIterator: NewRowsIterator[*sqlx.Rows, T](rows, opts...),
	}
}

// NewRowsIterator makes a StructIterator[T] from a ScannableIterator.
func NewRowsIterator[R ScannableIterator, T any](rows R, opts ...RowsIteratorOption) *RowsIterator[R, T] {
	return &RowsIterator[R, T]{
		rows:                rows,
		rowsIteratorOptions: rowsIteratorOptionsWithDefault(opts),
	}
}

func (ri *RowsIterator[R, T]) Close() error {
	ri.mx.Lock()
	defer ri.mx.Unlock()

	if ri.isClosed {
		return nil
	}
	ri.isClosed = true

	return ri.rows.Close()
}

func (ri *RowsIterator[R, T]) Next() bool {
	return ri.rows.Next()
}

func (ri *RowsIterator[R, T]) Item() (T, error) {
	var data T

	if err := ri.rows.StructScan(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (ri *RowsIterator[R, T]) Collect() ([]T, error) {
	collection := make([]T, 0, ri.preallocatedItems)

	for ri.rows.Next() {
		item, err := ri.Item()
		if err != nil {
			_ = ri.Close()

			return collection, err
		}

		collection = append(collection, item)
	}

	if err := ri.Close(); err != nil {
		return collection, err
	}

	return collection, nil
}

func (ri *RowsIterator[R, T]) CollectPtr() ([]*T, error) {
	collection := make([]*T, 0, ri.preallocatedItems)

	for ri.rows.Next() {
		item, err := ri.Item()
		if err != nil {
			_ = ri.Close()

			return collection, err
		}

		collection = append(collection, &item)
	}

	if err := ri.Close(); err != nil {
		return collection, err
	}

	return collection, nil
}
