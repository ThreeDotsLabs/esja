package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type storageEvent[A any] struct {
	stream.VersionedEvent[A]
	streamID string
	payload  []byte
}

type schemaAdapter[A any] interface {
	InitializeSchemaQuery() string
	SelectQuery(streamID string) (string, []any, error)
	InsertQuery(events []storageEvent[A]) (string, []any, error)
}

// SQLStore is an implementation of the EventStore interface using an SQLStore database.
type SQLStore[T stream.Stream[T]] struct {
	db     ContextExecutor
	config SQLConfig[T]
}

// NewSQLStore creates a new SQL EventStore.
// The streamType is used to identify the stream type in the database. It should be a constant string and not change.
// The serializer is used to translate the events to a database-friendly format and back.
func NewSQLStore[T stream.Stream[T]](
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

// Load loads the stream from the database events.
func (s SQLStore[T]) Load(ctx context.Context, id stream.ID) (*T, error) {
	query, args, err := s.config.SchemaAdapter.SelectQuery(id.String())
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

	var (
		streamID      stream.ID
		streamVersion int
		eventName     stream.EventName
		eventPayload  []byte
		events        []stream.VersionedEvent[T]
	)
	for results.Next() {
		err = results.Scan(&streamID, &streamVersion, &eventName, &eventPayload)
		if err != nil {
			return nil, fmt.Errorf("error reading row result: %w", err)
		}

		event, err := s.config.Mapper.New(eventName)
		if err != nil {
			return nil, fmt.Errorf("error creating new event: %w", err)
		}

		err = s.config.Marshaler.Unmarshal(streamID, eventPayload, event)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling event payload: %w", err)
		}

		mappedEvent, err := s.config.Mapper.FromStorage(streamID, event)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}

		events = append(events, stream.VersionedEvent[T]{
			Event:         mappedEvent,
			StreamVersion: streamVersion,
		})
	}

	if len(events) == 0 {
		return nil, ErrStreamNotFound
	}

	return stream.New(events)
}

// Save saves the stream's queued events to the database.
func (s SQLStore[T]) Save(ctx context.Context, t *T) (err error) {
	if t == nil {
		return errors.New("target to save must not be nil")
	}

	stm := *t

	events := stm.Events().PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	serializedEvents := make([]storageEvent[T], len(events))
	for i, event := range events {
		mapped, err := s.config.Mapper.ToStorage(stm.StreamID(), event.Event)
		if err != nil {
			return fmt.Errorf("error serializing event: %w", err)
		}

		payload, err := s.config.Marshaler.Marshal(stm.StreamID(), mapped)
		if err != nil {
			return fmt.Errorf("error marshaling event payload: %w", err)
		}

		serializedEvents[i] = storageEvent[T]{
			VersionedEvent: event,
			streamID:       stm.StreamID().String(),
			payload:        payload,
		}
	}

	query, args, err := s.config.SchemaAdapter.InsertQuery(serializedEvents)
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
