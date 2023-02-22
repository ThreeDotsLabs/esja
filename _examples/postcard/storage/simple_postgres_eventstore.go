package storage

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/ThreeDotsLabs/esja"
	"github.com/ThreeDotsLabs/esja/eventstore"
	"github.com/ThreeDotsLabs/esja/transport"
	"github.com/ThreeDotsLabs/pii"

	"postcard"
)

func NewDefaultSimplePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewPostgresSQLConfig[postcard.Postcard](
			[]esja.Event[postcard.Postcard]{
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
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard](),
			Mapper: transport.NewNoOpMapper(
				[]esja.Event[postcard.Postcard]{
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
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard](),
			Mapper: transport.NewAnonymizer[postcard.Postcard](
				transport.NewNoOpMapper[postcard.Postcard](
					[]esja.Event[postcard.Postcard]{
						postcard.Created{},
						postcard.Addressed{},
						postcard.Written{},
						postcard.Sent{},
					},
				),
				pii.NewStructAnonymizer[string, any](
					pii.NewAESAnonymizer[string](ConstantSecretProvider{}),
				),
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
			[]esja.Event[postcard.Postcard]{
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
			SchemaAdapter: eventstore.NewSQLiteSchemaAdapter[postcard.Postcard](),
			Mapper: transport.NewNoOpMapper[postcard.Postcard](
				[]esja.Event[postcard.Postcard]{
					postcard.Created{},
					postcard.Addressed{},
					postcard.Written{},
					postcard.Sent{},
				},
			),
			Marshaler: transport.GOBMarshaler{},
		},
	)
}

type ConstantSecretProvider struct{}

func (c ConstantSecretProvider) SecretForKey(_ context.Context, id string) ([]byte, error) {
	hash := md5.Sum([]byte(id))
	h, err := hex.DecodeString(hex.EncodeToString(hash[:]))
	if err != nil {
		return nil, err
	}

	return h, nil
}
