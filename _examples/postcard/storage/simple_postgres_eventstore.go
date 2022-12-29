package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"strings"

	"github.com/ThreeDotsLabs/esja/eventstore"
	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/esja/transport"

	"postcard"
)

func NewDefaultSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewPostgresSQLConfig[postcard.Postcard](
			[]stream.Event[postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
	)
}

func NewCustomSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard]("PostcardSimple"),
			Serializer: transport.NewSimpleSerializer(
				transport.JSONMarshaler{},
				[]stream.Event[postcard.Postcard]{
					postcard.Created{},
					postcard.Addressed{},
					postcard.Written{},
					postcard.Sent{},
				},
			),
		},
	)
}

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard]("PostcardSimpleAnonymizing"),
			Serializer: transport.NewAESAnonymizingSerializer[postcard.Postcard](
				transport.NewSimpleSerializer[postcard.Postcard](
					transport.JSONMarshaler{},
					[]stream.Event[postcard.Postcard]{
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

func NewSimpleSQLitePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewSQLiteConfig[postcard.Postcard](
			[]stream.Event[postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
	)
}

func NewGOBSQLitePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewSQLiteSchemaAdapter[postcard.Postcard](""),
			Serializer: transport.NewSimpleSerializer[postcard.Postcard](
				transport.GOBMarshaler{},
				[]stream.Event[postcard.Postcard]{
					postcard.Created{},
					postcard.Addressed{},
					postcard.Written{},
					postcard.Sent{},
				},
			),
		},
	)
}

type ConstantSecretProvider struct{}

func (c ConstantSecretProvider) SecretForKey(streamID stream.ID) ([]byte, error) {
	h, err := hex.DecodeString(strings.ReplaceAll(streamID.String(), "-", ""))
	if err != nil {
		return nil, err
	}

	return h, nil
}
