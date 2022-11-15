package pipelines

import (
	"context"
	"fmt"
	"reflect"
)

// FanInPipeline is a pipeline that gathers several input channels into one single ouput.
//
// The FanInPipeline does not transform the data type.
//
// There is a default runner implemented,
// with a few FanIn options available to introduce hooks, which could be used for example to:
// * collect metrics and logging
// * apply some data transform
type FanInPipeline[INOUT any, BUS any] struct {
	*commonPipeline[INOUT, INOUT, BUS]
	*fanInOptions[INOUT, BUS]
	inputs []<-chan INOUT
}

func NewFanIn[INOUT any, BUS any](opts ...Option) *FanInPipeline[INOUT, BUS] {
	return &FanInPipeline[INOUT, BUS]{
		commonPipeline: newCommonPipeline[INOUT, INOUT, BUS](opts...),
		fanInOptions:   defaultFanInOptions[INOUT, BUS](),
	}
}

func (p *FanInPipeline[INOUT, BUS]) WithFanInOptions(opts ...FanInOption[INOUT, BUS]) *FanInPipeline[INOUT, BUS] {
	for _, apply := range opts {
		apply(p.fanInOptions)
	}

	return p
}

func (p *FanOutPipeline[INOUT, BUS]) SetInput(_ chan INOUT) *FanOutPipeline[INOUT, BUS] {
	panic("unsupported")
}

func (p *FanInPipeline[INOUT, BUS]) SetInputs(inputs ...<-chan INOUT) {
	p.inputs = inputs
}

func (p *FanInPipeline[INOUT, BUS]) WithInputs(inputs ...<-chan INOUT) *FanInPipeline[INOUT, BUS] {
	p.SetInputs(inputs...)

	return p
}

// WithOutputTo alters the provided consumer to set its input to the current pipeline's output.
func (p *FanInPipeline[INOUT, BUS]) WithOutputTo(consumer SettableConsumer[INOUT]) *FanInPipeline[INOUT, BUS] {
	consumer.SetInput(p.Output())

	return p
}

func (p *FanInPipeline[INOUT, BUS]) WithOutput(out chan INOUT) *FanInPipeline[INOUT, BUS] {
	p.SetOutput(out)

	return p
}

func (p *FanInPipeline[INOUT, BUS]) WithBus(bus chan BUS) *FanInPipeline[INOUT, BUS] {
	p.SetBus(bus)

	return p
}

func (p *FanInPipeline[INOUT, BUS]) WithFanInFrom(producers ...Producer[INOUT]) *FanInPipeline[INOUT, BUS] {
	p.inputs = make([]<-chan INOUT, 0, len(producers))
	for _, producer := range producers {
		p.inputs = append(p.inputs, producer.Output())
	}

	p.in = nil // safeguard: don't use regular input
	p.runner = p.defaultRunner()

	return p
}

func (p *FanInPipeline[INOUT, BUS]) defaultRunner() func(context.Context, <-chan INOUT, chan<- INOUT, chan<- BUS) error {
	return func(ctx context.Context, _ <-chan INOUT, out chan<- INOUT, bus chan<- BUS) error {
		defer p.close()

		if out == nil {
			panic("dev error: a fan in pipeline must have an output channel")
		}

		cases := makeSelectCases(ctx, p.inputs)
		for {
			selected, value, ok := reflect.Select(cases) // dynamically built select { case ... } statement to listen on all participants to the bus while each of them close their outlet
			if selected == 0 {                           // case <-ctx.Done:
				return ctx.Err()
			}

			if !ok {
				// input closed: update the select case
				cases = removeElement(cases, selected)

				if len(cases) == 1 {
					return nil // we're done here: all producers are gone (only the case <- ctx.Done(): remains)
				}

				continue // enter the select again
			}

			output, ok := value.Interface().(INOUT)
			if !ok {
				panic(fmt.Errorf("dev error: expected value on channel to be of type %T", p.out))
			}

			for _, hook := range p.fanInHooks {
				if err := hook(ctx, output, bus); err != nil {
					return err
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case p.out <- output:
			}
		}
	}
}

// makeSelectCases dynamicaly builds a select over a slice of input channels.
// The first generate case is always on the channel being cancelled.
func makeSelectCases[IN any](ctx context.Context, inputs []<-chan IN) []reflect.SelectCase {
	cases := make([]reflect.SelectCase, 0, len(inputs)+1)

	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	for _, channel := range inputs {
		if channel == nil {
			continue
		}

		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(channel),
		})
	}

	return cases
}

func removeElement(s []reflect.SelectCase, i int) []reflect.SelectCase {
	// remove channel from slice (e.g. when closed). Does not preserve order
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
