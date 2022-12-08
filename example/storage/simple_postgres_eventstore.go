package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"github.com/ThreeDotsLabs/esja/stream"
	"strings"

	postcard2 "github.com/ThreeDotsLabs/esja/example/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/eventstore"
	"github.com/ThreeDotsLabs/esja/pkg/transport"
)

func NewDefaultSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.NewPostgresSQLConfig[*postcard2.Postcard](
			[]stream.Event[*postcard2.Postcard]{
				postcard2.Created{},
				postcard2.Addressed{},
				postcard2.Written{},
				postcard2.Sent{},
			},
		),
	)
}

func NewCustomSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard2.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard2.Postcard]("PostcardSimple"),
			Serializer: transport.NewSimpleSerializer(
				transport.JSONMarshaler{},
				[]stream.Event[*postcard2.Postcard]{
					postcard2.Created{},
					postcard2.Addressed{},
					postcard2.Written{},
					postcard2.Sent{},
				},
			),
		},
	)
}

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard2.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard2.Postcard]("PostcardSimpleAnonymizing"),
			Serializer: transport.NewAESAnonymizingSerializer[*postcard2.Postcard](
				transport.NewSimpleSerializer[*postcard2.Postcard](
					transport.JSONMarshaler{},
					[]stream.Event[*postcard2.Postcard]{
						postcard2.Created{},
						postcard2.Addressed{},
						postcard2.Written{},
						postcard2.Sent{},
					},
				),
				ConstantSecretProvider{},
			),
		},
	)
}

func NewSimpleSQLitePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.NewSQLiteConfig[*postcard.Postcard](
			[]aggregate.Event[*postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
	)
}

type ConstantSecretProvider struct{}

func (c ConstantSecretProvider) SecretForKey(aggregateID stream.ID) ([]byte, error) {
	h, err := hex.DecodeString(strings.ReplaceAll(aggregateID.String(), "-", ""))
	if err != nil {
		return nil, err
	}

	return h, nil
}
