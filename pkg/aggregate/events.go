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

// EventStore stores events.
type EventStore[A any] struct {
	aggregate A
	version   int
	queue     []VersionedEvent[A]
}

// NewEventStore creates a new EventStore with empty events queue.
func NewEventStore[A any](aggregate A) EventStore[A] {
	return EventStore[A]{
		aggregate: aggregate,
		version:   0,
		queue:     []VersionedEvent[A]{},
	}
}

// NewEventStoreFromEvents creates a new EventStore from given events.
// The events are applied to the aggregate.
// The version is set to the last event's version.
func NewEventStoreFromEvents[A any](aggregate A, events []VersionedEvent[A]) (EventStore[A], error) {
	if len(events) == 0 {
		return EventStore[A]{}, fmt.Errorf("no events to load")
	}

	for _, ev := range events {
		err := ev.Apply(aggregate)
		if err != nil {
			return EventStore[A]{}, err
		}
	}

	version := events[len(events)-1].AggregateVersion

	return EventStore[A]{
		aggregate: aggregate,
		version:   version,
		queue:     []VersionedEvent[A]{},
	}, nil
}

// Record puts a new Event on the events queue and applies it to the aggregate.
func (e *EventStore[A]) Record(event Event[A]) error {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:            event,
		AggregateVersion: e.version,
	})

	return event.Apply(e.aggregate)
}

// PopEvents returns the events recorded so far and clears the events queue.
func (e *EventStore[A]) PopEvents() []VersionedEvent[A] {
	var tmp = make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}
