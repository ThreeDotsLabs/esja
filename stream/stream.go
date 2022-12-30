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
type VersionedEvent[T any] struct {
	Event[T]
	StreamVersion int
}

// Stream stores stream.
// Zero-value is a valid state, ready to use.
type Stream[T any] struct {
	id         ID
	streamType string
	version    int
	queue      []VersionedEvent[T]
}

func NewStream[T any](id ID) (*Stream[T], error) {
	if id == "" {
		return nil, errors.New("empty id")
	}

	return &Stream[T]{
		id: id,
	}, nil
}

func NewStreamWithType[T any](id ID, streamType string) (*Stream[T], error) {
	s, err := NewStream[T](id)
	if err != nil {
		return nil, err
	}

	s.streamType = streamType

	return s, nil
}

func (s *Stream[T]) ID() ID {
	return s.id
}

func (s *Stream[T]) Type() string {
	return s.streamType
}

// Record applies the provided Event to the entity and puts it into the stream's queue with proper version.
func (s *Stream[T]) Record(entity *T, event Event[T]) error {
	err := event.ApplyTo(entity)
	if err != nil {
		return err
	}

	s.version += 1
	s.queue = append(s.queue, VersionedEvent[T]{
		Event:         event,
		StreamVersion: s.version,
	})

	return nil
}

// PopEvents returns the stream on the queue and clears it.
func (s *Stream[T]) PopEvents() []VersionedEvent[T] {
	tmp := make([]VersionedEvent[T], len(s.queue))
	copy(tmp, s.queue)
	s.queue = []VersionedEvent[T]{}

	return tmp
}

// HasEvents returns true if there are any queued stream.
func (s *Stream[T]) HasEvents() bool {
	return len(s.queue) > 0
}

func newStream[T any](id ID, events []VersionedEvent[T]) (*Stream[T], error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("no stream to load")
	}

	e, err := NewStream[T](id)
	if err != nil {
		return nil, err
	}
	e.version = events[len(events)-1].StreamVersion
	e.queue = events

	return e, nil
}
