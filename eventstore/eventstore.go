package eventstore

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/stream"
)

var ErrStreamNotFound = errors.New("stream not found by ID")

// EventStore loads and saves T implementing stream.Entity
type EventStore[T stream.Entity[T]] interface {
	// Load fetches all events for the stream id and returns a new instance of T based on them.
	Load(ctx context.Context, id stream.ID) (*T, error)

	// Save saves events recorded in the entity's stream.
	Save(ctx context.Context, entity *T) error
}
