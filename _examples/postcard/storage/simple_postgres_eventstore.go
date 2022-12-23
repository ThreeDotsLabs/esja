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
			"PostcardSimpleDefault",
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
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard]("PostcardSimpleCustom"),
			Mapper: transport.NewNoOpMapper(
				[]stream.Event[postcard.Postcard]{
					postcard.Created{},
					postcard.Addressed{},
					postcard.Written{},
					postcard.Sent{},
				},
			),
			Marshaler: transport.JSONMarshaler{},
		},
	)
}

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard]("PostcardSimpleAnonymized"),
			Mapper: transport.NewAESAnonymizer[postcard.Postcard](
				transport.NewNoOpMapper[postcard.Postcard](
					[]stream.Event[postcard.Postcard]{
						postcard.Created{},
						postcard.Addressed{},
						postcard.Written{},
						postcard.Sent{},
					},
				),
				ConstantSecretProvider{},
			),
			Marshaler: transport.JSONMarshaler{},
		},
	)
}

func NewSimpleSQLitePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewSQLiteConfig[postcard.Postcard](
			"PostcardMappingSQLite",
			[]stream.Event[postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
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
