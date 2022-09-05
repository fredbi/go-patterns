package iterators

type (
	// Iterator knows how to iterate over a collection.
	Iterator interface {
		Next() bool
		Close() error
	}

	// ScannableIterator is an iterator over DB records that can be scanned.
	ScannableIterator interface {
		Iterator
		StructScan(interface{}) error
	}

	// StructIterator is an iterator that delivers items of some type T.
	StructIterator[T any] interface {
		Iterator

		// Item return the current iterated item.
		//
		// Next() must have been called at least once.
		Item() (T, error)

		// Collect returns all items in one slice, then closes the iterator
		Collect() ([]T, error)

		// CollectPtr returns all items in one slice of pointers, then closes the iterator
		CollectPtr() ([]*T, error)
	}

	dummy struct{}
)
