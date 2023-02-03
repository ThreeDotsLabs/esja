package esja

import (
	"errors"
	"fmt"
)

// Stream represents a queue of events and basic stream properties.
type Stream[T any] struct {
	id         string
	streamType string
	version    int
	queue      []VersionedEvent[T]
}

// NewStream creates a new instance of a Stream with provided ID.
func NewStream[T any](id string) (*Stream[T], error) {
	if id == "" {
		return nil, errors.New("empty id")
	}

	return &Stream[T]{
		id: id,
	}, nil
}

// NewStreamWithType creates a new instance of a Stream with provided ID and custom type.
func NewStreamWithType[T any](id string, streamType string) (*Stream[T], error) {
	s, err := NewStream[T](id)
	if err != nil {
		return nil, err
	}

	s.streamType = streamType

	return s, nil
}

func (s *Stream[T]) ID() string {
	return s.id
}

func (s *Stream[T]) Type() string {
	return s.streamType
}

// Record applies the provided Event to the entity
// and puts it into the stream's event queue as a next VersionedEvent.
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

// PopEvents returns the slice of queued VersionedEvents and clears it.
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

func newStream[T any](id string, st string, events []VersionedEvent[T]) (*Stream[T], error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("no stream to load")
	}

	var e *Stream[T]
	var err error
	if st == "" {
		e, err = NewStream[T](id)
	} else {
		e, err = NewStreamWithType[T](id, st)
	}
	if err != nil {
		return nil, err
	}

	e.version = events[len(events)-1].StreamVersion
	e.queue = events

	return e, nil
}
