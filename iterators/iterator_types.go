package iterators

type (
	Iterator interface {
		Next() bool
		Close() error
	}

	ScannableIterator interface {
		Iterator
		StructScan(interface{}) error
	}

	StructIterator[T any] interface {
		Iterator
		Item() (T, error)
		Collect() ([]T, error)
		CollectPtr() ([]*T, error)
	}

	ScannableOption func(*scannableOptions)
	ChannelOption   func(*channelOptions)
)

type (
	scannableOptions struct {
		preallocatedItems int
	}

	channelOptions struct {
		buffers     int
		maxParallel int
	}

	sliceIterator[T any] struct {
		rows  []T
		index int
	}

	rowsIterator[T ScannableIterator] struct {
		rows T
	}
)

func NewSliceIterator[T any](rows []T) StructIterator[T] {
	return &sliceIterator[T]{
		rows: rows,
	}
}

func (si *sliceIterator[T]) Close() error {
	return nil
}

func (si *sliceIterator[T]) Next() bool {
	si.index++

	return si.index >= len(si.rows)
}

func (si *sliceIterator[T]) Item() (T, error) {
	return si.rows[si.index], nil
}

func (si *sliceIterator[T]) Collect() ([]T, error) {
	return si.rows, nil
}

func (si *sliceIterator[T]) CollectPtr() ([]*T, error) {
	ptrs := make([]*T, 0, len(si.rows))
	for _, toPin := range si.rows {
		val := toPin
		ptrs = append(ptrs, &val)
	}

	return ptrs, nil
}

func NewScannableIterator[T ScannableIterator](rows T, opts ...ScannableOption) StructIterator[T] {
	return &rowsIterator[T]{
		rows: rows,
	}
}

func NewFanInIterator[T any](iterators []StructIterator[T], opts ...ChannelOption) StructIterator[T] {
	return nil // TODO
}

func (ri *rowsIterator[T]) Close() error {
	return ri.rows.Close()
}

func (ri *rowsIterator[T]) Next() bool {
	return ri.rows.Next()
}

func (ri *rowsIterator[T]) Item() (T, error) {
	var data T

	if err := ri.rows.StructScan(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (ri *rowsIterator[T]) Collect() ([]T, error) {
	collection := make([]T, 0, 10) // TODO preallocate

	for ri.rows.Next() {
		item, err := ri.Item()
		if err != nil {
			_ = ri.Close()

			return nil, err
		}

		collection = append(collection, item)
	}

	if err := ri.Close(); err != nil {
		return nil, err
	}

	return collection, nil
}

func (ri *rowsIterator[T]) CollectPtr() ([]*T, error) {
	collection := make([]*T, 0, 10) // TODO preallocate

	for ri.rows.Next() {
		item, err := ri.Item()
		if err != nil {
			_ = ri.Close()

			return nil, err
		}

		collection = append(collection, &item)
	}

	if err := ri.Close(); err != nil {
		return nil, err
	}

	return collection, nil
}
