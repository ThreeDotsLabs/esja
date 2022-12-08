package eventstore

import (
	"fmt"
	"strings"
)

const eventsTableName = "events"

type PostgresSchemaAdapter[A any] struct {
	streamType string
}

func NewPostgresSchemaAdapter[A any](streamType string) PostgresSchemaAdapter[A] {
	return PostgresSchemaAdapter[A]{
		streamType: streamType,
	}
}

func (a PostgresSchemaAdapter[A]) InitializeSchemaQuery() string {
	query := `
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
	return fmt.Sprintf(query, eventsTableName)
}

func (a PostgresSchemaAdapter[A]) SelectQuery(streamID string) (string, []any, error) {
	query := `
SELECT 
	stream_id, stream_version, event_name, event_payload 
FROM "%s"
WHERE stream_id = $1 AND stream_type = $2
ORDER BY stream_version ASC;
`

	query = fmt.Sprintf(query, eventsTableName)

	args := []any{
		streamID, a.streamType,
	}

	return query, args, nil
}

func (a PostgresSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
	query := `
INSERT INTO %s (stream_id, stream_version, stream_type, event_name, event_payload)
VALUES %s`

	query = fmt.Sprintf(query, eventsTableName, defaultInsertMarkers(len(events)))

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

func defaultInsertMarkers(count int) string {
	result := strings.Builder{}

	index := 1
	for i := 0; i < count; i++ {
		result.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", index, index+1, index+2, index+3, index+4))
		index += 5
	}

	return strings.TrimRight(result.String(), ",")
}
