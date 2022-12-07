package eventstore

import (
	"fmt"
)

type SQLiteSchemaAdapter[A any] struct {
	aggregateType string
}

func NewSQLiteSchemaAdapter[A any](aggregateType string) SQLiteSchemaAdapter[A] {
	return SQLiteSchemaAdapter[A]{
		aggregateType: aggregateType,
	}
}

func (a SQLiteSchemaAdapter[A]) InitializeSchemaQuery() string {
	query := `
CREATE TABLE IF NOT EXISTS %s (
    id serial int NOT NULL,
    aggregate_id varchar NOT NULL,
    aggregate_version int NOT NULL,
    aggregate_type varchar NOT NULL,
    event_name varchar NOT NULL,
    event_payload blob NOT NULL,
    stored_at datetime NOT NULL, 
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_aggregate_id ON %s (aggregate_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_aggregate_id_version ON %s (aggregate_id, aggregate_version);
`
	return fmt.Sprintf(query, eventsTableName)
}

func (a SQLiteSchemaAdapter[A]) SelectQuery(aggregateID string) (string, []any, error) {
	query := `
SELECT 
	aggregate_id, aggregate_version, event_name, event_payload
FROM %s
WHERE aggregate_id = $1 AND aggregate_type = $2
ORDER BY aggregate_version ASC;
`

	query = fmt.Sprintf(query, eventsTableName)

	args := []any{
		aggregateID, a.aggregateType,
	}

	return query, args, nil
}

func (a SQLiteSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
	query := `
INSERT INTO %s (aggregate_id, aggregate_version, aggregate_type, event_name, event_payload)
VALUES %s
`

	query = fmt.Sprintf(query, eventsTableName, defaultInsertMarkers(len(events)))

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
