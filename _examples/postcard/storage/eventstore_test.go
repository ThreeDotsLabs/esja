package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/eventstore"

	"postcard"
	"postcard/storage"
)

var (
	senderAddress = postcard.Address{
		Name:  "Alice",
		Line1: "Foo Street 123",
		Line2: "Barville",
	}
	addresseeAddress = postcard.Address{
		Name:  "Bob",
		Line1: "987 Xyz Avenue",
		Line2: "Qux City",
	}
)

func TestPostcard_Repositories(t *testing.T) {
	postgresDB := testPostgresDB(t)
	sqliteDB := testSQLiteDB(t)

	testCases := []struct {
		name       string
		repository eventstore.EventStore[postcard.Postcard]
	}{
		{
			name:       "in_memory",
			repository: eventstore.NewInMemoryStore[postcard.Postcard](),
		},
		{
			name: "postgres_simple",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewDefaultSimplePostcardRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_simple_custom",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewCustomSimplePostcardRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_simple_anonymized",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewSimpleAnonymizingPostcardRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_mapping",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewDefaultMappingPostgresRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_mapping_custom",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewCustomMappingPostcardRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_mapping_anonymized",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewMappingAnonymizingPostcardRepository(context.Background(), postgresDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "sqlite_simple",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewSimpleSQLitePostcardRepository(context.Background(), sqliteDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "sqlite_simple_gob",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewGOBSQLitePostcardRepository(context.Background(), sqliteDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "sqlite_mapping",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewMappingSQLitePostcardRepository(context.Background(), sqliteDB)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "sqlite_mapping_gob",
			repository: func() eventstore.EventStore[postcard.Postcard] {
				repo, err := storage.NewGOBMappingSQLitePostcardRepository(context.Background(), sqliteDB)
				require.NoError(t, err)
				return repo
			}(),
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// A random id with more than 36 chars (UUID is 36 chars long).
			id := gofakeit.Generate("streamID-#############################??????????????????????????")

			pc, err := postcard.NewPostcard(id)
			assert.NoError(t, err)
			assert.Equal(t, id, pc.ID())

			err = pc.Address(senderAddress, addresseeAddress)
			require.NoError(t, err)

			_, err = tc.repository.Load(ctx, id)
			assert.ErrorIs(t, err, eventstore.ErrEntityNotFound, "expected entity not found yet")

			err = tc.repository.Save(ctx, pc)
			require.NoError(t, err, "should save the entity and it has some events already")

			fromRepo2, err := tc.repository.Load(ctx, id)
			assert.NoError(t, err, "should retrieve the entity, some events should have been saved")

			fromRepo2Duplicate, err := tc.repository.Load(ctx, id)
			assert.NoError(t, err)

			assert.Equal(t, pc.ID(), fromRepo2.ID())
			assert.Equal(t, pc.Addressee(), fromRepo2.Addressee())
			assert.Equal(t, pc.Sender(), fromRepo2.Sender())
			assert.Empty(t, fromRepo2.Stream().PopEvents())

			err = fromRepo2.Write("content")
			require.NoError(t, err)

			err = fromRepo2.Send()
			require.NoError(t, err)

			err = tc.repository.Save(ctx, fromRepo2)
			require.NoError(t, err)

			// Another path: send right away without writing
			err = fromRepo2Duplicate.Send()
			require.NoError(t, err)

			err = tc.repository.Save(ctx, fromRepo2Duplicate)
			require.Error(t, err, "should fail to save the same entity version")

			fromRepo3, err := tc.repository.Load(ctx, id)
			assert.NoError(t, err)

			assert.Equal(t, id, fromRepo3.ID())
			assert.Equal(t, senderAddress, fromRepo3.Sender())
			assert.Equal(t, addresseeAddress, fromRepo3.Addressee())
			assert.Equal(t, "content", fromRepo3.Content())
			assert.True(t, fromRepo3.Sent())
			assert.Empty(t, fromRepo3.Stream().PopEvents())
		})
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func testPostgresDB(t *testing.T) *sql.DB {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	postgresDB, err := sql.Open("postgres", conn)
	require.NoError(t, err)

	return postgresDB
}

func testSQLiteDB(t *testing.T) *sql.DB {
	dbFile, err := os.CreateTemp("", "tmp_*.db")
	require.NoError(t, err)

	sqliteDB, err := sql.Open("sqlite3", dbFile.Name())
	require.NoError(t, err)

	return sqliteDB
}
