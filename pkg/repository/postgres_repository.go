package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

const eventsTableName = "events"

// EventMarshaler translates the event to a database-friendly format and back.
type EventMarshaler[T any] interface {
	SupportedEvent() aggregate.Event[T]
	Marshal(event aggregate.Event[T]) ([]byte, error)
	Unmarshal([]byte) (aggregate.Event[T], error)
}

// PostgresRepository is an opinionated implementation of the Repository interface using PostgreSQL.
// It saves all events in a single table (`events`).
type PostgresRepository[T any] struct {
	db            *sql.DB
	aggregateType string
	marshalers    map[aggregate.EventName]EventMarshaler[T]
}

// NewPostgresRepository creates a new PostgresRepository.
// The aggregateType is used to identify the aggregate type in the database. It should be a constant string and not change.
// The marshalers are used to translate the events to a database-friendly format and back.
func NewPostgresRepository[T any](
	ctx context.Context,
	db *sql.DB,
	aggregateType string,
	marshalers []EventMarshaler[T],
) (PostgresRepository[T], error) {
	if db == nil {
		return PostgresRepository[T]{}, errors.New("db must not be nil")
	}

	marshalersMap := map[aggregate.EventName]EventMarshaler[T]{}
	for _, s := range marshalers {
		marshalersMap[s.SupportedEvent().EventName()] = s
	}

	r := PostgresRepository[T]{
		db:            db,
		aggregateType: aggregateType,
		marshalers:    marshalersMap,
	}

	err := r.initializeSchema(ctx)
	if err != nil {
		return PostgresRepository[T]{}, err
	}

	return r, nil
}

func (r PostgresRepository[T]) initializeSchema(ctx context.Context) error {
	q := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s (
		id serial NOT NULL PRIMARY KEY,
		aggregate_id uuid NOT NULL, -- assuming uuid will be used; if you have a different id, implement your own adapter
		aggregate_version int NOT NULL,
		aggregate_type varchar(255) NOT NULL,
		event_name varchar(255) NOT NULL,
		event_payload JSONB NOT NULL,
		stored_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_aggregate_id ON %[1]s (aggregate_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_aggregate_id_version ON %[1]s (aggregate_id, aggregate_version);
`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(q, eventsTableName))
	if err != nil {
		return fmt.Errorf("error initializing schema: %w", err)
	}

	return nil
}

// Load loads the aggregate from the database events.
// The target should be a pointer to the aggregate.
func (r PostgresRepository[T]) Load(ctx context.Context, id aggregate.ID, target aggregate.Aggregate[T]) error {
	q := `
SELECT 
	aggregate_id, aggregate_version, event_name, event_payload 
FROM "%s"
WHERE aggregate_id = $1 AND aggregate_type = $2
ORDER BY aggregate_version ASC;
`
	results, err := r.db.QueryContext(ctx, fmt.Sprintf(q, eventsTableName), id.String(), r.aggregateType)
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

		marshaler, ok := r.marshalers[eventName]
		if !ok {
			return fmt.Errorf("no marshaler found for event name %s", eventName)
		}

		event, err := marshaler.Unmarshal(eventPayload)
		if err != nil {
			return fmt.Errorf("error unmarshaling event: %w", err)
		}

		versionedEvent := aggregate.VersionedEvent[T]{
			Event:            event,
			AggregateVersion: aggregateVersion,
		}
		events = append(events, versionedEvent)
	}

	if len(events) == 0 {
		return ErrAggregateNotFound
	}

	return target.FromEvents(events)
}

// Save saves the aggregate's queued events to the database.
func (r PostgresRepository[T]) Save(ctx context.Context, agg aggregate.Aggregate[T]) (err error) {
	events := agg.PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	var tx *sql.Tx
	tx, err = r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error beginning tx for inserting events: %w", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	q := `
INSERT INTO %s (aggregate_id, aggregate_version, aggregate_type, event_name, event_payload)
VALUES ($1, $2, $3, $4, $5);`

	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(q, eventsTableName))
	if err != nil {
		return fmt.Errorf("error preparing stmt for inserting events: %w", err)
	}

	for _, e := range events {
		marshaler, ok := r.marshalers[e.EventName()]
		if !ok {
			return fmt.Errorf("could not find marshaler for event %s", e.EventName())
		}

		payload, err := marshaler.Marshal(e.Event)
		if err != nil {
			return fmt.Errorf("error marshaling event: %w", err)
		}

		result, err := stmt.ExecContext(ctx, agg.AggregateID(), e.AggregateVersion, r.aggregateType, e.EventName(), payload)
		if err != nil {
			return fmt.Errorf("error inserting rows for events: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected != 1 {
			return fmt.Errorf("insert did not work")
		}
	}

	return nil
}
