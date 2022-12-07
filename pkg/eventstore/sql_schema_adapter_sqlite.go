package eventstore

import (
	"fmt"
)

const sqliteInitializeSchemaQuery = `
CREATE TABLE IF NOT EXISTS %[1]s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    aggregate_id varchar NOT NULL,
    aggregate_version int NOT NULL,
    aggregate_type varchar NOT NULL,
    event_name varchar NOT NULL,
    event_payload blob NOT NULL,
    stored_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_aggregate_id ON %[1]s (aggregate_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_aggregate_id_version ON %[1]s (aggregate_id, aggregate_version);
`

type SQLiteSchemaAdapter[A any] struct {
	aggregateType string
}

func NewSQLiteSchemaAdapter[A any](aggregateType string) SQLiteSchemaAdapter[A] {
	return SQLiteSchemaAdapter[A]{
		aggregateType: aggregateType,
	}
}

func (a SQLiteSchemaAdapter[A]) InitializeSchemaQuery() string {
	return fmt.Sprintf(sqliteInitializeSchemaQuery, defaultEventsTableName)
}

func (a SQLiteSchemaAdapter[A]) SelectQuery(aggregateID string) (string, []any, error) {
	query := fmt.Sprintf(defaultSelectQuery, defaultEventsTableName)

	args := []any{
		aggregateID, a.aggregateType,
	}

	return query, args, nil
}

func (a SQLiteSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
	query := fmt.Sprintf(defaultInsertQuery, defaultEventsTableName, defaultInsertMarkers(len(events)))

	var args []any
	for _, e := range events {
		args = append(
			args,
			e.aggregateID,
			e.AggregateVersion,
			a.aggregateType,
			e.EventName(),
			e.payload,
		)
	}

	return query, args, nil

}
