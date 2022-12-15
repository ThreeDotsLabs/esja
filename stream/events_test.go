package stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ThreeDotsLabs/esja/stream"
)

type Stream struct {
	events *stream.Events[Stream]
}

func (s Stream) StreamID() stream.ID {
	return "ID"
}

func (s Stream) Events() *stream.Events[Stream] {
	return s.events
}

func (s Stream) WithEvents(events *stream.Events[Stream]) Stream {
	s.events = events
	return s
}

type Event struct {
	ID int
}

func (e Event) EventName() stream.EventName {
	return "Event"
}

func (e Event) ApplyTo(_ *Stream) error {
	return nil
}

func TestNewEventsQueue(t *testing.T) {
	event1 := Event{ID: 1}
	event2 := Event{ID: 2}

	es := new(stream.Events[Stream])
	s := Stream{
		events: es,
	}

	assert.False(t, es.HasEvents())

	events := es.PopEvents()
	assert.Len(t, events, 0)

	err := stream.Record(&s, event1)
	assert.NoError(t, err)
	err = stream.Record(&s, event2)
	assert.NoError(t, err)

	assert.True(t, es.HasEvents())

	events = es.PopEvents()
	assert.Len(t, events, 2)
	assert.False(t, es.HasEvents())

	assert.Equal(t, event1, events[0].Event)
	assert.Equal(t, 1, events[0].StreamVersion)

	assert.Equal(t, event2, events[1].Event)
	assert.Equal(t, 2, events[1].StreamVersion)

	events = es.PopEvents()
	assert.Len(t, events, 0)

	event3 := Event{ID: 3}
	err = stream.Record(&s, event3)
	assert.NoError(t, err)

	events = es.PopEvents()
	assert.Len(t, events, 1)

	assert.Equal(t, event3, events[0].Event)
	assert.Equal(t, 3, events[0].StreamVersion)
}
