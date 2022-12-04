package postcard_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
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

func TestPostcard_Lifecycle(t *testing.T) {
	id := uuid.NewString()

	assert := assert.New(t)

	pc, err := postcard.NewPostcard(id)
	assert.Equal(id, pc.ID())
	assert.NoError(err)

	assert.Empty(pc.Addressee())
	assert.Empty(pc.Sender())
	assert.Empty(pc.Content())
	assert.False(pc.Sent())

	err = pc.Address(senderAddress, addresseeAddress)
	require.NoError(t, err)
	assert.Equal(senderAddress, pc.Sender())
	assert.Equal(addresseeAddress, pc.Addressee())

	err = pc.Write("content")
	require.NoError(t, err)
	assert.Equal("content", pc.Content())

	events := pc.PopEvents()
	assert.Len(events, 3)

	expectedEvents := []aggregate.VersionedEvent[*postcard.Postcard]{
		{Event: &postcard.Created{ID: id}, AggregateVersion: 1},
		{Event: &postcard.Addressed{Sender: senderAddress, Addressee: addresseeAddress}, AggregateVersion: 2},
		{Event: &postcard.Written{Content: "content"}, AggregateVersion: 3},
	}
	assert.Equal(expectedEvents, events)

	pcLoaded := postcard.Postcard{}
	err = pcLoaded.FromEvents(events)
	assert.NoError(err)

	assert.Equal(senderAddress, pcLoaded.Sender())
	assert.Equal(addresseeAddress, pcLoaded.Addressee())
	assert.Equal("content", pcLoaded.Content())
	assert.False(pcLoaded.Sent())

	events = pc.PopEvents()
	assert.Len(events, 0)

	events = pcLoaded.PopEvents()
	assert.Len(events, 0)

	err = pcLoaded.Write("new content")
	require.NoError(t, err)

	err = pcLoaded.Send()
	require.NoError(t, err)
	assert.True(pcLoaded.Sent())

	events = pcLoaded.PopEvents()
	assert.Len(events, 2)

	expectedEvents = []aggregate.VersionedEvent[*postcard.Postcard]{
		{Event: &postcard.Written{Content: "new content"}, AggregateVersion: 4},
		{Event: &postcard.Sent{}, AggregateVersion: 5},
	}

	assert.Equal(expectedEvents, events)
}
