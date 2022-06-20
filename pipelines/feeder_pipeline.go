package pipelines

import (
	"context"
)

type (
	// FeederPipeline is a pipeline that only outputs, with a message generating Feeder method.
	//
	// A FeederPipeline is a starting point in your pipeline dependency graph.
	// TODO: use inner to shut inapplicable methods
	FeederPipeline[OUT any, BUS any] struct {
		*commonPipeline[dummy, OUT, BUS]
	}

	dummy struct{}
)

// NewFeeder builds a new FeederPipeline.
func NewFeeder[OUT any, BUS any](opts ...Option) *FeederPipeline[OUT, BUS] {
	return &FeederPipeline[OUT, BUS]{
		commonPipeline: newCommonPipeline[dummy, OUT, BUS](opts...),
	}
}

// TODO: should be part of the constructor
func (p *FeederPipeline[OUT, BUS]) WithFeeder(feeder Feeder[OUT, BUS]) *FeederPipeline[OUT, BUS] {
	p.runner = func(ctx context.Context, _ <-chan dummy, out chan<- OUT, bus chan<- BUS) error {
		return feeder(ctx, out, bus)
	}

	return p
}

func (p *FeederPipeline[OUT, BUS]) WithOutput(out chan OUT) *FeederPipeline[OUT, BUS] {
	p.SetOutput(out)

	return p
}

func (p *FeederPipeline[OUT, BUS]) WithBus(bus chan BUS) *FeederPipeline[OUT, BUS] {
	p.SetBus(bus)

	return p
}
