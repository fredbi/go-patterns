package pipelines

import (
	"context"
)

type (
	// CollectorPipeline is a dead-end pipeline, only consuming from input and producing not output channel.
	//
	// This is a termination point in your pipeline dependency graph.
	CollectorPipeline[IN any, BUS any] struct {
		*commonPipeline[IN, dummy, BUS]
	}
)

// NewCollector builds a CollectorPipeline that reads from a channel of type IN.
func NewCollector[IN any, BUS any](opts ...Option) *CollectorPipeline[IN, BUS] {
	return &CollectorPipeline[IN, BUS]{
		commonPipeline: newCommonPipeline[IN, dummy, BUS](opts...),
	}
}

func (p *CollectorPipeline[IN, BUS]) WithCollector(collector Collector[IN, BUS]) *CollectorPipeline[IN, BUS] {
	p.runner = func(ctx context.Context, in <-chan IN, _ chan<- dummy, bus chan<- BUS) error {
		return collector(ctx, in, bus)
	}

	return p
}

func (p *CollectorPipeline[IN, BUS]) WithInputFrom(producer Producer[IN]) *CollectorPipeline[IN, BUS] {
	return p.WithInput(producer.Output())
}

func (p *CollectorPipeline[IN, BUS]) WithInput(in chan IN) *CollectorPipeline[IN, BUS] {
	p.SetInput(in)

	return p
}

func (p *CollectorPipeline[IN, BUS]) WithBus(bus chan BUS) *CollectorPipeline[IN, BUS] {
	p.SetBus(bus)

	return p
}
