# go-patterns

Musings with go1.18 generics to implement a few simple algorithms.

## iterators

A collection of iterator utitilies:
* an iterator is essentially something that knows what's `Next() bool` and collect the next `Item() T`
* the `iterators` package exposes 3 generic variants:
  1. Simple iterator from a slice `[]T` (e.g. to build mocks, etc)
  2. SQL rows iterator using sqlx.Rows and the `StructScan(interface{}) error` method.
     (this is used to iterate over unmarshaled structs scanned from a SQL cursor).
  3. FanIn iterator that joins a collection of input iterators in parallel (the result is unordered).
     Options: inner channel buffers


## batchers

* A batcher is used to execute some repeated action on a batch of T, every time there are enough elements pushed to the batch.
* 2 goroutine-safe methods: Push(T), Flush() 
* The `executor func([]T)` is assumed to handle errors etc. It is executed when the batch size is reached or on Flush().
* Options to consider (TODO):
    * Shallow clone of batch elements
    * Timeout on buffering wait
* TODO: insertBatcher
  * a common specialized usage of the batcher to construct Postgres multi-values batch INSERTs
* TODO: parallelBatcher
  * run executors as parallel go routines with a throttle


# multi-sorter (TODO)
(no need for generics here)

* produces a `Less() bool` function for multiple criteria.
  Example: a < b, if a == b then c < d,  if c == d, then e < f etc.
* interface: `CompoundLes(criteria ...func(int, int) int) func(int, int) bool`
