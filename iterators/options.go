package iterators

type (
	// RowsIteratorOption provides options to the RowsIterator
	RowsIteratorOption func(*rowsIteratorOptions)

	// ChanIteratorOption provides options to the ChanIterator
	ChanIteratorOption func(*chanIteratorOptions)

	rowsIteratorOptions struct {
		preallocatedItems int
	}

	chanIteratorOptions struct {
		*rowsIteratorOptions

		fanInBuffers int
	}
)

func rowsIteratorOptionsWithDefault(opts []RowsIteratorOption) *rowsIteratorOptions {
	options := &rowsIteratorOptions{
		preallocatedItems: 1000,
	}

	for _, apply := range opts {
		apply(options)
	}

	return options
}

// WithRowsPreallocatedItems preallocate n items in the returned slice when
// using the Collect and CollectPtr methods.
//
// The default value is 1000.
func WithRowsPreallocatedItems(n int) RowsIteratorOption {
	return func(o *rowsIteratorOptions) {
		o.preallocatedItems = n
	}
}

func chanIteratorOptionsWithDefault(opts []ChanIteratorOption) *chanIteratorOptions {
	options := &chanIteratorOptions{
		rowsIteratorOptions: rowsIteratorOptionsWithDefault(nil),
		fanInBuffers:        -1,
	}
	for _, apply := range opts {
		apply(options)
	}

	return options
}

// WithChanPreallocatedItems preallocate n items in the returned slice when
// using the Collect and CollectPtr methods.
func WithChanPreallocatedItems(n int) ChanIteratorOption {
	return func(o *chanIteratorOptions) {
		o.preallocatedItems = n
	}
}

// WithChanFanInBuffers allocates buffers to fan-in the input results.
//
// The default value is the number of underlying iterators.
func WithChanFanInBuffers(n int) ChanIteratorOption {
	return func(o *chanIteratorOptions) {
		o.fanInBuffers = n
	}
}
