package post_office_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/esja/pkg/repository"
	"github.com/ThreeDotsLabs/esja/example/post_office"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	id := aggregate.ID("123")

	postcard, err := post_office.NewPostcardAggregate(id)
	require.NoError(t, err)

	assert.Equal(t, id, postcard.ID())

	err = postcard.Record(post_office.Addressed{
		Sender:    senderAddress,
		Addressee: addresseeAddress,
	})
	require.NoError(t, err)

	pc := postcard.Base()
	assert.NotEmpty(t, pc.Sender())
	assert.NotEmpty(t, pc.Addressee())
}

func TestPostcard_InMemoryRepository(t *testing.T) {
	id := aggregate.ID("123")

	postcard, err := post_office.NewPostcardAggregate(id)
	require.NoError(t, err)

	assert.Equal(t, id, postcard.ID())

	err = postcard.Record(post_office.Addressed{
		Sender:    senderAddress,
		Addressee: addresseeAddress,
	})
	require.NoError(t, err)

	repo := repository.NewInMemoryRepository[*post_office.Postcard]()
	ctx := context.Background()

	// fromRepo will be the target of repo.Load
	fromRepo, err := post_office.NewPostcardAggregate(id)
	require.NoError(t, err)

	err = repo.Load(ctx, id, fromRepo)
	assert.Equal(t, repository.ErrAggregateNotFound, err)

	err = repo.Save(ctx, postcard)
	require.NoError(t, err)

	err = repo.Load(ctx, id, fromRepo)
	assert.NoError(t, err)

	assert.Equal(t, postcard.ID(), fromRepo.ID())
	assert.Equal(t, postcard.Base().Addressee(), fromRepo.Base().Addressee())
	assert.Empty(t, fromRepo.PopEvents())
}
