package stream

import "fmt"

// EventName identifies the type of the event and the version of its schema, e.g. "FooCreated_v1".
type EventName string

type Event[T any] interface {
	// EventName should identify the event and the version of its schema.
	//
	// Example implementation:
	// 	func (e FooCreated) EventName() EventName {
	// 		return "FooCreated_v1"
	// 	}
	EventName() EventName

	// ApplyTo applies the event to the stream.
	ApplyTo(T) error
}

// VersionedEvent is an event with a corresponding stream version.
type VersionedEvent[A any] struct {
	Event[A]
	StreamVersion int
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

	version := events[len(events)-1].StreamVersion

	return Events[A]{
		version: version,
		queue:   events,
	}, nil
}

// Record puts a new Event on the queue with proper version.
func (e *Events[A]) Record(event Event[A]) {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:         event,
		StreamVersion: e.version,
	})
}

// PopEvents returns the events on the queue and clears it.
func (e *Events[A]) PopEvents() []VersionedEvent[A] {
	var tmp = make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}

// HasEvents returns true if there are any queued events.
func (e *Events[A]) HasEvents() bool {
	return len(e.queue) > 0
}
