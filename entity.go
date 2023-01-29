package esja

import "fmt"

// Entity represents the event-sourced type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Entity:
//   - Keep *Stream in a field.
//   - Implement the interface methods in accordance with its description.
//
// Then an EventStore will be able to store and load it.
//
// Example:
//
//	type User struct {
//	    stream *esja.Stream[User]
//	    id     string
//	}
//
//	func (u User) Stream() *esja.Stream[User] {
//	    return u.stream
//	}
//
//	func (u User) NewWithStream(stream *esja.Stream[User]) *User {
//		return &User{stream: stream}
//	}
type Entity[T any] interface {
	// Stream exposes a pointer to the internal entity's Stream.
	Stream() *Stream[T]

	// NewWithStream returns a new instance of T
	// with the provided Stream queue injected.
	NewWithStream(*Stream[T]) *T
}

// NewEntityWithSnapshot instantiates a new T with the given snapshot and events applied to it.
// At the same time the entity's internal Stream is initialised,
// so it can record new upcoming events.
func NewEntityWithSnapshot[T Entity[T]](
	id string,
	snapshot VersionedSnapshot[T],
	events []VersionedEvent[T],
) (*T, error) {
	var t T

	stream, err := NewStream[T](id)
	if err != nil {
		return nil, err
	}

	stream.queue = events
	stream.version = snapshot.StreamVersion
	if len(events) != 0 {
		stream.version = events[len(events)-1].StreamVersion
	}

	target := t.NewWithStream(stream)

	err = snapshot.ApplyTo(target)
	if err != nil {
		return nil, err
	}

	events = stream.PopEvents()
	for _, e := range events {
		err := e.ApplyTo(target)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}

// NewEntity instantiates a new T with the given events applied to it.
// At the same time the entity's internal Stream is initialised,
// so it can record new upcoming events.
func NewEntity[T Entity[T]](
	id string,
	events []VersionedEvent[T],
) (*T, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("no stream to load")
	}

	var t T

	stream, err := NewStream[T](id)
	if err != nil {
		return nil, err
	}

	stream.queue = events
	stream.version = events[len(events)-1].StreamVersion

	target := t.NewWithStream(stream)

	events = stream.PopEvents()
	for _, e := range events {
		err := e.ApplyTo(target)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}
