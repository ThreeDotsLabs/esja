package stream

// Entity represents the type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Entity:
//   - Embed pointer to Stream queue.
//   - Implement the interface methods in accordance with its description.
//
// Then an EventStore will be able to store and load it.
//
// Example:
//
//	type User struct {
//	    events *stream.Stream[User]
//	    id string
//	}
//
//	func (u User) StreamID() stream.ID {
//	    return stream.ID(u.id)
//	}
//
//	func (u User) Stream() *stream.Stream[User] {
//	    return u.events
//	}
//
//	func (u User) NewFromEvents(events *stream.Stream[User]) *User {
//		return &User{events: events}
//	}
type Entity[T any] interface {
	// Stream exposes a pointer to the Stream queue.
	Stream() *Stream[T]

	// NewFromEvents returns a new instance with the provided Stream queue.
	NewFromStream(events *Stream[T]) *T
}

// ID is the unique identifier of a stream.
type ID string

func (i ID) String() string {
	return string(i)
}

// Record applies a provided Event and puts that into the stream's internal Stream queue.
func Record[T Entity[T]](stream *T, e Event[T]) error {
	err := e.ApplyTo(stream)
	if err != nil {
		return err
	}

	(*stream).Stream().Record(e)

	return nil
}

// New instantiates a new T with all events applied to it.
// At the same time the stream's internal Stream queue is initialised,
// so it can record new upcoming events.
func New[T Entity[T]](id ID, eventsSlice []VersionedEvent[T]) (*T, error) {
	var t T

	events, err := newEvents(id, eventsSlice)
	if err != nil {
		return nil, err
	}

	eventsSlice = events.PopEvents()

	target := t.NewFromStream(events)
	for _, e := range eventsSlice {
		err := e.ApplyTo(target)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}
