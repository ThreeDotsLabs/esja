package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/eventstore"

	"counter"
)

func TestCounter_Repositories(t *testing.T) {
	testCases := []struct {
		name       string
		repository eventstore.EventStore[counter.Counter]
	}{
		{
			name:       "in_memory",
			repository: eventstore.NewInMemoryStore[counter.Counter](eventstore.InMemoryStoreConfig{MakeSnapshotEveryNVersions: 100}),
		},
	}

	ctx := context.Background()
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.NewString()

			c, err := counter.NewCounter(id)
			assert.NoError(t, err)
			assert.Equal(t, id, c.ID())
			assert.Equal(t, 0, c.CurrentValue())

			_, err = tc.repository.Load(ctx, id)
			assert.ErrorIs(t, err, eventstore.ErrEntityNotFound)

			err = tc.repository.Save(ctx, c)
			require.NoError(t, err)

			fromRepo, err := tc.repository.Load(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, c.ID(), fromRepo.ID())
			assert.Equal(t, c.CurrentValue(), fromRepo.CurrentValue())

			incrementFor := 300
			for i := 0; i < incrementFor; i++ {
				c, err = tc.repository.Load(ctx, id)
				require.NoError(t, err)

				err = c.IncrementBy(1)
				require.NoError(t, err)

				err = tc.repository.Save(ctx, c)
				require.NoError(t, err)
			}

			fromRepo, err = tc.repository.Load(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, c.ID(), fromRepo.ID())
			assert.Equal(t, incrementFor, fromRepo.CurrentValue())
		})
	}
}
