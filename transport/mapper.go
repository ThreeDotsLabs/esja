package transport

import "github.com/ThreeDotsLabs/esja/stream"

// Mapper translates the event into a serializable transport model.
type Mapper[T any] interface {
	FromStorage(
		stream.ID,
		stream.EventName,
		interface{},
	) (stream.Event[T], error)
	ToStorage(
		stream.ID,
		stream.Event[T],
	) (interface{}, error)
}
