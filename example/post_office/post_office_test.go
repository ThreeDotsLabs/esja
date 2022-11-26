package post_office_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/ThreeDotsLabs/esja/example/post_office"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
	"github.com/ThreeDotsLabs/esja/pkg/repository/sql"

	stdSQL "database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/example/post_office"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
	"github.com/ThreeDotsLabs/esja/pkg/repository/sql"
)

var (
	senderAddress = post_office.Address{
		Name:  "Alice",
		Line1: "Foo Street 123",
		Line2: "Barville",
	}
	addresseeAddress = post_office.Address{
		Name:  "Bob",
		Line1: "987 Xyz Avenue",
		Line2: "Qux City",
	}
)

func TestPostcard_Lifecycle(t *testing.T) {
	id := uuid.NewString()

	postcard, err := post_office.NewPostcard(id)
	assert.Equal(t, id, postcard.ID())
	assert.NoError(t, err)

	err = postcard.Address(senderAddress, addresseeAddress)
	require.NoError(t, err)

	assert.NotEmpty(t, postcard.Sender())
	assert.NotEmpty(t, postcard.Addressee())
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func TestPostcard_Repositories(t *testing.T) {
	testCases := []struct {
		name       string
		repository repository.Repository[*post_office.Postcard]
	}{
		{
			name:       "in_memory",
			repository: repository.NewInMemoryRepository[*post_office.Postcard](post_office.NewPostcardFromEvents),
		},
		{
			name: "postgres",
			repository: func() repository.Repository[*post_office.Postcard] {
				conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
				db, err := stdSQL.Open("postgres", conn)
				require.NoError(t, err)

				marshaler := repository.NewEventsMarshaler(
					post_office.Created{},
					post_office.Addressed{},
					post_office.Written{},
					post_office.Sent{},
				)

				schemaAdapter := sql.NewPostgresSchemaAdapter("events", marshaler)

				repo, err := sql.NewRepository[*post_office.Postcard](db, schemaAdapter, post_office.NewPostcardFromEvents)
				require.NoError(t, err)

				return repo
			}(),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			id := aggregate.ID(uuid.NewString())

			postcard, err := post_office.NewPostcard(id)
			assert.NoError(t, err)

			assert.Equal(t, id, postcard.ID())

			err = postcard.Address(senderAddress, addresseeAddress)
			require.NoError(t, err)

			ctx := context.Background()

			_, err = tc.repository.Load(ctx, aggregate.ID(id))
			assert.ErrorIs(t, err, repository.ErrAggregateNotFound, "expected aggregate not found yet")

			err = tc.repository.Save(ctx, postcard)
			require.NoError(t, err, "should save the aggregate and it has some events already")

			fromRepo, err := tc.repository.Load(ctx, aggregate.ID(id))
			assert.NoError(t, err, "should retrieve the aggregate, some events should have been saved")

			assert.Equal(t, postcard.ID(), fromRepo.ID())
			assert.Equal(t, postcard.Addressee(), fromRepo.Addressee())
			assert.Empty(t, fromRepo.PopEvents())
		})
	}
}
