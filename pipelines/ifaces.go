package pipelines

import (
	"context"
)

type (
	// NOBUS is a placeholder when the BUS feature is unused.
	NOBUS dummy
	// NOBUSCHAN is an alias for chan <- NOBUS when the BUS feature is unused.
	NOBUSCHAN = chan<- NOBUS

	NOINPUT  = dummy
	NOOUTPUT = dummy

	// Consumer knows what is the input channel.
	Consumer[IN any] interface {
		Input() chan IN
	}

	// SettableConsumer is a Consumer that may have its input set.
	SettableConsumer[IN any] interface {
		Consumer[IN]
		SetInput(chan IN)
	}

	// Producer knows what is the output channel.
	Producer[OUT any] interface {
		Output() chan OUT
	}

	// SettableProducer is a Producer that may have its output set.
	SettableProducer[OUT any] interface {
		Producer[OUT]
		SetOutput(chan OUT)
	}

	// Buser knows about its bus channel
	Buser[BUS any] interface {
		Bus() chan BUS
	}

	// Runnable knows how to execute a runner with a context.
	//
	// Runnables SHOULD relinquish resources and exit when the context is cancelled or when the input channel is closed.
	// They are not expected to close the channels.
	//
	// Whenever the bus feature is not put to use, you may use the placeholder type NOBUS: "Runnable[NOBUS]"
	Runnable[BUS any] interface {
		RunWithContext(context.Context) func() error
		Buser[BUS]
	}

	// Runner is the function executed inside a pipeline.
	//
	// The runner SHOULD relinquish resources and exit when the context is cancelled or when the input channel is closed.
	//
	// If the runner takes care of closing the ouput channel, the executing pipeline should explictly no do so
	// redundantly (create the pipeline with "WithAutoCloseOutput(false)" in this case).
	Runner[IN any, OUT any, BUS any] func(context.Context, <-chan IN, chan<- OUT, chan<- BUS) error

	// Joiner is the two-ways join executed inside a JoinerPipeline.
	//
	// The joiner SHOULD relinquish resources and exit when the context is cancelled or when the input channel is closed.
	//
	// If the joiner takes care of closing the ouput channel, the executing pipeline should explictly no do so
	// redundantly (create the pipeline with "WithAutoCloseOutput(false)" in this case).
	Joiner[IN any, OTHER any, OUT any, BUS any] func(context.Context, <-chan IN, <-chan OTHER, chan<- OUT, chan<- BUS) error

	// Feeder is the output-only function executed inside a FeederPipeline.
	//
	// The feeder SHOULD relinquish resources and exit when the context is cancelled.
	//
	// If the feeder takes care of closing the ouput channel, the executing pipeline should explictly no do so
	// redundantly (create the pipeline with "WithAutoCloseOutput(false)" in this case).
	Feeder[OUT any, BUS any] func(context.Context, chan<- OUT, chan<- BUS) error

	// Collector is the input-only function executed inside a CollectorPipeline.
	//
	// The collector SHOULD relinquish resources and exit when the context is cancelled or when the input channel is closed.
	Collector[IN any, BUS any] func(context.Context, <-chan IN, chan<- BUS) error

	// FanHook is a hook function executed by the FanInPipeline or the FanOutPipeline runner.
	FanHook[INOUT any, BUS any] func(context.Context, INOUT, chan<- BUS) error

	// BusListener defines a special collector for BUS entries.
	BusListener[BUS any] func(context.Context, BUS) error
)
