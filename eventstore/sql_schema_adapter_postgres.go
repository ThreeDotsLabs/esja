package eventstore

import (
	"fmt"
)

const postgresInitializeSchemaQuery = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s (
		id serial NOT NULL PRIMARY KEY,
		stream_id uuid NOT NULL, -- assuming uuid will be used; if you have a different id, implement your own adapter
		stream_version int NOT NULL,
		stream_type varchar(255) NOT NULL,
		event_name varchar(255) NOT NULL,
		event_payload JSONB NOT NULL,
		stored_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_stream_id ON %[1]s (stream_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stream_id_version ON %[1]s (stream_id, stream_version);
`

type PostgresSchemaAdapter[A any] struct {
	streamType string
}

func NewPostgresSchemaAdapter[A any](streamType string) PostgresSchemaAdapter[A] {
	return PostgresSchemaAdapter[A]{
		streamType: streamType,
	}
}

func (a PostgresSchemaAdapter[A]) InitializeSchemaQuery() string {
	return fmt.Sprintf(postgresInitializeSchemaQuery, defaultEventsTableName)
}

func (a PostgresSchemaAdapter[A]) SelectQuery(streamID string) (string, []any, error) {
	query := fmt.Sprintf(defaultSelectQuery, defaultEventsTableName)

	args := []any{
		streamID, a.streamType,
	}

	return query, args, nil
}

func (a PostgresSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
	query := fmt.Sprintf(defaultInsertQuery, defaultEventsTableName, defaultInsertMarkers(len(events)))

	var args []any
	for _, e := range events {
		args = append(
			args,
			e.streamID,
			e.StreamVersion,
			a.streamType,
			e.EventName(),
			e.payload,
		)
	}

	return query, args, nil

}
