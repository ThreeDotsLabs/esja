package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Aggregate struct{}

type Event struct {
	ID int
}

func (e Event) EventName() aggregate.EventName {
	return "Event"
}

func (e Event) Apply(a Aggregate) error {
	return nil
}

func TestNewEventsQueue(t *testing.T) {
	event1 := Event{ID: 1}
	event2 := Event{ID: 2}

	es := aggregate.EventsQueue[Aggregate]{}

	events := es.PopEvents()
	assert.Len(t, events, 0)

	es.Record(event1)
	es.Record(event2)

	events = es.PopEvents()
	assert.Len(t, events, 2)

	assert.Equal(t, event1, events[0].Event)
	assert.Equal(t, 1, events[0].AggregateVersion)

	assert.Equal(t, event2, events[1].Event)
	assert.Equal(t, 2, events[1].AggregateVersion)

	events = es.PopEvents()
	assert.Len(t, events, 0)

	event3 := Event{ID: 3}

	es.Record(event3)

	events = es.PopEvents()
	assert.Len(t, events, 1)

	assert.Equal(t, event3, events[0].Event)
	assert.Equal(t, 3, events[0].AggregateVersion)
}
