package pipelines

import (
	"container/list"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type (
	// NamedRunnable is a Runnable identifiable by a name.
	NamedRunnable[BUS any] interface {
		Runnable[BUS]
		Name() string
	}

	// Collection is a chained list of named runnable pipelines with some common
	// BUS type to handle out-of-band notifications.
	Collection[BUS any] struct {
		chain *list.List
	}

	namedElement[BUS any] struct {
		name string
		Runnable[BUS]
	}
)

func newNamedElement[BUS any](name string, runnable Runnable[BUS]) *namedElement[BUS] {
	return &namedElement[BUS]{
		name:     name,
		Runnable: runnable,
	}
}

func (e *namedElement[BUS]) Name() string {
	return e.name
}

func NewCollection[BUS any]() *Collection[BUS] {
	return &Collection[BUS]{
		chain: list.New(),
	}
}

func (c *Collection[BUS]) mustBeNamedElement(e *list.Element) *namedElement[BUS] {
	if e == nil {
		panic("dev error: unexpected nil element in list")
	}

	n, ok := e.Value.(*namedElement[BUS])
	if !ok {
		panic(fmt.Sprintf("dev error: expected %T in list but got %T", n, e.Value))
	}

	return n
}

// Lookup in a collection for a given pipeline by its name.
func (c *Collection[BUS]) Lookup(name string) (NamedRunnable[BUS], bool) {
	for e := c.chain.Front(); e != nil; e = e.Next() {
		n := c.mustBeNamedElement(e)

		if n.Name() == name {
			return n, true
		}
	}

	return nil, false
}

func (c *Collection[BUS]) Add(piped ...Runnable[BUS]) {
	for _, toPin := range piped {
		pipe := toPin
		c.AddNamed("", pipe)
	}
}
func (c *Collection[BUS]) AddNamed(name string, piped Runnable[BUS]) {
	c.chain.PushBack(newNamedElement[BUS](name, piped))
}

func (c *Collection[BUS]) Append(piped ...Runnable[BUS]) *Collection[BUS] {
	c.Add(piped...)

	return c
}

func (c *Collection[BUS]) AppendNamed(name string, piped Runnable[BUS]) *Collection[BUS] {
	c.AddNamed(name, piped)

	return c
}

// AddBusCollector adds a BusCollector pipeline connected to all the bus channels published
// by the pipelines currently in the collection.
func (c *Collection[BUS]) AddBusCollector(collector *BusCollector[BUS]) {
	c.AddNamedBusCollector("", collector)
}

func (c *Collection[BUS]) AddNamedBusCollector(name string, collector *BusCollector[BUS]) {
	for _, buser := range c.AllRunnables() {
		collector = collector.WithInputsFrom(buser)
	}

	c.AddNamed(name, collector)
}

// AllRunnables returns all the runnable pipelines in the collection.
func (c *Collection[BUS]) AllRunnables() []Runnable[BUS] {
	runnables := make([]Runnable[BUS], 0, c.chain.Len())

	for e := c.chain.Front(); e != nil; e = e.Next() {
		n := c.mustBeNamedElement(e)
		runnables = append(runnables, n)
	}

	return runnables
}

// AllRunnersWithContext returns all runners as a collection of func() error suitable to
// execute in some errgroup.Group.
func (c *Collection[BUS]) AllRunnersWithContext(ctx context.Context) []func() error {
	runners := make([]func() error, 0, c.chain.Len())

	for e := c.chain.Front(); e != nil; e = e.Next() {
		n := c.mustBeNamedElement(e)

		runners = append(runners, n.RunWithContext(ctx))
	}

	return runners
}

// RunInGroup starts the pipeline inside some provided errgroup.Group.
//
// The caller is responsible for waiting on the group: RunInGroup only launches all runners as goroutines.
//
// The context passed should be the context produced when creating the group.
func (c *Collection[BUS]) RunInGroup(ctx context.Context, group *errgroup.Group) {
	// TODO: add options to add a hook (e.g. logger), using name as input
	for e := c.chain.Front(); e != nil; e = e.Next() {
		n := c.mustBeNamedElement(e)

		group.Go(n.RunWithContext(ctx))
	}
}
