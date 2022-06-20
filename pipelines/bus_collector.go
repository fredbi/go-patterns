package pipelines

import (
	"context"
	"fmt"
	"reflect"
)

type (
	// BusCollector is a fan-in collector on many BUS channels, to listen on
	// all bus notifications from other pipelines.
	BusCollector[BUS any] struct {
		*commonPipeline[BUS, dummy, BUS]
		inputs   []<-chan BUS
		listener BusListener[BUS]
	}
)

// NewBusCollector builds a BusCollector to listen on bus channels.
func NewBusCollector[BUS any](opts ...Option) *BusCollector[BUS] {
	return &BusCollector[BUS]{
		commonPipeline: newCommonPipeline[BUS, dummy, BUS](opts...),
	}
}

func (p *BusCollector[BUS]) WithInputs(inputs ...<-chan BUS) *BusCollector[BUS] {
	p.inputs = append(p.inputs, inputs...)

	return p
}

func (p *BusCollector[BUS]) WithInputsFrom(busers ...Buser[BUS]) *BusCollector[BUS] {
	for _, buser := range busers {
		_ = p.WithInputs(buser.Bus())
	}

	return p
}

func (p *BusCollector[BUS]) WithBusListener(listener BusListener[BUS]) *BusCollector[BUS] {
	p.listener = listener
	p.in = nil // safeguard: don't use regular input
	p.out = nil
	p.bus = nil

	p.runner = p.defaultRunner()

	return p
}

func (p *BusCollector[BUS]) defaultRunner() func(context.Context, <-chan BUS, chan<- dummy, chan<- BUS) error {
	return func(ctx context.Context, _ <-chan BUS, _ chan<- dummy, _ chan<- BUS) error {
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

			busValue, ok := value.Interface().(BUS)
			if !ok {
				panic(fmt.Errorf("dev error: expected value on channel to be of type %T", p.bus))
			}

			if err := p.listener(ctx, busValue); err != nil {
				return err
			}
		}
	}
}
