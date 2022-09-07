package iterators

import (
	"context"
	"errors"
	"io"
	"sync"

	"golang.org/x/sync/errgroup"
)

var _ StructIterator[dummy] = &ChanIterator[dummy]{}

// ChanIterator is a channel-based iterator that may be used to run a collection of StructIterators in parallel.
//
// Notice that its asynchronous working does not make it suitable to collect ordered items.
//
// The ChanIterator is goroutine-safe and may be iterated by several concurrent goroutines.
//
// The input is collected from a collection of input StructIterators, then may be collected by one or several goroutines
// reading from the ChanIterator using Next() and Item().
//
// Item() may return io.EOF is the iterator is done with producing records (e.g. some other consumer reached the end of the stream).
//
// WithChanFanInBuffers may be used to pre-fetch from input iterators asynchronously.
//
// Methods Collect() and CollectPrt() can't be used by concurrent goroutines and are protected against such a misuse.
type ChanIterator[T any] struct {
	fanIn       chan T
	workerGroup *errgroup.Group
	ctx         context.Context
	mx          sync.Mutex
	done        chan struct{}

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
		done:                make(chan struct{}),
	}

	if iter.fanInBuffers < 0 {
		iter.fanInBuffers = len(iterators)
	}

	iter.fanIn = make(chan T, iter.fanInBuffers)

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
		close(iter.done)

		return nil
	})

	return iter
}

func (d *ChanIterator[T]) Next() bool {
	select {
	case <-d.done:
		return len(d.fanIn) > 0
	case <-d.ctx.Done():
		return false
	default:
		return true
	}
}

func (d *ChanIterator[T]) Item() (T, error) {
	var empty T

	select {
	case item, ok := <-d.fanIn:
		if !ok {
			return empty, io.EOF
		}

		return item, nil
	case <-d.ctx.Done():
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
			if errors.Is(err, io.EOF) {
				break
			}

			return results, preferErrorOverContext(err, d.Close())
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
			if errors.Is(err, io.EOF) {
				break
			}

			return results, preferErrorOverContext(err, d.Close())
		}
		results = append(results, &item)
	}

	return results, d.Close()
}

// preferErrorOverContext returns a specific error preferrably
// to the generic "context cancelled" error, whenever available.
func preferErrorOverContext(err1, err2 error) error {
	isCancelled1 := errors.Is(err1, context.Canceled)
	isCancelled2 := errors.Is(err2, context.Canceled)

	if isCancelled1 && !isCancelled2 {
		return err2
	}

	return err1
}
