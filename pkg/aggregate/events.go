package aggregate

import "fmt"

// EventName identifies the type of the event and the version of its schema, e.g. "FooCreated_v1".
type EventName string

type Event[A any] interface {
	// EventName should identify the event and the version of its schema.
	//
	// Example implementation:
	// 	func (e FooCreated) EventName() EventName {
	// 		return "FooCreated_v1"
	// 	}
	EventName() EventName

	// Apply applies the event to the aggregate.
	Apply(A) error
}

// VersionedEvent is an event with a corresponding aggregate version.
type VersionedEvent[A any] struct {
	Event[A]
	AggregateVersion int
}

// EventsQueue stores events.
// Zero-value is a valid empty queue.
type EventsQueue[A any] struct {
	version int
	queue   []VersionedEvent[A]
}

// NewEventsQueueFromEvents creates a new EventsQueue with version set to the last event's version.
func NewEventsQueueFromEvents[A any](events []VersionedEvent[A]) (EventsQueue[A], error) {
	if len(events) == 0 {
		return EventsQueue[A]{}, fmt.Errorf("no events to load")
	}

	version := events[len(events)-1].AggregateVersion

	return EventsQueue[A]{
		version: version,
		queue:   events,
	}, nil
}

// Record puts a new Event on the queue with proper version.
func (e *EventsQueue[A]) Record(event Event[A]) {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:            event,
		AggregateVersion: e.version,
	})
}

// PopEvents returns the events on the queue and clears it.
func (e *EventsQueue[A]) PopEvents() []VersionedEvent[A] {
	var tmp = make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}
