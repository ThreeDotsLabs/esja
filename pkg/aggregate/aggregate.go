package aggregate

import (
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Aggregate struct {
	id          ID
	version     int
	eventsQueue []event.Event
}

type EventSourced interface {
	AggregateID() ID
	PopEvents() []event.Event
}

func NewAggregate(id ID) (*Aggregate, error) {
	return &Aggregate{
		id:          id,
		eventsQueue: []event.Event{},
	}, nil
}

func (a *Aggregate) PopEvents() []event.Event {
	var tmp = make([]event.Event, len(a.eventsQueue))
	copy(tmp, a.eventsQueue)
	a.eventsQueue = []event.Event{}

	return tmp
}

func (a *Aggregate) RecordEvent(ev event.Event) {
	a.eventsQueue = append(a.eventsQueue, ev)
}
