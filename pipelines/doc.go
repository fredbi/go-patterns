// Package pipelines exposes a generic implementation of asynchronous pipelines.
//
// A pipeline is defined by:
// * an input channel of some type: <-chan IN
// * an output channel of some type: chan <- OUT
// * an optional outgoing out-of-band notification bus channel: chan <- BUS
// * a runner function func(context.CONTEXT, <- chan IN, chan <- OUT, chan <- BUS)
//
// This implementation puts forward type safety when connecting channels of different types with each other:
// the consistency of your pipeline scaffolding should be ensured at build time.
//
// Error handling
//
// Pipelines are not akin to "promises" as there is no handling of the rejected state.
// Blocking errors should be handled by having the runner return an error.
// Non-blocking errors should be handled by sending some notification to the bus channel.
//
// The package allows to define and manipulate pipeline with different characteristics:
// * Pipeline is a plain IN/OUT pipeline
// * FeederPipeline is a pipeline that only produces output (e.g. from some other non-channel input, such as a io.Reader or DB records)
// * CollectorPipeline is a pipeline that only consumes input (e.g. it produces some final non-channel output)
// * FanInPipeline implements a (many IN -> single OUT) fan-in pattern
// * FanOutPipeline implements a (single IN -> many OUT) fan-out pattern
// * JoinerPipeline provides a means to implement a 2-way join pattern: (IN, OTHER) -> OUT.
// * ChainedPipeline provides a means to declare pipelines easily using the BeginsWith, Then and Finally builder methods.
package pipelines
