package pipelines_test

// "strings"

/*
  TODO
func ExampleChainedPipeline_withoutBus() {
	// This example creates a simple pipeline operation that feeds 3 integers to a pipeline that computes their squares,
	// then to collector that prints out the result.
	pipes := pipelines.NewCollection[pipelines.NOBUS]()

	// a runner that computes the square of an integer
	squarer := func(ctx context.Context, in <-chan int, out chan<- uint, _ pipelines.NOBUSCHAN) error {
		// Runner[IN any, OUT any, BUS any] func(context.Context, <-chan IN, chan<- OUT, chan<- BUS) error
		//
		// Notice how both input and output are guarded with context-cancellation cases.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case received, isOpen := <-in:
				if !isOpen {
					return nil
				}

				result := uint(received * received)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- result:
				}
			}
		}
	}

	// TODO: cuber

	collector := func(ctx context.Context, in <-chan uint, _ pipelines.NOBUSCHAN) error {
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

	begin := pipelines.NewInitialChained[int, pipelines.NOBUS]().
		BeginsWith(
			pipelines.NewFeeder[int, pipelines.NOBUS]().
				WithFeeder(exampleGenerator), // a feeder that produces integers [int]
		)

	// a pipeline that consumes [int] and produces [uint]
	chain := pipelines.Then[pipelines.NOINPUT, int, uint, pipelines.NOBUS]( // the use of this extra call is required to explicit the new output type to the compiler
		begin,
		pipelines.NewPipeline[int, uint, pipelines.NOBUS]().WithRunner(squarer),
	).
		Finally( // a collector that consumes [uint]
			pipelines.NewCollector[uint, pipelines.NOBUS]().
				WithCollector(collector), // a collector that prints integers
		)

	spew.Dump(chain)
	pipes.Add(chain)

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
*/
