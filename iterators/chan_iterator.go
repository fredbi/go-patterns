package iterators

import (
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"sync"
)

var _ StructIterator[dummy] = &ChanIterator[dummy]{}

// ChanIterator is a channel-based iterator that may be used to run a collection of StructIterators in parallel.
//
// Notice that its asynchronous working does not make it suitable to collect ordered items.
//
// The input is collected from a collection of input StructIterators, then may be collected by one or several goroutines
// reading from the ChanIterator using Next() and Item().
//
// If iterating in parallel from several goroutines, the WithChanFanOutBuffers option must be used to instruct the iterator to
// use at least as many buffers as there are readers and avoid blocking due to a starved channel.
//
// WithChanFanInBuffers may be used to pre-fetch from input iterators asynchronously.
//
// Methods Collect() and CollectPrt() can't be used by concurrent goroutines and are protected against such a misuse.
type ChanIterator[T any] struct {
	fanIn       chan T
	fanOut      chan T
	current     *T
	workerGroup *errgroup.Group
	ctx         context.Context
	mx          sync.Mutex

	*chanIteratorOptions
}

// NewChanIterator builds a ChanIterator and starts the goroutines pumping items from the input iterators.
//
// All goroutines are terminated and input iterators closed if the context is cancelled.
func NewChanIterator[T any](ctx context.Context, iterators []StructIterator[T], opts ...ChanIteratorOption) *ChanIterator[T] {
	var pendingWorkers sync.WaitGroup // rendez-vous to close the fan-in channel
	workerGroup, groupCtx := errgroup.WithContext(ctx)

	iter := &ChanIterator[T]{
		ctx:                 groupCtx,
		workerGroup:         workerGroup,
		chanIteratorOptions: chanIteratorOptionsWithDefault(opts),
	}

	iter.fanIn = make(chan T, iter.fanInBuffers)
	iter.fanOut = make(chan T, iter.fanOutBuffers)

	for i := range iterators {
		idx := i
		pendingWorkers.Add(1)

		workerGroup.Go(func() error {
			iterator := iterators[idx]

			defer func() {
				_ = iterator.Close()
				pendingWorkers.Done()
			}()

			for iterator.Next() {
				item, err := iterator.Item()
				if err != nil {
					return err
				}

				select {
				case <-groupCtx.Done():
					return groupCtx.Err()
				case iter.fanIn <- item:
				}
			}

			return nil
		})
	}

	workerGroup.Go(func() error {
		pendingWorkers.Wait()
		close(iter.fanIn)

		return nil
	})

	return iter
}

func (d *ChanIterator[T]) Next() bool {
	select {
	case item, ok := <-d.fanIn:
		if !ok {
			return false
		}
		select {
		case d.fanOut <- item:
		case <-d.ctx.Done():
			return false
		}

		return true

	case <-d.ctx.Done():
		return false
	}
}

func (d *ChanIterator[T]) Item() (T, error) {
	if d.current == nil {
		var empty T
		return empty, io.EOF
	}

	select {
	case item := <-d.fanOut:
		return item, nil
	case <-d.ctx.Done():
		var empty T
		return empty, d.ctx.Err()
	}
}

func (d *ChanIterator[T]) Close() error {
	return d.workerGroup.Wait()
}

func (d *ChanIterator[T]) Collect() ([]T, error) {
	d.mx.Lock()
	defer d.mx.Unlock()

	results := make([]T, 0, d.preallocatedItems)
	for d.Next() {
		item, err := d.Item()
		if err != nil {
			_ = d.Close()

			return nil, err
		}
		results = append(results, item)
	}

	return results, d.Close()
}

func (d *ChanIterator[T]) CollectPtr() ([]*T, error) {
	d.mx.Lock()
	defer d.mx.Unlock()

	results := make([]*T, 0, d.preallocatedItems)

	for d.Next() {
		item, err := d.Item()
		if err != nil {
			_ = d.Close()

			return nil, err
		}
		results = append(results, &item)
	}

	return results, d.Close()
}
