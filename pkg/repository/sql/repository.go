package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
)

// EventSerializer translates the event to a database-friendly format and back.
type EventSerializer[T any] interface {
	Serialize(aggregate.ID, aggregate.Event[T]) ([]byte, error)
	Deserialize(aggregate.ID, aggregate.EventName, []byte) (aggregate.Event[T], error)
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type storageEvent[A any] struct {
	aggregate.VersionedEvent[A]
	aggregateID string
	payload     []byte
}

type schemaAdapter[A any] interface {
	InitializeSchemaQuery() string
	SelectQuery(aggregateID string) (string, []any, error)
	InsertQuery(events []storageEvent[A]) (string, []any, error)
}

// Repository is an implementation of the Repository interface using an SQL database.
type Repository[T any] struct {
	db            ContextExecutor
	schemaAdapter schemaAdapter[T]
	serializer    EventSerializer[T]
}

// NewRepository creates a new Repository.
// The aggregateType is used to identify the aggregate type in the database. It should be a constant string and not change.
// The serializer is used to translate the events to a database-friendly format and back.
func NewRepository[T any](
	ctx context.Context,
	db ContextExecutor,
	schemaAdapter schemaAdapter[T],
	serializer EventSerializer[T],
) (Repository[T], error) {
	if db == nil {
		return Repository[T]{}, errors.New("db must not be nil")
	}

	r := Repository[T]{
		db:            db,
		schemaAdapter: schemaAdapter,
		serializer:    serializer,
	}

	err := r.initializeSchema(ctx)
	if err != nil {
		return Repository[T]{}, err
	}

	return r, nil
}

func (r Repository[T]) initializeSchema(ctx context.Context) error {
	query := r.schemaAdapter.InitializeSchemaQuery()
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error initializing schema: %w", err)
	}
	return nil
}

// Load loads the aggregate from the database events.
// The target should be a pointer to the aggregate.
func (r Repository[T]) Load(ctx context.Context, id aggregate.ID, target aggregate.Aggregate[T]) error {
	query, args, err := r.schemaAdapter.SelectQuery(id.String())
	if err != nil {
		return fmt.Errorf("error building select query: %w", err)
	}
	results, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error retrieving rows for events: %w", err)
	}

	defer func() {
		_ = results.Close()
	}()

	var (
		aggregateID      aggregate.ID
		aggregateVersion int
		eventName        aggregate.EventName
		eventPayload     []byte
		events           []aggregate.VersionedEvent[T]
	)
	for results.Next() {
		err = results.Scan(&aggregateID, &aggregateVersion, &eventName, &eventPayload)
		if err != nil {
			return fmt.Errorf("error reading row result: %w", err)
		}

		event, err := r.serializer.Deserialize(aggregateID, eventName, eventPayload)
		if err != nil {
			return fmt.Errorf("error deserializing event: %w", err)
		}

		versionedEvent := aggregate.VersionedEvent[T]{
			Event:            event,
			AggregateVersion: aggregateVersion,
		}
		events = append(events, versionedEvent)
	}

	if len(events) == 0 {
		return repository.ErrAggregateNotFound
	}

	return target.FromEvents(events)
}

// Save saves the aggregate's queued events to the database.
func (r Repository[T]) Save(ctx context.Context, agg aggregate.Aggregate[T]) (err error) {
	events := agg.PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	serializedEvents := make([]storageEvent[T], len(events))
	for i, event := range events {
		payload, err := r.serializer.Serialize(agg.AggregateID(), event.Event)
		if err != nil {
			return fmt.Errorf("error serializing event: %w", err)
		}

		serializedEvents[i] = storageEvent[T]{
			VersionedEvent: event,
			aggregateID:    agg.AggregateID().String(),
			payload:        payload,
		}
	}

	query, args, err := r.schemaAdapter.InsertQuery(serializedEvents)
	if err != nil {
		return fmt.Errorf("error building insert query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(events)) {
		return fmt.Errorf("insert did not work")
	}

	return nil
}
