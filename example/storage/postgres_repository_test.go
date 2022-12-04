package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/example/storage"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
	"github.com/ThreeDotsLabs/esja/pkg/repository/inmemory"
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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func TestPostcard_Repositories(t *testing.T) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		repository repository.Repository[*postcard.Postcard]
	}{
		{
			name:       "in_memory",
			repository: inmemory.NewRepository[*postcard.Postcard](),
		},
		{
			name: "postgres_simple",
			repository: func() repository.Repository[*postcard.Postcard] {
				repo, err := storage.NewSimplePostcardRepository(context.Background(), db)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_simple_anonymized",
			repository: func() repository.Repository[*postcard.Postcard] {
				repo, err := storage.NewSimpleAnonymizingPostcardRepository(context.Background(), db)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_mapping",
			repository: func() repository.Repository[*postcard.Postcard] {
				repo, err := storage.NewMappingPostcardRepository(context.Background(), db)
				require.NoError(t, err)
				return repo
			}(),
		},
		{
			name: "postgres_mapping_anonymized",
			repository: func() repository.Repository[*postcard.Postcard] {
				repo, err := storage.NewMappingAnonymizingPostcardRepository(context.Background(), db)
				require.NoError(t, err)
				return repo
			}(),
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.NewString()

			pc, err := postcard.NewPostcard(id)
			assert.NoError(t, err)
			assert.Equal(t, id, pc.ID())

			err = pc.Address(senderAddress, addresseeAddress)
			require.NoError(t, err)

			var fromRepo1 postcard.Postcard
			err = tc.repository.Load(ctx, aggregate.ID(id), &fromRepo1)
			assert.ErrorIs(t, err, repository.ErrAggregateNotFound, "expected aggregate not found yet")

			err = tc.repository.Save(ctx, pc)
			require.NoError(t, err, "should save the aggregate and it has some events already")

			var fromRepo2 postcard.Postcard
			err = tc.repository.Load(ctx, aggregate.ID(id), &fromRepo2)
			assert.NoError(t, err, "should retrieve the aggregate, some events should have been saved")

			var fromRepo2Duplicate postcard.Postcard
			err = tc.repository.Load(ctx, aggregate.ID(id), &fromRepo2Duplicate)
			assert.NoError(t, err)

			assert.Equal(t, pc.ID(), fromRepo2.ID())
			assert.Equal(t, pc.Addressee(), fromRepo2.Addressee())
			assert.Equal(t, pc.Sender(), fromRepo2.Sender())
			assert.Empty(t, fromRepo2.PopEvents())

			err = fromRepo2.Write("content")
			require.NoError(t, err)

			err = fromRepo2.Send()
			require.NoError(t, err)

			err = tc.repository.Save(ctx, &fromRepo2)
			require.NoError(t, err)

			// Another path: send right away without writing
			err = fromRepo2Duplicate.Send()
			require.NoError(t, err)

			err = tc.repository.Save(ctx, &fromRepo2Duplicate)
			require.Error(t, err, "should fail to save the same aggregate version")

			var fromRepo3 postcard.Postcard
			err = tc.repository.Load(ctx, aggregate.ID(id), &fromRepo3)
			assert.NoError(t, err)

			assert.Equal(t, id, fromRepo3.ID())
			assert.Equal(t, senderAddress, fromRepo3.Sender())
			assert.Equal(t, addresseeAddress, fromRepo3.Addressee())
			assert.Equal(t, "content", fromRepo3.Content())
			assert.True(t, fromRepo3.Sent())
			assert.Empty(t, fromRepo3.PopEvents())
		})
	}
}
