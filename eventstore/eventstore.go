package eventstore

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/stream"
)

var ErrEntityNotFound = errors.New("entity not found by ID")

// EventStore loads and saves T implementing stream.Entity.
type EventStore[T stream.Entity[T]] interface {
	// Load fetches all events for the ID and returns a new instance of T based on them.
	Load(ctx context.Context, id string) (*T, error)

	// Save saves events recorded in the entity's stream.
	Save(ctx context.Context, entity *T) error
}
