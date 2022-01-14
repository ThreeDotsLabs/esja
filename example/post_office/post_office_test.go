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

}
