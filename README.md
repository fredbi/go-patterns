# go-patterns

Musings with go1.18 generics to implement a few simple algorithms.

This repository reflects ongoing research to explore the capabilities of the new go1.18 generics.

This is _applied_ research, seeking to reduce the boiler plate and/or generated code on useful patterns.

The patters that I am exploring here mostly revolve around the idea of slicing or regrouping elements from an iterable data stream.

## iterators

A collection of iterator utitilies:
* an iterator is essentially something that knows what's `Next() bool` and collect the next `Item() T`
* the `iterators` package exposes 3 generic variants:
  1. Simple iterator from a slice `[]T` (e.g. to build mocks, etc)
  2. SQL rows iterator using sqlx.Rows and the `StructScan(interface{}) error` method.
     (this is used to iterate over unmarshaled structs scanned from a SQL cursor).
  3. FanIn iterator that joins a collection of input iterators in parallel (the result is unordered).
     Options: inner channel buffers

> NOTE: this is typically a piece of code to replace generated iterators.

### TODOs on iterator

* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.
* [ ] implement the channel-based fan-in iterator
* [ ] write testable examples

## batchers

* A batcher is used to execute some repeated action on a batch of T, every time there are enough elements pushed to the batch.
* 2 goroutine-safe methods: Push(T), Flush() 
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

A batcher is something that execute some processing in batches.

* [ ] assert performance - I expect that using a generic struct, not method, reduces the performance penalty due to the compiler's stencilinh.
* [x] introduce variations to shallow clone batched input elements (e.g. when we have `[]*TYPE` slices)
* [x] write testable examples

> Findings: at the moment, there is no easy way to perform a type assertion on the parametric types.
>
> For instance, it's unclear what kind of type I can pass. I don't know how to check that with built-in type constraints.
> As of go1.18.3, it looks like I have to build methods like `method(p TYPE)` and `methodPtr(p *TYPE)` specifically.
>
> I've been disappointed by the difficulty to use `any|*any` and be able in the code to figure out when this is a pointer type.


## pipelines

This is essentially intended to make my async code more readable and easier to guard against type error when mixing channels conveying messages of different types.

I don't want to mimic node's Promises. I just want plain async that can run multiple IN/OUT channels and check at a glance that types are correct.
I want the pipelines to be able to send out-of-band notifications to some listener ("bus").

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

# multi-sorter (TODO)

Objective: to produce a compound Less(int,int) bool method out of multiple individual criteria.

Type constraints:
* either `ordered`
* or `comparable` and implement some `Compare(T1,T2)`` int interface

e.g.: implement for any number of criteria (a,b,c ...):
```go 
Less := func(i, ) bool {
    if a[i] == a[j] {
        if b[i] == b[j] {
            return c[i] < c[j]
        }

        return b[i] < c[j]
    }

    return a[i] < a[j]
}
```

* produces a `Less() bool` function for multiple criteria.
  Example: a < b, if a == b then c < d,  if c == d, then e < f etc.
* API: `CompoundLes(criteria ...func(int, int) int) func(int, int) bool`
