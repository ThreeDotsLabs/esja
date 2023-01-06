package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	sql2 "github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
)

type WatermillSchemaAdapter struct{}

func (s WatermillSchemaAdapter) SchemaInitializingQueries(topic string) []string {
	return nil
}

func (s WatermillSchemaAdapter) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	return "", nil, errors.New("not supported")
}

func (s WatermillSchemaAdapter) SelectQuery(topic string, consumerGroup string, offsetsAdapter sql2.OffsetsAdapter) (string, []interface{}) {
	nextOffsetQuery, nextOffsetArgs := offsetsAdapter.NextOffsetQuery(topic, consumerGroup)
	selectQuery := `
		SELECT id, stream_id, event_payload FROM events
		WHERE
			id > (` + nextOffsetQuery + `)
		ORDER BY
			id ASC
		LIMIT 1`

	return selectQuery, nextOffsetArgs
}

type schemaRow struct {
	ID       int64
	StreamID []byte
	Payload  []byte
}

func (s WatermillSchemaAdapter) UnmarshalMessage(row *sql.Row) (offset int, msg *message.Message, err error) {
	r := schemaRow{}
	err = row.Scan(&r.ID, &r.StreamID, &r.Payload)
	if err != nil {
		return 0, nil, err
	}

	// TODO Event UUID?
	msg = message.NewMessage(uuid.NewString(), r.Payload)

	// TODO Metadata?

	return int(r.ID), msg, nil
}
