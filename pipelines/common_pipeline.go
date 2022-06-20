package pipelines

import (
	"context"
)

type commonPipeline[IN any, OUT any, BUS any] struct {
	in  chan IN
	out chan OUT
	bus chan BUS

	runner Runner[IN, OUT, BUS]

	*options
}

func newCommonPipeline[IN any, OUT any, BUS any](opts ...Option) *commonPipeline[IN, OUT, BUS] {
	options := defaultOptions()
	for _, apply := range opts {
		apply(options)
	}

	return &commonPipeline[IN, OUT, BUS]{
		options: options,
		in:      make(chan IN, options.inBuffers),
		out:     make(chan OUT, options.outBuffers),
		bus:     make(chan BUS, options.busBuffers),
	}
}

func (c *commonPipeline[IN, OUT, BUS]) close() {
	if c.out != nil && c.autoCloseOutput {
		close(c.out)
	}

	if c.bus != nil {
		close(c.bus)
	}
}

func (c *commonPipeline[IN, OUT, BUS]) Input() chan IN {
	return c.in
}

func (c *commonPipeline[IN, OUT, BUS]) Output() chan OUT {
	return c.out
}

func (c *commonPipeline[IN, OUT, BUS]) Bus() chan BUS {
	return c.bus
}

func (c *commonPipeline[IN, OUT, BUS]) SetInput(in chan IN) {
	c.in = in
}

func (c *commonPipeline[IN, OUT, BUS]) SetOutput(out chan OUT) {
	c.out = out
}

func (c *commonPipeline[IN, OUT, BUS]) SetBus(bus chan BUS) {
	c.bus = bus
}

func (c *commonPipeline[IN, OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return func() error {
		defer c.close()

		return c.runner(ctx, c.in, c.out, c.bus)
	}
}
