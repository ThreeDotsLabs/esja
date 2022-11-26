package aggregate

import (
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

// Aggregate stores events.
// In order for your aggregate to implement EventSourced, make it implement `AggregateID`
// and embed Aggregate, like this:
//
// type MyAggregate struct {
// 		aggregate.Aggregate
//  	id string
// }
//
// func (a MyAggregate) AggregateID() aggregate.ID {
// 		return aggregate.ID(a.id)
// }
//
// Then repository.Repository will be able to store and load it.
type Aggregate struct {
	id          ID
	version     int
	eventsQueue []event.Event
}

type EventSourced interface {
	AggregateID() ID
	PopEvents() []event.Event
}

// PopEvents returns the events recorded so far and clears the events queue.
func (a *Aggregate) PopEvents() []event.Event {
	var tmp = make([]event.Event, len(a.eventsQueue))
	copy(tmp, a.eventsQueue)
	a.eventsQueue = []event.Event{}

	return tmp
}

// RecordEvent puts a new event.Event on the events queue.
func (a *Aggregate) RecordEvent(ev event.Event) {
	a.version += 1
	a.eventsQueue = append(a.eventsQueue, ev)
}

func (a Aggregate) Version() int {
	return a.version
}
