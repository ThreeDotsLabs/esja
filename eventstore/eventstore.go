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
//
// An example implementation of EventStore:
// 1. Load would fetch all events for `StreamID()` and use them to instantiate a pointer to `T` using `FromEvents()` and return it.
// 2. Save would call `PopEvents()` and then save them under the stream's id from `StreamID()`.
type EventStore[T stream.Stream[T]] interface {
	Load(ctx context.Context, id stream.ID) (*T, error)
	Save(ctx context.Context, stream *T) error
}
