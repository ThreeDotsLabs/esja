package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/esja/pkg/event"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
)

type Repository[T aggregate.EventSourced] struct {
	db            database
	schemaAdapter SchemaAdapter
	marshaler     repository.EventsMarshaler
	constructor   func(events []event.Event) (T, error)
}

func NewRepository[T aggregate.EventSourced](
	db database,
	schemaAdapter SchemaAdapter,
	constructor func(events []event.Event) (T, error),
) (*Repository[T], error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	if schemaAdapter == nil {
		return nil, errors.New("schema adapter must not be nil")
	}

	// todo: better ideas for this?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := schemaAdapter.InitializeSchema(ctx, db)
	if err != nil {
		return nil, err
	}

	return &Repository[T]{
		db:            db,
		schemaAdapter: schemaAdapter,
		constructor:   constructor,
	}, nil
}

func (r Repository[T]) Load(
	ctx context.Context,
	id aggregate.ID,
) (T, error) {
	var agg T

	events, err := r.schemaAdapter.EventsForAggregate(ctx, r.db, id)
	if err != nil {
		return agg, fmt.Errorf("error fetching events for aggregate '%s': %w", id, err)
	}

	return r.constructor(events)
}

func (r Repository[T]) Save(
	ctx context.Context,
	agg T,
) error {
	events := agg.PopEvents()

	err := r.schemaAdapter.InsertEvents(ctx, r.db, agg.AggregateID(), events...)
	if err != nil {
		return err
	}

	return nil
}
