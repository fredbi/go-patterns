package iterators_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/fredbi/go-patterns/iterators"
)

func ExampleChanIterator_Next() {
	baseIterators := []iterators.StructIterator[SampleStruct]{
		iterators.NewSliceIterator[SampleStruct](testSlice()),
		iterators.NewSliceIterator[SampleStruct](testSlice()),
	}
	count := 0

	iterator := iterators.NewChanIterator[SampleStruct](context.Background(), baseIterators)
	defer func() {
		_ = iterator.Close()
	}()
	items := make(SortableStructs, 0, 4)

	for iterator.Next() {
		item, err := iterator.Item()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("err: %v\n", err)
			}

			break
		}
		count++
		items = append(items, item)
	}

	sort.Sort(items)
	fmt.Printf("count: %d\n", count)
	fmt.Printf("items: %v\n", items)

	// Output:
	// count: 4
	// items: [{1 x}, {1 x}, {2 y}, {2 y}]
}

func ExampleChanIterator() {
	baseIterators := []iterators.StructIterator[SampleStruct]{
		iterators.NewSliceIterator[SampleStruct](testSlice()),
		iterators.NewSliceIterator[SampleStruct](testSlice()),
		iterators.NewSliceIterator[SampleStruct](testSlice()),
	}

	group, ctx := errgroup.WithContext(context.Background())
	iterator := iterators.NewChanIterator[SampleStruct](ctx, baseIterators)
	defer func() {
		_ = iterator.Close()
	}()
	var mx sync.Mutex
	items := make(SortableStructs, 0, 6)
	latch := make(chan struct{})

	// In this example, we iterate in parallel:
	// 3 producer iterators and 3 consumer iterators are running in parallel against
	// a single inner channel.
	for i := 0; i < 3; i++ {
		group.Go(func() error {
			<-latch
			count := 0
			defer func() {
				fmt.Fprintf(os.Stderr, "goroutine count: %d\n", count) // stderr doesn't count for example asserted output
			}()

			for iterator.Next() {
				item, err := iterator.Item()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return nil
					}

					return err
				}
				mx.Lock()
				items = append(items, item)
				mx.Unlock()
				count++
			}

			return nil
		})
	}
	close(latch)
	if err := group.Wait(); err != nil {
		fmt.Printf("an error occured: %v", err)
	}

	sort.Sort(items)
	fmt.Printf("count: %d\n", len(items))
	fmt.Printf("items: %v\n", items)

	// Output:
	// count: 6
	// items: [{1 x}, {1 x}, {1 x}, {2 y}, {2 y}, {2 y}]
}
