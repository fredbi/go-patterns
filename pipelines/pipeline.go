package pipelines

import (
	"context"
	"fmt"
	"reflect"
)

type (
	Consumer[IN chan any] interface {
		Input() IN
	}

	Producer[OUT chan any] interface {
		Output() OUT
	}

	Pipeline[IN chan any, OUT chan any, BUS <-chan any] struct {
		name   string
		runner func(context.Context) error
		in     IN
		out    OUT
		bus    BUS

		*options
	}
)

func New[IN chan any, OUT chan any, BUS <-chan any](name string, opts ...Option) *Pipeline[IN, OUT, BUS] {
	options := defaultOptions()
	for _, apply := range opts {
		apply(options)
	}

	return &Pipeline[IN, OUT, BUS]{
		name:    name,
		options: options,
	}
}

func (p *Pipeline[IN, OUT, BUS]) Name() string {
	return p.name
}

func (p *Pipeline[IN, OUT, BUS]) RunWithContext(ctx context.Context) func() error {
	return func() error {
		defer func() {
			if p.out != nil {
				close(p.out)
			}
		}()

		return p.runner(ctx)
	}
}

func (p *Pipeline[IN, OUT, BUS]) Input() IN {
	return p.in
}

func (p *Pipeline[IN, OUT, BUS]) Output() OUT {
	return p.out
}

func (p *Pipeline[IN, OUT, BUS]) Bus() BUS {
	return p.bus
}

func (p *Pipeline[IN, OUT, BUS]) WithRunner(runner func(context.Context) error) *Pipeline[IN, OUT, BUS] {
	p.runner = runner

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithNewInput() *Pipeline[IN, OUT, BUS] {
	p.in = make(IN, p.inBuffers)

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithNewOutput() *Pipeline[IN, OUT, BUS] {
	p.out = make(OUT, p.outBuffers)

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithNewBus() *Pipeline[IN, OUT, BUS] {
	p.bus = make(BUS)

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithInputFrom(producer Producer[IN]) *Pipeline[IN, OUT, BUS] {
	p.in = producer.Output()

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithRunFanOut(consumers ...Consumer[IN]) *Pipeline[IN, OUT, BUS] {
	outputs := make([]IN, 0, len(consumers))
	for _, consumer := range consumers {
		outputs = append(outputs, consumer.Input())
	}
	p.out = nil // safeguard don't use regular output
	p.runner = func(ctx context.Context) error {
		defer func() {
			for _, output := range outputs {
				close(output)
			}
		}()
		if p.in == nil {
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

				for _, output := range outputs {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- inout:
					}
				}
			}
		}

	}

	return p
}

// WithRunFanIn builds a pipeline that fans in the output from a list of producers
// into a single output channel.
//
// The output must be defined.
func (p *Pipeline[IN, OUT, BUS]) WithRunFanIn(producers ...Producer[OUT]) *Pipeline[IN, OUT, BUS] {
	inputs := make([]OUT, 0, len(producers))
	for _, producer := range producers {
		inputs = append(inputs, producer.Output())
	}
	p.in = nil // safeguard don't use regular input
	p.runner = func(ctx context.Context) error {
		defer func() {
			close(p.out)
		}()
		if p.out == nil {
			panic("dev error: a fan in pipeline must have an output channel")
		}

		cases := makeSelectCases(ctx, inputs)
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

			output, ok := value.Interface().(OUT)
			if !ok {
				panic(fmt.Errorf("dev error: expected value on channel to be of type %T", p.out))
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case p.out <- output:
			}
		}

		return nil
	}

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithInput(in IN) *Pipeline[IN, OUT, BUS] {
	p.in = in

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithOutput(out OUT) *Pipeline[IN, OUT, BUS] {
	p.out = out

	return p
}

func (p *Pipeline[IN, OUT, BUS]) WithBus(bus BUS) *Pipeline[IN, OUT, BUS] {
	p.bus = bus

	return p
}

// makeSelectCases dynamicaly builds a select over a slice of input channels.
// The first generate case is always on the channel being cancelled.
func makeSelectCases[IN chan any](ctx context.Context, inputs []IN) []reflect.SelectCase {
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
