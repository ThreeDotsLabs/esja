package stream

// Entity represents the type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Entity:
//   - Embed pointer to the Stream.
//   - Implement the interface methods in accordance with its description.
//
// Then an EventStore will be able to store and load it.
//
// Example:
//
//	type User struct {
//	    stream *stream.Stream[User]
//	    id     string
//	}
//
//	func (u User) Stream() *stream.Stream[User] {
//	    return u.stream
//	}
//
//	func (u User) NewWithStream(stream *stream.Stream[User]) *User {
//		return &User{stream: stream}
//	}
type Entity[T any] interface {
	// Stream exposes a pointer to the Stream.
	Stream() *Stream[T]

	// NewWithStream returns a new instance with the provided Stream queue.
	NewWithStream(*Stream[T]) *T
}

// NewEntity instantiates a new T with the given events applied to it.
// At the same time the entity's internal Stream is initialised,
// so it can record new upcoming stream.
func NewEntity[T Entity[T]](id string, eventsSlice []VersionedEvent[T]) (*T, error) {
	var t T

	stream, err := newStream(id, eventsSlice)
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
