package pipelines

import (
	"context"
)

// JoinerPipeline is a pipeline that takes input for 2 channels of different types and applies a Joiner method to this input.
//
// The implementation of the Joiner must be provided.
type JoinerPipeline[IN any, OTHER any, OUT any, BUS any] struct {
	*commonPipeline[IN, OUT, BUS]
	other  chan OTHER
	joiner Joiner[IN, OTHER, OUT, BUS]
}

func NewJoiner[IN any, OTHER any, OUT any, BUS any](opts ...Option) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	return &JoinerPipeline[IN, OTHER, OUT, BUS]{
		commonPipeline: newCommonPipeline[IN, OUT, BUS](opts...),
	}
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) WithJoiner(joiner Joiner[IN, OTHER, OUT, BUS]) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	p.joiner = joiner

	return p
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return func() error {
		defer p.close()

		return p.joiner(ctx, p.in, p.other, p.out, p.bus)
	}
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) WithInputsFrom(producer Producer[IN], otherProducer Producer[OTHER]) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	p.in = producer.Output()
	p.other = otherProducer.Output()

	return p
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) SetOther(other chan OTHER) {
	p.other = other
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) WithInputs(in chan IN, other chan OTHER) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	p.in = in
	p.other = other

	return p
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) WithOutput(out chan OUT) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	p.SetOutput(out)

	return p
}

func (p *JoinerPipeline[IN, OTHER, OUT, BUS]) WithBus(bus chan BUS) *JoinerPipeline[IN, OTHER, OUT, BUS] {
	p.SetBus(bus)

	return p
}
