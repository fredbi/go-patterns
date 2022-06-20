package pipelines

type (
	// Option to tune the behavior of a pipeline.
	Option func(*options)

	options struct {
		name            string
		autoCloseOutput bool
		inBuffers       int
		outBuffers      int
		busBuffers      int
	}

	// FanInOption alters the behavior of the fan-in runner.
	FanInOption[INOUT any, BUS any] func(*fanInOptions[INOUT, BUS])

	// FanOutOption alters the behavior of the fan-out runner.
	FanOutOption[INOUT any, BUS any] func(*fanOutOptions[INOUT, BUS])

	fanInOptions[INOUT any, BUS any] struct {
		fanInHooks []FanHook[INOUT, BUS]
	}

	fanOutOptions[INOUT any, BUS any] struct {
		fanOutHooks []FanHook[INOUT, BUS]
	}
)

func defaultOptions() *options {
	return &options{
		autoCloseOutput: true,
	}
}

func defaultFanOutOptions[INOUT any, BUS any]() *fanOutOptions[INOUT, BUS] {
	return &fanOutOptions[INOUT, BUS]{}
}

func defaultFanInOptions[INOUT any, BUS any]() *fanInOptions[INOUT, BUS] {
	return &fanInOptions[INOUT, BUS]{}
}

func (o *options) Name() string {
	return o.name
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithOutputBuffers sets a buffered output cha,nel.
func WithOutputBuffers(channelBuffers int) Option {
	return func(o *options) {
		o.outBuffers = channelBuffers
	}
}

// WithInputBuffers sets a buffered input channel.
func WithInputBuffers(channelBuffers int) Option {
	return func(o *options) {
		o.inBuffers = channelBuffers
	}
}

// WithBusBuffers sets a buffered bus channel.
func WithBusBuffers(channelBuffers int) Option {
	return func(o *options) {
		o.busBuffers = channelBuffers
	}
}

// WithAutoCloseOutput sets the behavior to close the output channel after the runner is complete.
//
// This defaults to true and can be disabled if the runner is already taking care of closing its output channel.
func WithAutoCloseOutput(enabled bool) Option {
	return func(o *options) {
		o.autoCloseOutput = enabled
	}
}

func withCloneOptions(opts *options) Option {
	return func(o *options) {
		*o = *opts
	}
}

func WithFanOutHooks[INOUT any, BUS any](hooks ...FanHook[INOUT, BUS]) FanOutOption[INOUT, BUS] {
	return func(o *fanOutOptions[INOUT, BUS]) {
		o.fanOutHooks = append(o.fanOutHooks, hooks...)
	}
}

func WithFanInHooks[INOUT any, BUS any](hooks ...FanHook[INOUT, BUS]) FanInOption[INOUT, BUS] {
	return func(o *fanInOptions[INOUT, BUS]) {
		o.fanInHooks = append(o.fanInHooks, hooks...)
	}
}
