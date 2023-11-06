# batchers

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
