package eventstore

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja"
)

var ErrEntityNotFound = errors.New("entity not found by ID")

// EventStore loads and saves T implementing esja.Entity.
type EventStore[T esja.Entity[T]] interface {
	// Load fetches all events for the stream id and returns a new instance of T based on them.
	Load(ctx context.Context, id string) (*T, error)

	// Save saves events recorded in the entity's stream.
	Save(ctx context.Context, entity *T) error
}
