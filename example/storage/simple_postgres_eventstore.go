package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"strings"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/eventstore"
	"github.com/ThreeDotsLabs/esja/pkg/transport"
)

func NewDefaultSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.NewPostgresSQLConfig[*postcard.Postcard](
			[]aggregate.Event[*postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
	)
}

func NewCustomSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimple"),
			Serializer: transport.NewSimpleSerializer(
				transport.JSONMarshaler{},
				[]aggregate.Event[*postcard.Postcard]{
					postcard.Created{},
					postcard.Addressed{},
					postcard.Written{},
					postcard.Sent{},
				},
			),
		},
	)
}

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimpleAnonymizing"),
			Serializer: transport.NewAESAnonymizingSerializer[*postcard.Postcard](
				transport.NewSimpleSerializer[*postcard.Postcard](
					transport.JSONMarshaler{},
					[]aggregate.Event[*postcard.Postcard]{
						postcard.Created{},
						postcard.Addressed{},
						postcard.Written{},
						postcard.Sent{},
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

func (c ConstantSecretProvider) SecretForKey(aggregateID aggregate.ID) ([]byte, error) {
	h, err := hex.DecodeString(strings.ReplaceAll(aggregateID.String(), "-", ""))
	if err != nil {
		return nil, err
	}

	return h, nil
}
