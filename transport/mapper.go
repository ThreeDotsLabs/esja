package transport

import "github.com/ThreeDotsLabs/esja/stream"

// Mapper translates the event into a serializable transport model.
type Mapper[T any] interface {
	// New returns a new instance of a transport model
	// corresponding to the provided event name.
	New(string) (any, error)

	// FromTransport maps corresponding transport model
	// into an instance of a stream.Event.
	FromTransport(
		string,
		any,
	) (stream.Event[T], error)

	// ToTransport maps a stream.Event into an instance of
	// a corresponding transport model.
	ToTransport(
		string,
		stream.Event[T],
	) (any, error)
}
