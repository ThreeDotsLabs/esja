package transport

import "github.com/ThreeDotsLabs/esja/stream"

// Mapper translates the event into a serializable transport model.
type Mapper[T any] interface {
	// New returns an instance of a transport model
	// corresponding to provided stream.EventName.
	New(stream.EventName) (any, error)

	// FromTransport maps corresponding transport model
	// into an instance of a stream.Event.
	FromTransport(
		stream.ID,
		any,
	) (stream.Event[T], error)

	// ToTransport maps a stream.Event into an instance of
	// a corresponding transport model.
	ToTransport(
		stream.ID,
		stream.Event[T],
	) (any, error)
}
