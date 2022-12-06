package eventstore

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

var (
	ErrAggregateNotFound = errors.New("aggregate not found by ID")
)

// EventStore loads and saves T implementing aggregate.Aggregate.
//
// An example implementation of EventStore:
// 1. Load would fetch all events for `AggregateID()` and use them to instantiate a `T` using `FromEvents()` and return it.
// 2. Save would call `PopEvents()` and then save them under the aggregate's id from `AggregateID()`.
type EventStore[T aggregate.Aggregate[T]] interface {
	Load(ctx context.Context, id aggregate.ID) (T, error)
	Save(ctx context.Context, aggregate T) error
}
