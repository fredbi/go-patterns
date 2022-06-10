package pipelines

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type (
	NamedRunnable interface {
		Name() string
		RunWithContext(context.Context) func() error
	}

	Collection []NamedRunnable
)

func (collection Collection) Lookup(name string) NamedRunnable {
	for _, element := range collection {
		if element.Name() == name {
			return element
		}
	}

	return nil
}

func (collection Collection) AllRunnersWithContext(ctx context.Context) []func() error {
	runners := make([]func() error, 0, len(collection))
	for _, element := range collection {
		runners = append(runners, element.RunWithContext(ctx))
	}

	return runners
}

func (collection Collection) RunInGroup(ctx context.Context, group errgroup.Group) {
	// TODO: add options to add a hook (e.g. logger)
	for _, runner := range collection.AllRunnersWithContext(ctx) {
		group.Go(runner)
	}
}
