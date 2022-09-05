package batchers

type (
	// Option for an executor.
	//
	// At this moment, no options are provided.
	Option func(*options)

	options struct {
	}
)

func defaultOptions() *options {
	return &options{}
}
