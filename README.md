# go-patterns

Musings with go1.18 generics to implement a few simple algorithms.

This repository reflects ongoing research to explore the capabilities of the new go1.18 generics.

This is _applied_ research, seeking to reduce the boiler plate and/or generated code on useful patterns.

The patterns that I am exploring here mostly revolve around the idea of slicing or regrouping elements from an iterable data stream.

## Iterators

A collection of iterator utitilies:
* an iterator is essentially something that knows what's `Next() bool` and collect the next `Item() T`
* the `iterators` package exposes 3 generic variants:
  1. Simple iterator from a slice `[]T` (e.g. to build mocks, etc)
  2. SQL rows iterator using github.com/jmoiron/sqlx.Rows and the `StructScan(interface{}) error` method.
     (this is used to iterate over unmarshaled structs scanned from a SQL cursor).
  3. ChanIterator that joins a collection of input iterators in parallel (the result is unordered).
  4. TransformIterator that applies a data transform on the iterations of some other base iterator.

> NOTE: I like the iterator pattern a lot when it comes to fetch from a database an arbitrary number of rows.
> Iterators allow a stream of data to traverse all the layers of an app without undue intermediary buffering.

### TODOs on iterator

* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.

## Batchers

* A batcher is used to execute some repeated action on a batch of T, every time there are enough elements pushed to the batch.
* It is intended to apply to a stream of inputs, some executor function on a slice of fixed maximum size.
* Simple interface: 2 goroutine-safe methods: Push(T), Flush() 
* The `executor func([]T)` is assumed to handle errors etc. It is executed when the batch size is reached or on Flush().

* Options to consider:
    * Optional shallow clone of batch elements (atm cloned by default on pointers)
    * Timeout on buffering wait
* TODO: InsertBatcher
  * a common specialized usage of the batcher to construct Postgres multi-values batch INSERTs
* TODO: ParallelBatcher
  * run executors as parallel go routines with a throttle
* TODO: ErrBatcher
  * executor may return an error

### TODOs on batchers

A batcher is something that executes some processing in batches.

* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.
* [x] introduce variations to shallow clone batched input elements (e.g. when we have `[]*TYPE` slices)
* [x] write testable examples

> Findings: at the moment, there is no easy way to perform a type assertion on the parametric types.
>
> For instance, it's unclear what kind of type I can pass. I don't know how to check that with built-in type constraints.
> As of go1.18.3, it looks like I have to build methods like `method(p TYPE)` and `methodPtr(p *TYPE)` specifically.
>
> I've been disappointed by the difficulty to use `any|*any` and be able in the code to figure out when this is a pointer type.

# multi-sorter

This package provides a way to apply compound sorting criteria to a collection of any type.

It provides utilities to build comparison operators for common types and pointers.

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

# upserter (TODO)

The upserter knows how to carry out batch inserts (resp. upserts) from some input channel, in parallel.
This leverage the Postgres multiple `VALUES()` syntax. There is also a slightly different variant for `cockroachDB`.

### TODOs on upserters

[ ] Publish generic implementation
