package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"time"
)

type Repository[T aggregate.EventSourced] struct {
	db            beginner
	schemaAdapter SchemaAdapter
}

func NewRepository[T aggregate.EventSourced](db beginner, schemaAdapter SchemaAdapter) (*Repository, error) {
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

	return &Repository{
		db:            db,
		schemaAdapter: schemaAdapter,
	}, nil
}

func (r Repository[T]) Load(ctx context.Context, id aggregate.ID, a *aggregate.Aggregate[T]) error {
	events, err := r.schemaAdapter.EventsForAggregate(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("error fetching events for aggregate '%s': %w", id, err)
	}

	err = a.ApplyEvents(events...)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository[T]) Save(ctx context.Context, a aggregate.Aggregate[T]) error {
	events := a.PopEvents()

	err := r.schemaAdapter.InsertEvents(ctx, r.db, events...)
	if err != nil {
		return err
	}

	return nil
}
