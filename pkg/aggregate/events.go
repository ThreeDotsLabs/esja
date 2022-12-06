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

// Events stores events.
// Zero-value is a valid state, ready to use.
type Events[A any] struct {
	version int
	queue   []VersionedEvent[A]
}

// LoadEvents creates a new Events with version set to the last event's version.
func LoadEvents[A any](events []VersionedEvent[A]) (Events[A], error) {
	if len(events) == 0 {
		return Events[A]{}, fmt.Errorf("no events to load")
	}

	version := events[len(events)-1].AggregateVersion

	return Events[A]{
		version: version,
		queue:   events,
	}, nil
}

// Record puts a new Event on the queue with proper version.
func (e *Events[A]) Record(event Event[A]) {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:            event,
		AggregateVersion: e.version,
	})
}

// PopEvents returns the events on the queue and clears it.
func (e *Events[A]) PopEvents() []VersionedEvent[A] {
	var tmp = make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}
