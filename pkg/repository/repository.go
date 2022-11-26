package repository

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

var (
	ErrAggregateNotFound = errors.New("aggregate not found by ID")
)

// Repository loads and saves T implementing aggregate.EventSourced.
//
// An example implementation of Repository would:
// 1. Take a T constructor from events like `func(events []event.Event) (T, error)`
// 2. Load would then fetch all events for `id` and use them to instantiate a T and return it
// 3. Save may call aggregate.EventSourced.PopEvents() and then save them under the aggregate's id `T.AggregateID()`.
type Repository[T aggregate.EventSourced] interface {
	Load(ctx context.Context, id aggregate.ID) (T, error)
	Save(ctx context.Context, aggregate T) error
}
