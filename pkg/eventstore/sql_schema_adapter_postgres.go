package eventstore

import (
	"fmt"
)

const postgresInitializeSchemaQuery = `
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

type PostgresSchemaAdapter[A any] struct {
	aggregateType string
}

func NewPostgresSchemaAdapter[A any](aggregateType string) PostgresSchemaAdapter[A] {
	return PostgresSchemaAdapter[A]{
		aggregateType: aggregateType,
	}
}

func (a PostgresSchemaAdapter[A]) InitializeSchemaQuery() string {
	return fmt.Sprintf(postgresInitializeSchemaQuery, defaultEventsTableName)
}

func (a PostgresSchemaAdapter[A]) SelectQuery(aggregateID string) (string, []any, error) {
	query := fmt.Sprintf(defaultSelectQuery, defaultEventsTableName)

	args := []any{
		aggregateID, a.aggregateType,
	}

	return query, args, nil
}

func (a PostgresSchemaAdapter[A]) InsertQuery(events []storageEvent[A]) (string, []any, error) {
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
