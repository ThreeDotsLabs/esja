package esja

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

// NewEntity instantiates a new T with the given events applied to it.
// At the same time the entity's internal Stream is initialised,
// so it can record new upcoming stream.
func NewEntity[T Entity[T]](
	id string,
	eventsSlice []VersionedEvent[T],
) (*T, error) {
	return NewEntityWithStringType(id, "", eventsSlice)
}

// NewEntityWithStringType instantiates a new T with the given
// stream type and events applied to it.
// At the same time the entity's internal Stream is initialised,
// so it can record new upcoming stream.
func NewEntityWithStringType[T Entity[T]](
	id string,
	streamType string,
	eventsSlice []VersionedEvent[T],
) (*T, error) {
	var t T

	stream, err := newStream(
		id,
		streamType,
		eventsSlice,
	)
	if err != nil {
		return nil, err
	}

	eventsSlice = stream.PopEvents()

	target := t.NewWithStream(stream)
	for _, e := range eventsSlice {
		err := e.ApplyTo(target)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}
