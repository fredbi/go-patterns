package pipelines

type (
	Option func(*options)

	options struct {
		inBuffers  int
		outBuffers int
	}
)

func defaultOptions() *options {
	return &options{}
}

func WithOutputBuffers(channelBuffers int) Option {
	return func(o *options) {
		o.outBuffers = channelBuffers
	}
}

func WithInputBuffers(channelBuffers int) Option {
	return func(o *options) {
		o.inBuffers = channelBuffers
	}
}
