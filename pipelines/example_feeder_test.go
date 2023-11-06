// nolint:forbidigo
package pipelines_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/fredbi/go-patterns/pipelines"
	"golang.org/x/sync/errgroup"
)

type (
	exampleNotification struct {
		Msg   string
		Value int
	}

	exampleNotifications struct {
		mx    sync.Mutex
		inner []exampleNotification
	}
)

func (e *exampleNotifications) Add(in exampleNotification) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.inner = append(e.inner, in)
}

func (e *exampleNotifications) String() string {
	b := new(strings.Builder)
	e.mx.Lock()
	defer e.mx.Unlock()

	for _, in := range e.inner {
		fmt.Fprintf(b, "warn: notified of %s: %d\n", in.Msg, in.Value)
	}

	return b.String()
}

func (e *exampleNotifications) Sort() {
	sort.Slice(e.inner, func(i, j int) bool {
		return e.inner[i].Value < e.inner[j].Value
	})
}

var (
	exampleGenerator = func(ctx context.Context, out chan<- int, _ pipelines.NOBUSCHAN) error {
		// Feeder[OUT any, BUS any] func(context.Context, chan<- OUT, chan<- BUS) error
		inputs := []int{1, 2, 3, 4}
		for _, generated := range inputs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- generated:
			}
		}

		return nil
	}

	exampleCollector = func(ctx context.Context, in <-chan int, _ pipelines.NOBUSCHAN) error {
		// Collector[IN any, BUS any] func(context.Context, <-chan IN, chan<- BUS) error
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case received, isOpen := <-in:
				if !isOpen {
					return nil
				}
				fmt.Printf("received: %d\n", received)
			}
		}
	}
)

func ExampleFeederPipeline_withoutBus() {
	// This example creates a simple pipeline operation that feeds 3 integers to a collector that prints these out.
	pipes := pipelines.NewCollection[pipelines.NOBUS]()

	feeder := pipelines.NewFeeder[int, pipelines.NOBUS]().
		WithFeeder(exampleGenerator) // a feeder that produces integers
	pipes.Add(feeder)

	final := pipelines.NewCollector[int, pipelines.NOBUS]().
		WithInputFrom(feeder).
		WithCollector(exampleCollector) // a collector that prints integers
	pipes.Add(final)

	// executes the pipeline in some goroutines group
	group, ctx := errgroup.WithContext(context.Background())
	pipes.RunInGroup(ctx, group)

	if err := group.Wait(); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// Output:
	// received: 1
	// received: 2
	// received: 3
	// received: 4
}

func ExampleFeederPipeline_withBus() {
	// This example creates a simple pipeline operation that feeds 3 integers to a collector that prints these out.
	// Even integer are rejected an generate a notification.

	generator := func(ctx context.Context, out chan<- int, _ chan<- exampleNotification) error {
		// Feeder[OUT any, BUS any] func(context.Context, chan<- OUT, chan<- BUS) error
		return exampleGenerator(ctx, out, nil)
	}

	oddCollector := func(ctx context.Context, in <-chan int, bus chan<- exampleNotification) error {
		// Collector[IN any, BUS any] func(context.Context, <-chan IN, chan<- BUS) error
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case received, isOpen := <-in:
				if !isOpen {
					return nil
				}

				if received%2 == 0 {
					notice := exampleNotification{
						Msg:   "rejected even input",
						Value: received,
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case bus <- notice:
					}

					continue
				}

				fmt.Printf("received: %d\n", received)
			}
		}
	}

	pipes := pipelines.NewCollection[exampleNotification]()

	feeder := pipelines.NewFeeder[int, exampleNotification]().
		WithFeeder(generator) // a feeder that produces integers

	final := pipelines.NewCollector[int, exampleNotification]().
		WithInputFrom(feeder).      // connects the collector pipeline to the feeder
		WithCollector(oddCollector) // a collector that prints odd integers

	pipes.Add(feeder, final)

	notificationsOutlet := new(exampleNotifications)

	// registers a bus listener to process out-of-band notifications
	// NOTE: bus messages are processed in parallel
	pipes.AddBusCollector(
		// func NewBusCollector[BUS any](opts ...Option) *BusCollector[BUS]{
		pipelines.NewBusCollector[exampleNotification]().
			WithBusListener(
				// BusListener[BUS any] func(context.Context, BUS) error
				func(ctx context.Context, in exampleNotification) error {
					notificationsOutlet.Add(in)

					return nil
				},
			),
	)

	// executes the pipeline in some goroutines group
	group, ctx := errgroup.WithContext(context.Background())
	pipes.RunInGroup(ctx, group)

	if err := group.Wait(); err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// Notice that bus notices are processed asynchronously in parallel.
	// Therefore, their ordering is not guaranteed.
	notificationsOutlet.Sort()

	fmt.Println(notificationsOutlet.String())

	// Output:
	// received: 1
	// received: 3
	// warn: notified of rejected even input: 2
	// warn: notified of rejected even input: 4
}
