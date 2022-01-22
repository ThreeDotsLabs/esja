package repository

import (
	"errors"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

var (
	ErrAggregateNotFound = errors.New("aggregate not found by ID")
)

type Repository[T aggregate.EventSourced] interface {
	Load(id aggregate.ID, a *aggregate.Aggregate[T]) error
	Save(a aggregate.Aggregate[T]) error
}
