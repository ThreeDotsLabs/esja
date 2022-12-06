package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"strings"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	sql2 "github.com/ThreeDotsLabs/esja/pkg/repository/sql"
)

func NewDefaultSimplePostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.NewPostgresConfig[*postcard.Postcard](
			[]aggregate.Event[*postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			},
		),
	)
}

func NewCustomSimplePostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.Config[*postcard.Postcard]{
			SchemaAdapter: sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimple"),
			Serializer: sql2.NewSimpleSerializer(
				sql2.JSONMarshaler{},
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

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.Config[*postcard.Postcard]{
			SchemaAdapter: sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimpleAnonymizing"),
			Serializer: sql2.NewAESAnonymizingSerializer[*postcard.Postcard](
				sql2.NewSimpleSerializer[*postcard.Postcard](
					sql2.JSONMarshaler{},
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

type ConstantSecretProvider struct{}

func (c ConstantSecretProvider) SecretForKey(aggregateID aggregate.ID) ([]byte, error) {
	h, err := hex.DecodeString(strings.ReplaceAll(aggregateID.String(), "-", ""))
	if err != nil {
		return nil, err
	}

	return h, nil
}
