package repository

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

var (
	ErrAggregateNotFound = errors.New("aggregate not found by ID")
)

type Repository[T aggregate.EventSourced] interface {
	Load(ctx context.Context, id aggregate.ID) (T, error)
	Save(ctx context.Context, aggregate T) error
}
