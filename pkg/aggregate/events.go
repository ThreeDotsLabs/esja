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
type EventsQueue[A any] struct {
	aggregate A
	version   int
	queue     []VersionedEvent[A]
}

// NewEventsQueue creates a new empty EventsQueue.
func NewEventsQueue[A any](aggregate A) EventsQueue[A] {
	return EventsQueue[A]{
		aggregate: aggregate,
		version:   0,
		queue:     []VersionedEvent[A]{},
	}
}

// NewEventsQueueFromEvents creates a new EventsQueue from given events.
// The events are applied to the aggregate.
// The version is set to the last event's version.
func NewEventsQueueFromEvents[A any](aggregate A, events []VersionedEvent[A]) (EventsQueue[A], error) {
	if len(events) == 0 {
		return EventsQueue[A]{}, fmt.Errorf("no events to load")
	}

	for _, ev := range events {
		err := ev.Apply(aggregate)
		if err != nil {
			return EventsQueue[A]{}, err
		}
	}

	version := events[len(events)-1].AggregateVersion

	return EventsQueue[A]{
		aggregate: aggregate,
		version:   version,
		queue:     []VersionedEvent[A]{},
	}, nil
}

// PushAndApply puts a new Event on the queue and applies it to the aggregate.
func (e *EventsQueue[A]) PushAndApply(event Event[A]) error {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:            event,
		AggregateVersion: e.version,
	})

	return event.Apply(e.aggregate)
}

// PopEvents returns the events on the queue and clears it.
func (e *EventsQueue[A]) PopEvents() []VersionedEvent[A] {
	var tmp = make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}
