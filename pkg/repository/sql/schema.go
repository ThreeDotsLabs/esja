package sql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/ThreeDotsLabs/esja/pkg/repository"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

// UnmarshalFn provides a way to unmarshal event payloads onto objects corresponding to a particular aggregate.
// The caller of EventsForAggregate is responsible for implementing this.
// payload will be marshaled JSON.
type UnmarshalFn func(id aggregate.ID, eventName string, payload []byte) (event.Event, error)

type SchemaAdapter interface {
	InitializeSchema(ctx context.Context, db database) error
	InsertEvents(ctx context.Context, db database, id aggregate.ID, events ...event.Event) error
	EventsForAggregate(ctx context.Context, db database, id aggregate.ID) ([]event.Event, error)
}

type postgresSchemaAdapter struct {
	EventsTableName string
	Marshaler       repository.EventsMarshaler
}

var tableNameRegexp = regexp.MustCompile("[a-z]+")

func NewPostgresSchemaAdapter(eventsTableName string, marshaler repository.EventsMarshaler) SchemaAdapter {
	if eventsTableName == "" {
		panic(errors.New("empty events table name"))
	}

	if !tableNameRegexp.MatchString(eventsTableName) {
		panic(fmt.Errorf("invalid events table name; use [a-z]+"))
	}

	return postgresSchemaAdapter{
		EventsTableName: eventsTableName,
		Marshaler:       marshaler,
	}
}

func (p postgresSchemaAdapter) InitializeSchema(ctx context.Context, db database) error {
	q := `
CREATE TABLE IF NOT EXISTS %s (
		id serial NOT NULL,
		aggregate_id varchar(36) NOT NULL, -- assuming uuid will be used; if you have a different id, implement your own adapter
		event_name varchar(255) NOT NULL,
		event_payload BYTEA NOT NULL,
		event_payload_debug JSON,
		stored_at TIMESTAMP NOT NULL DEFAULT NOW(),
	
		PRIMARY KEY (id)
);`

	_, err := db.ExecContext(ctx, fmt.Sprintf(q, p.EventsTableName))
	if err != nil {
		return fmt.Errorf("error initializing schema with postgresSchemaAdapter: %w", err)
	}

	q = `CREATE INDEX IF NOT EXISTS idx_aggregate_id ON %s ( aggregate_id );`
	_, err = db.ExecContext(ctx, fmt.Sprintf(q, p.EventsTableName))
	if err != nil {
		return fmt.Errorf("error creating index with postgresSchemaAdapter: %w", err)
	}

	return nil
}

func (p postgresSchemaAdapter) InsertEvents(ctx context.Context, db database, id aggregate.ID, events ...event.Event) (err error) {
	var tx *sql.Tx
	tx, err = db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error beginning tx for inserting events with postgresSchemaAdapter: %w", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	var debugPayload bytes.Buffer
	encoder := json.NewEncoder(&debugPayload)

	q := `
INSERT INTO %s (aggregate_id, event_name, event_payload, event_payload_debug)
VALUES  ($1, $2, $3, $4);
`

	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(q, p.EventsTableName))
	if err != nil {
		return fmt.Errorf("error preparing stmt for inserting events with postgresSchemaAdapter: %w", err)
	}

	for _, e := range events {
		err = encoder.Encode(e)
		if err != nil {
			return fmt.Errorf("error encoding event '%s'\n%+v\n%w", e.EventName(), e, err)
		}

		payload, err := p.Marshaler.Marshal(e)
		if err != nil {
			return err
		}

		result, err := stmt.ExecContext(ctx, id, e.EventName(), payload, debugPayload.Bytes())
		if err != nil {
			return fmt.Errorf("error inserting rows for events with postgresSchemaAdapter: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("insert did not work")
		}

		debugPayload.Reset()
	}

	return nil
}

func (p postgresSchemaAdapter) EventsForAggregate(ctx context.Context, db database, id aggregate.ID) ([]event.Event, error) {
	q := `
SELECT 
	aggregate_id, event_name, event_payload 
FROM "%s"
WHERE aggregate_id = $1 
ORDER BY id ASC;
`
	results, err := db.QueryContext(ctx, fmt.Sprintf(q, p.EventsTableName), string(id))
	if err != nil {
		return nil, fmt.Errorf("error retrieving rows for events with postgresSchemaAdapter: %w", err)
	}

	defer func() {
		_ = results.Close()
	}()

	var (
		aggregateID  aggregate.ID
		eventName    event.Name
		eventPayload []byte
		evts         []event.Event
	)
	for results.Next() {
		err = results.Scan(&aggregateID, &eventName, &eventPayload)
		if err != nil {
			return nil, fmt.Errorf("error reading row result with postgresSchemaAdapter: %w", err)
		}

		evt, err := p.Marshaler.Unmarshal(eventName, eventPayload)
		if err != nil {
			return nil, err
		}

		evts = append(evts, evt)
	}

	if len(evts) == 0 {
		return nil, repository.ErrAggregateNotFound
	}

	return evts, nil
}
