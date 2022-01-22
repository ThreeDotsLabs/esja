package post_office_test

import (
	"testing"

	"github.com/ThreeDotsLabs/esja/example/post_office"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostcard_Lifecycle(t *testing.T) {
	id := aggregate.ID("123")

	postcard, err := post_office.NewPostcardAggregate(id)
	require.NoError(t, err)

	assert.Equal(t, id, postcard.ID())

	err = postcard.Handle(post_office.Addressed{
		Sender: post_office.Address{
			Name:  "Alice",
			Line1: "Foo Street 123",
			Line2: "Barville",
		},
		Addressee: post_office.Address{
			Name:  "Bob",
			Line1: "987 Xyz Avenue",
			Line2: "Qux City",
		},
	})
	require.NoError(t, err)

	pc := postcard.Base()
	assert.NotEmpty(t, pc.Sender())
	assert.NotEmpty(t, pc.Addressee())
}
