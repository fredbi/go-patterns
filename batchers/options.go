package batchers

type (
	// Option for an executor.
	//
	// At this moment, no options are provided.
	Option func(*options)

	options struct {
		// TODO
	}
)

func defaultOptions() *options {
	return &options{}
}
