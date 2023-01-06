package esja_test

import (
	"testing"

	"github.com/ThreeDotsLabs/esja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Entity struct {
	stream *esja.Stream[Entity]
}

func (s Entity) Stream() *esja.Stream[Entity] {
	return s.stream
}

func (s Entity) NewWithStream(stream *esja.Stream[Entity]) *Entity {
	return &Entity{stream: stream}
}

type Event struct {
	ID int
}

func (e Event) EventName() string {
	return "Event"
}

func (e Event) ApplyTo(_ *Entity) error {
	return nil
}

func TestNewStream(t *testing.T) {
	var event1 esja.Event[Entity] = Event{ID: 1}
	var event2 esja.Event[Entity] = Event{ID: 2}

	stm, err := esja.NewStreamWithType[Entity]("ID", "Stream")
	require.NoError(t, err)
	assert.Equal(t, "ID", stm.ID())
	assert.Equal(t, "Stream", stm.Type())

	entity := &Entity{
		stream: stm,
	}

	assert.False(t, stm.HasEvents())

	events := stm.PopEvents()
	assert.Len(t, events, 0)

	err = stm.Record(entity, event1)
	assert.NoError(t, err)
	err = stm.Record(entity, event2)
	assert.NoError(t, err)

	assert.True(t, stm.HasEvents())

	events = stm.PopEvents()
	assert.Len(t, events, 2)
	assert.False(t, stm.HasEvents())

	assert.Equal(t, event1, events[0].Event)
	assert.Equal(t, 1, events[0].StreamVersion)

	assert.Equal(t, event2, events[1].Event)
	assert.Equal(t, 2, events[1].StreamVersion)

	events = stm.PopEvents()
	assert.Len(t, events, 0)

	var event3 esja.Event[Entity] = Event{ID: 3}

	err = stm.Record(entity, event3)
	assert.NoError(t, err)

	events = stm.PopEvents()
	assert.Len(t, events, 1)

	assert.Equal(t, event3, events[0].Event)
	assert.Equal(t, 3, events[0].StreamVersion)
}
