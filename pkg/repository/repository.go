package repository

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

var (
	ErrAggregateNotFound = errors.New("aggregate not found by ID")
)

// Repository loads and saves T implementing aggregate.Aggregate.
//
// An example implementation of Repository:
// 1. Load would fetch all events for `AggregateID()` and use them to instantiate a `A` using `FromEvents()` and return it.
// 2. Save would call `PopEvents()` and then save them under the aggregate's id from `AggregateID()`.
type Repository[A any] interface {
	Load(ctx context.Context, id aggregate.ID, target aggregate.Aggregate[A]) error
	Save(ctx context.Context, aggregate aggregate.Aggregate[A]) error
}
