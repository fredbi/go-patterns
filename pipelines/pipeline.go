package pipelines

type (
	// Pipeline is a generic pipe from channel IN to channel OUT and
	// reporting errors or notices on channel BUS.
	Pipeline[IN any, OUT any, BUS any] struct {
		*commonPipeline[IN, OUT, BUS]
	}
)

// NewPipeline builds a generic pipeline that collects from a channel of type IN and outputs into a channel of type OUT,
// without out-of-band notifications on channel of type BUS.
func NewPipeline[IN any, OUT any, BUS any](opts ...Option) *Pipeline[IN, OUT, BUS] {

	return &Pipeline[IN, OUT, BUS]{
		commonPipeline: newCommonPipeline[IN, OUT, BUS](opts...),
	}
}

// WithRunner returns a pipeline with a runner function to transform IN into OUT.
func (p *Pipeline[IN, OUT, BUS]) WithRunner(runner Runner[IN, OUT, BUS]) *Pipeline[IN, OUT, BUS] {
	p.runner = runner

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithInputFrom(producer Producer[IN]) *Pipeline[IN, OUT, BUS] {
	return p.WithInput(producer.Output())
}

// WithOutputTo alters the provided consumer to set its input to the current pipeline's output.
func (p *Pipeline[IN, OUT, BUS]) WithOutputTo(consumer SettableConsumer[OUT]) *Pipeline[IN, OUT, BUS] {
	consumer.SetInput(p.Output())

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithInput(in chan IN) *Pipeline[IN, OUT, BUS] {
	p.SetInput(in)

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithOutput(out chan OUT) *Pipeline[IN, OUT, BUS] {
	p.SetOutput(out)

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithBus(bus chan BUS) *Pipeline[IN, OUT, BUS] {
	p.SetBus(bus)

	return p
}
