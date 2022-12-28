package eventstore

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

const sqliteInitializeSchemaQuery = `
CREATE TABLE IF NOT EXISTS %[1]s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stream_id TEXT NOT NULL,
    stream_version INTEGER NOT NULL,
    stream_type TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_payload BLOB NOT NULL,
    stored_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_stream_id ON %[1]s (stream_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stream_id_version ON %[1]s (stream_id, stream_version);
`

type SQLiteSchemaAdapter[A any] struct {
	streamType string
}

func NewSQLiteSchemaAdapter[A any]() SQLiteSchemaAdapter[A] {
	return SQLiteSchemaAdapter[A]{
		streamType: stream.GetStreamType[A](),
	}
}

func (a SQLiteSchemaAdapter[A]) InitializeSchemaQuery() string {
	return fmt.Sprintf(sqliteInitializeSchemaQuery, defaultEventsTableName)
}

func (a SQLiteSchemaAdapter[A]) SelectQuery(streamID string) (string, []any, error) {
	query := fmt.Sprintf(defaultSelectQuery, defaultEventsTableName)

	args := []any{
		streamID, a.streamType,
	}

	return query, args, nil
}

func (a SQLiteSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
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
