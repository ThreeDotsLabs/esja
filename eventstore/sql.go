package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/esja"
)

// ContextExecutor can perform SQL queries with context.
type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type storageEvent[A any] struct {
	esja.VersionedEvent[A]
	streamID string
	payload  []byte
}

type schemaAdapter[A any] interface {
	InitializeSchemaQuery() string
	SelectQuery(streamID string) (string, []any, error)
	InsertQuery(streamType string, events []storageEvent[A]) (string, []any, error)
}

// SQLStore is an implementation of the EventStore interface using an SQLStore database.
type SQLStore[T esja.Entity[T]] struct {
	db     ContextExecutor
	config SQLConfig[T]
}

// NewSQLStore creates a new SQL EventStore.
func NewSQLStore[T esja.Entity[T]](
	ctx context.Context,
	db ContextExecutor,
	config SQLConfig[T],
) (SQLStore[T], error) {
	if db == nil {
		return SQLStore[T]{}, errors.New("db must not be nil")
	}

	err := config.validate()
	if err != nil {
		return SQLStore[T]{}, fmt.Errorf("invalid config: %w", err)
	}

	r := SQLStore[T]{
		db:     db,
		config: config,
	}

	err = r.initializeSchema(ctx)
	if err != nil {
		return SQLStore[T]{}, err
	}

	return r, nil
}

func (s SQLStore[T]) initializeSchema(ctx context.Context) error {
	query := s.config.SchemaAdapter.InitializeSchemaQuery()
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error initializing schema: %w", err)
	}
	return nil
}

type event struct {
	streamID      string
	streamVersion int
	streamType    string
	eventName     string
	eventPayload  []byte
}

// Load loads the entity from the database events.
func (s SQLStore[T]) Load(ctx context.Context, id string) (*T, error) {
	query, args, err := s.config.SchemaAdapter.SelectQuery(id)
	if err != nil {
		return nil, fmt.Errorf("error building select query: %w", err)
	}

	results, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving rows for events: %w", err)
	}

	defer func() {
		_ = results.Close()
	}()

	var streamType string

	var dbEvents []event
	for results.Next() {
		e := event{}

		err = results.Scan(&e.streamID, &e.streamVersion, &e.streamType, &e.eventName, &e.eventPayload)
		if err != nil {
			return nil, fmt.Errorf("error reading row result: %w", err)
		}

		if e.streamType != "" {
			streamType = e.streamType
		}

		dbEvents = append(dbEvents, e)
	}

	if len(dbEvents) == 0 {
		return nil, ErrEntityNotFound
	}

	var events []esja.VersionedEvent[T]
	for _, e := range dbEvents {
		event, err := s.config.Mapper.New(e.eventName)
		if err != nil {
			return nil, fmt.Errorf("error creating new event instance: %w", err)
		}

		err = s.config.Marshaler.Unmarshal(e.eventPayload, event)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling event payload: %w", err)
		}

		mappedEvent, err := s.config.Mapper.FromTransport(ctx, e.streamID, event)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}

		events = append(events, esja.VersionedEvent[T]{
			Event:         mappedEvent,
			StreamVersion: e.streamVersion,
		})
	}

	return esja.NewEntityWithStringType(id, streamType, events)
}

// Save saves the entity's queued events to the database.
func (s SQLStore[T]) Save(ctx context.Context, t *T) (err error) {
	if t == nil {
		return errors.New("target to save must not be nil")
	}

	stm := *t

	events := stm.Stream().PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	serializedEvents := make([]storageEvent[T], len(events))
	for i, event := range events {
		mapped, err := s.config.Mapper.ToTransport(ctx, stm.Stream().ID(), event.Event)
		if err != nil {
			return fmt.Errorf("error serializing event: %w", err)
		}

		payload, err := s.config.Marshaler.Marshal(mapped)
		if err != nil {
			return fmt.Errorf("error marshaling event payload: %w", err)
		}

		serializedEvents[i] = storageEvent[T]{
			VersionedEvent: event,
			streamID:       stm.Stream().ID(),
			payload:        payload,
		}
	}

	stmType := stm.Stream().Type()
	query, args, err := s.config.SchemaAdapter.InsertQuery(stmType, serializedEvents)
	if err != nil {
		return fmt.Errorf("error building insert query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
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
