package eventstore

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/stream"
)

var (
	ErrStreamNotFound = errors.New("stream not found by ID")
)

// EventStore loads and saves T implementing stream.Stream
type EventStore[T stream.Stream[T]] interface {
	// Load will fetch all events for `StreamID()` and use them
	// to instantiate a pointer to `T` using `FromEvents()` and return it.
	Load(ctx context.Context, id stream.ID) (*T, error)

	// Save will call `PopEvents()` and then save them
	// under the stream's id from `StreamID()`.
	Save(ctx context.Context, stream *T) error
}
