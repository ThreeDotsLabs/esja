package eventstore

import (
	"fmt"
)

const postgresInitializeSchemaQuery = `
CREATE TABLE IF NOT EXISTS %[1]s (
		id serial NOT NULL PRIMARY KEY,
		stream_id varchar(255) NOT NULL,
		stream_version int NOT NULL,
		stream_type varchar(255) NOT NULL,
		event_name varchar(255) NOT NULL,
		event_payload JSONB NOT NULL,
		stored_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_stream_id ON %[1]s (stream_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stream_id_version ON %[1]s (stream_id, stream_version);
`

type PostgresSchemaAdapter[A any] struct{}

func NewPostgresSchemaAdapter[A any]() PostgresSchemaAdapter[A] {
	return PostgresSchemaAdapter[A]{}
}

func (a PostgresSchemaAdapter[A]) InitializeSchemaQuery() string {
	return fmt.Sprintf(postgresInitializeSchemaQuery, defaultEventsTableName)
}

func (a PostgresSchemaAdapter[A]) SelectQuery(streamID string) (string, []any, error) {
	query := fmt.Sprintf(defaultSelectQuery, defaultEventsTableName)

	args := []any{
		streamID,
	}

	return query, args, nil
}

func (a PostgresSchemaAdapter[A]) InsertQuery(streamType string, events []storageEvent[A]) (string, []any, error) {
	query := fmt.Sprintf(defaultInsertQuery, defaultEventsTableName, defaultInsertMarkers(len(events)))

	var args []any
	for _, e := range events {
		args = append(
			args,
			e.streamID,
			e.StreamVersion,
			streamType,
			e.EventName(),
			e.payload,
		)
	}

	return query, args, nil
}
