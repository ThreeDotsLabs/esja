package transport

import (
	"context"

	"github.com/ThreeDotsLabs/esja"
)

// Mapper translates the event into a serializable transport model.
type Mapper[T any] interface {
	// New returns a new instance of a transport model
	// corresponding to the provided event name.
	New(eventName string) (any, error)

	// FromTransport maps corresponding transport model
	// into an instance of an esja.Event.
	FromTransport(ctx context.Context, eventName string, transportEvent any) (esja.Event[T], error)

	// ToTransport maps an esja.Event into an instance of
	// a corresponding transport model.
	ToTransport(ctx context.Context, eventName string, event esja.Event[T]) (any, error)
}
