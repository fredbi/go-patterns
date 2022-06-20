package pipelines

import (
	"context"
)

// FanOutPipeline takes one input and duplicate the input messages to several output channels.
//
// FanOutPipeline does not change the data type.
//
// There is a default runner implemented,
// with a few options available to introduce hooks, which could be used to:
// * collect metrics & logging
// * cloning output messages if they pass references that are mutated downstream
//
type FanOutPipeline[INOUT any, BUS any] struct {
	*commonPipeline[INOUT, INOUT, BUS]
	*fanOutOptions[INOUT, BUS]
	outputs []chan<- INOUT
}

func NewFanOut[INOUT any, BUS any](opts ...Option) *FanOutPipeline[INOUT, BUS] {
	return &FanOutPipeline[INOUT, BUS]{
		commonPipeline: newCommonPipeline[INOUT, INOUT, BUS](opts...),
		fanOutOptions:  defaultFanOutOptions[INOUT, BUS](),
	}
}

func (p *FanOutPipeline[INOUT, BUS]) WithFanOutOptions(opts ...FanOutOption[INOUT, BUS]) *FanOutPipeline[INOUT, BUS] {
	for _, apply := range opts {
		apply(p.fanOutOptions)
	}

	return p
}

func (p *FanOutPipeline[INOUT, BUS]) SetOutput(_ chan INOUT) *FanOutPipeline[INOUT, BUS] {
	panic("unsupported")
}

func (p *FanOutPipeline[INOUT, BUS]) WithInputFrom(producer Producer[INOUT]) *FanOutPipeline[INOUT, BUS] {
	return p.WithInput(producer.Output())
}

func (p *FanOutPipeline[INOUT, BUS]) WithInput(in chan INOUT) *FanOutPipeline[INOUT, BUS] {
	p.SetInput(in)

	return p
}

func (p *FanOutPipeline[INOUT, BUS]) WithBus(bus chan BUS) *FanOutPipeline[INOUT, BUS] {
	p.SetBus(bus)

	return p
}

func (p *FanOutPipeline[INOUT, BUS]) WithFanOutTo(consumers ...Consumer[INOUT]) *FanOutPipeline[INOUT, BUS] {
	p.outputs = make([]chan<- INOUT, 0, len(consumers))
	for _, consumer := range consumers {
		p.outputs = append(p.outputs, consumer.Input())
	}

	p.out = nil // safeguard don't use regular output
	p.runner = p.defaultRunner()

	return p
}

func (p *FanOutPipeline[INOUT, BUS]) defaultRunner() func(context.Context, <-chan INOUT, chan<- INOUT, chan<- BUS) error {
	return func(ctx context.Context, in <-chan INOUT, _ chan<- INOUT, bus chan<- BUS) error {
		defer func() {
			if p.autoCloseOutput {
				for _, output := range p.outputs {
					if output != nil {
						close(output)
					}
				}
			}

			if p.bus != nil {
				close(p.bus)
			}
		}()

		if in == nil {
			panic("dev error: a fan out pipeline must have an input channel")
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case inout, isOpen := <-p.in:
				if !isOpen {
					return nil
				}

				for _, hook := range p.fanOutHooks {
					if err := hook(ctx, inout, bus); err != nil {
						return err
					}
				}

				for _, output := range p.outputs {
					if output == nil {
						continue
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- inout:
					}
				}
			}
		}
	}
}
