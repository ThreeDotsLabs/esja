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

func NewSimplePostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimple"),
		sql2.NewSimpleSerializer(
			sql2.JSONMarshaler{},
			[]aggregate.Event[*postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			}),
	)
}

func NewSimpleAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardSimpleAnonymizing"),
		sql2.NewSimpleSerializer(
			sql2.NewAnonymizingMarshaler(
				sql2.JSONMarshaler{},
				sql2.NewAESAnonymizer(ConstantSecretProvider{}),
			),
			[]aggregate.Event[*postcard.Postcard]{
				postcard.Created{},
				postcard.Addressed{},
				postcard.Written{},
				postcard.Sent{},
			}),
	)
}

type ConstantSecretProvider struct{}

func (c ConstantSecretProvider) SecretForAggregate(aggregateID aggregate.ID) ([]byte, error) {
	h, err := hex.DecodeString(strings.ReplaceAll(aggregateID.String(), "-", ""))
	if err != nil {
		return nil, err
	}

	return h, nil
}
