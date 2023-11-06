![Lint](https://github.com/fredbi/go-patterns/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/go-patterns/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/go-patterns/badge.svg)](https://coveralls.io/github/fredbi/go-patterns)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/go-patterns)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-patterns.svg)](https://pkg.go.dev/github.com/fredbi/go-patterns)

# go-patterns

This repository exposes personal musings with go1.18+ generics to implement a few algorithms.

> This repository reflects ongoing research to explore the capabilities of generic types available since go 1.18.
>
> Most patterns that I am exploring here revolve around the idea of slicing or regrouping elements from an iterable data stream.

## Outcome

I've hit two slightly annoying limitations of the current implementation of go generics, as of go1.19:

1. It is not possible to perform type assertions on generics. The workaround is awkward.
2. It is not possible to make a method with a parameterized type. There is no working around that.

## Iterators

A collection of [iterators](iterators/README.md) to iterate in different ways over collections.

**Main interface:**
```go
	// StructIterator is an iterator that delivers items of some type T.
	StructIterator[T any] interface {
		Iterator

		// Item returns the current iterated item.
		//
		// Next() must have been called at least once.
		Item() (T, error)

		// Collect returns all items in one slice, then closes the iterator
		Collect() ([]T, error)

		// CollectPtr returns all items in one slice of pointers, then closes the iterator
		CollectPtr() ([]*T, error)
	}
```

**Features:**
* an iterator based on slice (no big deal, this is intended for testing...)
* an iterator that knows how to deal with `sqlx.Rows` fetched from a database
* an iterator that knows how to transform the stream from another iterator, possibly hiding some state in the context

[Testable example](iterators/example_rows_iterator_test.go)

## Batchers

A [batcher](batchers/README.md) is something that executes some processing in batches, slicing an input stream into batches of fixed length.

**Features:**

TODO[Testable example](batchers/example_batcher_test.go)

# multi-sorter

This package provides a way to apply compound sorting criteria to a collection of any type.

**Use-case**: we want to sort a slice of `struct{A int, B int, C int}` by (i) A, then (ii) B,  then (iii) C,
providing semantics similar to the SQL `ORDER BY A,B,C` statement.


**Main interface:**
```go
	// Comparison knows how to compare 2 values of the same type.
	//
	// Return values:
	//
	// if a == b: 0
	// if a < b: -1
	// if a > b: 1
	Comparison[T any] func(a, b T) int

    // NewMulti produces a sortable object, that supports multiple sorting criteria.
    func NewMulti[T any](collection []T, criteria ...Comparison[T]) *Multi[T]
```

**Features:**
* multi-criteria sort
* utilities to build comparison operators for common types and pointers.
* a method to revert the ordering (e.g. like `ORDER BY A DESC`)
* sorting options to apply language-specifics when sorting `strings` or `[]byte`

[Testable example](sorters/example_multi_sorter_test.go)

## Pipelines

The [`pipelines` package](pipelines/README.md) offers an abstraction to build asynchronous pipelines.

This was initially intended to make my async code more readable and easier to guard against 
type error when mixing channels conveying messages of different types.

**Main interface:**

**Features:**
* fan-in/fan-out

## upserter (TODO)

The upserter knows how to carry out batch inserts (resp. upserts) from some input channel, in parallel.
It is a special case of the `batcher` pattern.
This leverage the Postgres multiple `VALUES()` syntax. There is also a slightly different variant for `cockroachDB`.

TODOs:
[ ] Publish generic implementation
