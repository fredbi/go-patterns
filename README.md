![Lint](https://github.com/fredbi/go-patterns/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/go-patterns/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/go-patterns/badge.svg)](https://coveralls.io/github/fredbi/go-patterns)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/go-patterns)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-patterns.svg)](https://pkg.go.dev/github.com/fredbi/go-patterns)

# go-patterns

Musings with go1.18 generics to implement a few simple algorithms.

> This repository reflects ongoing research to explore the capabilities of the new go1.18 generics.
>
> This is _applied_ research, seeking to reduce the boiler plate and/or generated code on useful patterns.
> Most patterns that I am exploring here revolve around the idea of slicing or regrouping elements from an iterable data stream.

## Iterators

A collection of iterator utitilies:
* an iterator is essentially something that knows what's `Next() bool` and collects the next available `Item() T`
* the `iterators` package exposes 3 generic variants:
  1. A simple iterator over an underlying slice `[]T` (e.g. to build mocks, etc)
  2. SQL rows iterator using `github.com/jmoiron/sqlx.Rows` and the `StructScan(interface{}) error` method.
     (this is used to iterate over unmarshaled structs scanned from a SQL cursor).
  3. A `ChanIterator` that joins a collection of input iterators in parallel (the result is unordered).
  4. A `TransformIterator` that applies a data transform on the iterations of some other base iterator.

> NOTE: I like the iterator pattern a lot when it comes to fetch from a database an arbitrary number of rows.
> Iterators allow a stream of data to traverse all the layers of an app without undue intermediary buffering.

Sample code: see the full [testable example](iterators/example_rows_iterator_test.go)
```go
    // create DB with some data
    ...

	// open a cursor selecting over the test data.
    // Typically, rows come from a SQL query.
	rows, err := testdb.OpenDBCursor(db)
	if err != nil {
		log.Fatalf("could not query DB: %v", err)
	}

    // testdb.DummyRow is the go type receiving unmarshalled data
	iterator := iterators.NewRowsIterator[*sqlx.Rows, testdb.DummyRow](rows)
	defer func() {
		_ = iterator.Close()
	}()
	count := 0

	for iterator.Next() {
		item, err := iterator.Item()
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		count++
		fmt.Printf("item: %#v\n", item)
	}
```

### TODOs on iterator

* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.

## Batchers

A batcher is something that executes some processing in batches.

* A batcher is used to execute some repeated action on a batch of elements with type `T`, every time there are enough elements pushed to the batch.
* It is intended to apply an executor function to a stream of inputs. The executor is a function operating on a slice `[]T`of fixed maximum size.
* The interface is  minimal: `Push(T)`, `Flush()` (safe to execute from concurrent go routines)
* The `executor func([]T)` is assumed to handle errors etc. It is executed when the batch size is reached or on `Flush()`.

Sample code: [testable example](batchers/batcher_examples_test.go)

### TODOs on batchers

* Options to consider:
    * Optional shallow clone of batch elements (atm cloned by default on pointers)
    * Timeout on buffering wait
* TODO: InsertBatcher
  * a common specialized usage of the batcher to construct Postgres multi-values batch INSERTs
* TODO: ParallelBatcher
  * run executors as parallel go routines with a throttle
* TODO: ErrBatcher
  * executor may return an error
* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.
* [x] introduce variations to shallow clone batched input elements (e.g. when we have `[]*TYPE` slices)
* [x] write testable examples

> Findings: at the moment, there is no easy way to perform a type assertion on the parametric types.
>
> For instance, it's unclear what kind of type I can pass. I don't know how to check that with built-in type constraints.
> As of go1.18.3, it looks like I have to build methods like `method(p TYPE)` and `methodPtr(p *TYPE)` specifically.
>
> I've been disappointed by the difficulty to use `any|*any` while keeping the ability in the code to figure out when this is a pointer type.

# multi-sorter

This package provides a way to apply compound sorting criteria to a collection of any type.

Use-case: we want to sort a slice of `struct{A int, B int, C int}` by (i) A, then (ii) B,  then (iii) C,
providing semantics similar to the SQL `ORDER BY A,B,C` statement.

It provides utilities to build comparison operators for common types and pointers.

Sample code: [testable example](sorters/example_multi_sorter_test.go)

## Pipelines

This is intended to make my async code more readable and easier to guard against type error when mixing channels conveying messages of different types.

I don't want to mimic node's Promises.
I just want plain async that can run multiple IN/OUT channels and check at a glance that types are correct.
I want the pipelines to be able to send out-of-band notifications to some listener ("bus").

TODO: basic pipelines work, but it is still difficult to get a nice, idiomatic chain of pipelines running smooth as I had expected.

Pipelines patterns to support:
* in/out/bus
* fan-int
* fan-out
* 2-way hetereogenous join
* feeder (no input)
* collector (no output)

> Findings
> I realize the implications of the limitation that no method can be itself parametric:
> this totally prevents me from building a fluent pipeline chain with a method like `Then[NEWOUT](next *Pipeline[OUT,NEWOUT]) *ChainedPipeline[IN, NEWOUT]`
> Ugh. Need to reflect more on that one.

## upserter (TODO)

The upserter knows how to carry out batch inserts (resp. upserts) from some input channel, in parallel.
It is a special case of the `batcher` pattern.
This leverage the Postgres multiple `VALUES()` syntax. There is also a slightly different variant for `cockroachDB`.

### TODOs on upserters

[ ] Publish generic implementation
