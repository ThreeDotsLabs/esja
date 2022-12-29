package stream

import (
	"errors"
	"fmt"
)

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
	ApplyTo(*T) error
}

// VersionedEvent is an event with a corresponding stream version.
type VersionedEvent[A any] struct {
	Event[A]
	StreamVersion int
}

// Stream stores events.
// Zero-value is a valid state, ready to use.
type Stream[A any] struct {
	id         ID
	streamType string
	metadata   any
	version    int
	queue      []VersionedEvent[A]
}

func NewStream[A any](id ID) (*Stream[A], error) {
	if id == "" {
		return nil, errors.New("empty id")
	}

	return &Stream[A]{
		id: id,
	}, nil
}

func NewStreamWithMetadata[A any](id ID, streamType string, metadata any) (*Stream[A], error) {
	s, err := NewStream[A](id)
	if err != nil {
		return nil, err
	}

	s.streamType = streamType
	s.metadata = metadata

	return s, nil
}

func (e *Stream[A]) ID() ID {
	return e.id
}

func (e *Stream[A]) Type() string {
	return e.streamType
}

// Record puts a new Event on the queue with proper version.
func (e *Stream[A]) Record(event Event[A]) {
	e.version += 1
	e.queue = append(e.queue, VersionedEvent[A]{
		Event:         event,
		StreamVersion: e.version,
	})
}

// PopEvents returns the events on the queue and clears it.
func (e *Stream[A]) PopEvents() []VersionedEvent[A] {
	tmp := make([]VersionedEvent[A], len(e.queue))
	copy(tmp, e.queue)
	e.queue = []VersionedEvent[A]{}

	return tmp
}

// HasEvents returns true if there are any queued events.
func (e *Stream[A]) HasEvents() bool {
	return len(e.queue) > 0
}

func newEvents[A any](id ID, events []VersionedEvent[A]) (*Stream[A], error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("no events to load")
	}

	e, err := NewStream[A](id)
	if err != nil {
		return nil, err
	}
	e.version = events[len(events)-1].StreamVersion
	e.queue = events

	return e, nil
}
