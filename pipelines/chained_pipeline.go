package pipelines

import (
	"container/list"
	"context"

	// "fmt"

	"golang.org/x/sync/errgroup"
)

type (
	// ChainedPipeline allows to construct chains of asynchronous pipelines with the Then idiom.
	//
	// The input of the chain is specified with "WithInput".
	//
	// The whole chain may run asynchronously with RunWithContext(), which blocks until all the elements
	// in the chain are done.
	//
	// The input of a ChainedPipeline corresponds to the input of its first pipeline.
	// The output of a ChainedPipeline corresponds to the output of its last pipeline.
	ChainedPipeline[IN any, OUT any, BUS any] struct {
		*commonPipeline[IN, OUT, BUS]
		chain *list.List
	}

	// InitialChainedPipeline is an output-only chained pipeline.
	//
	// The input cannot be defined on InitialChainedPipeline.
	InitialChainedPipeline[OUT any, BUS any] struct {
		inner *commonPipeline[dummy, OUT, BUS]
	}

	// FinalChainedPipeline is an input-only chained pipeline, that is returned when chaining with "Eventually".
	//
	// The output cannot be defined on FinalChainedPipeline.
	FinalChainedPipeline[IN any, BUS any] struct {
		inner *commonPipeline[IN, dummy, BUS]
		chain *list.List
	}

	/*
			Settable[IN any, OUT any] interface {
				SetInput(chan IN)
				SetOutput(chan OUT)
			}

		chainable interface {
			getChain() *list.List
		}
	*/
)

// Then allows to chain outputs by declaring explicit types.
//
// Notice that unfortunately, we can't make this a fluent method of ChainedPipeline as of go1.18.
func Then[IN any, OUT any, CHAINEDOUT any, BUS any](current *ChainedPipeline[IN, OUT, BUS], next *Pipeline[OUT, CHAINEDOUT, BUS]) *ChainedPipeline[IN, CHAINEDOUT, BUS] {
	next.SetInput(current.Output())

	n := NewChained[IN, CHAINEDOUT, BUS]()
	n.SetOutput(next.Output())

	n.chain.PushBack(next)
	n.chain.PushFrontList(current.chain)

	return n
}

func NewChained[IN any, OUT any, BUS any](opts ...Option) *ChainedPipeline[IN, OUT, BUS] {
	return &ChainedPipeline[IN, OUT, BUS]{
		commonPipeline: newCommonPipeline[IN, OUT, BUS](opts...),
		chain:          list.New(),
	}
}

// NewInitialChained builds a new initial chain of pipelines.
func NewInitialChained[OUT any, BUS any](opts ...Option) *InitialChainedPipeline[OUT, BUS] {
	return &InitialChainedPipeline[OUT, BUS]{
		inner: newCommonPipeline[dummy, OUT, BUS](opts...),
	}
}

func newFinalChained[OUT any, BUS any](opts ...Option) *FinalChainedPipeline[OUT, BUS] {
	return &FinalChainedPipeline[OUT, BUS]{
		inner: newCommonPipeline[OUT, dummy, BUS](opts...),
		chain: list.New(),
	}
}

// BeginsWith feeds a chained pipeline with some initial input pipeline.
func (p *InitialChainedPipeline[OUT, BUS]) BeginsWith(start *FeederPipeline[OUT, BUS]) *ChainedPipeline[dummy, OUT, BUS] {
	n := NewChained[dummy, OUT, BUS]()

	n.SetInput(nil)
	n.SetOutput(start.Output())
	n.chain.PushFront(start)

	return n
}

func (p *InitialChainedPipeline[OUT, BUS]) Output() chan OUT {
	return p.inner.Output()
}

func (p *InitialChainedPipeline[OUT, BUS]) Bus() chan BUS {
	return p.inner.Bus()
}

func (p *InitialChainedPipeline[OUT, BUS]) SetOutput(out chan OUT) {
	p.inner.SetOutput(out)
}

func (p *InitialChainedPipeline[OUT, BUS]) WithOutput(out chan OUT) *InitialChainedPipeline[OUT, BUS] {
	p.inner.SetOutput(out)

	return p
}

func (p *InitialChainedPipeline[OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return p.inner.RunWithContext(ctx)
}

func (p *FinalChainedPipeline[OUT, BUS]) Input() chan OUT {
	return p.inner.Input()
}

func (p *FinalChainedPipeline[OUT, BUS]) Bus() chan BUS {
	return p.inner.Bus()
}

func (p *FinalChainedPipeline[OUT, BUS]) SetInput(in chan OUT) {
	p.inner.SetInput(in)
}

func (p *FinalChainedPipeline[OUT, BUS]) WithInput(in chan OUT) *FinalChainedPipeline[OUT, BUS] {
	p.inner.SetInput(in)

	return p
}

func (p *FinalChainedPipeline[OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return p.inner.RunWithContext(ctx)
}

// Finally terminates a chained pipeline with some final collector pipeline.
func (p *ChainedPipeline[IN, OUT, BUS]) Finally(end *CollectorPipeline[OUT, BUS]) *FinalChainedPipeline[OUT, BUS] {
	end.SetInput(p.Output())

	n := newFinalChained[OUT, BUS](withCloneOptions(p.options))
	n.inner.options = p.options
	n.chain = p.chain
	n.inner.SetOutput(nil)
	n.chain.PushBack(end)

	return n
}

func (p *ChainedPipeline[IN, OUT, BUS]) mustBeRunnable(e *list.Element) Runnable[BUS] {
	val, ok := p.chain.Back().Value.(Runnable[BUS])
	if !ok {
		panic("dev error")
	}

	return val
}

// WithInputFrom sets the the input channel for this chain.
func (p *ChainedPipeline[IN, OUT, BUS]) WithInputFrom(producer Producer[IN]) *ChainedPipeline[IN, OUT, BUS] {
	p.in = producer.Output()

	return p
}

// WithOutputTo sets the input of the given consumer to the output of the current chain.
func (p *ChainedPipeline[IN, OUT, BUS]) WithOutputTo(consumer SettableConsumer[OUT]) *ChainedPipeline[IN, OUT, BUS] {
	consumer.SetInput(p.out)

	return p
}

// RunWithContext executes all the runners of a chain of pipelines as goroutines.
//
// The execution is interrupted if the context is cancelled.
//
// It blocks until all inner runners exit.
func (p *ChainedPipeline[IN, OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return func() error {
		group, groupCtx := errgroup.WithContext(ctx)
		for e := p.chain.Front(); e != nil; e = e.Next() {
			runnable := p.mustBeRunnable(e)
			group.Go(runnable.RunWithContext(groupCtx))
		}

		return group.Wait()
	}
}
