package stream

// Stream represents the type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Stream:
//   - Embed pointer to Events queue.
//   - Implement the interface methods in accordance with its description.
//
// Then an EventStore will be able to store and load it.
//
// Example:
//
//	type User struct {
//	    events *stream.Events[User]
//	    id string
//	}
//
//	func (u User) StreamID() stream.ID {
//	    return stream.ID(u.id)
//	}
//
//	func (u User) Events() *stream.Events[User] {
//	    return u.events
//	}
//
//	func (u User) NewFromEvents(events *stream.Events[User]) *User {
//		return &User{events: events}
//	}
type Stream[T any] interface {
	// StreamID returns a unique identifier (usually the same as your stream's internal ID).
	StreamID() ID

	// Events exposes a pointer to the Events queue.
	Events() *Events[T]

	// NewFromEvents returns a new instance with the provided Events queue.
	NewFromEvents(events *Events[T]) *T
}

// ID is the unique identifier of a stream.
type ID string

func (i ID) String() string {
	return string(i)
}

// Record applies a provided Event and puts that into the stream's internal Events queue.
func Record[T Stream[T]](stream *T, e Event[T]) error {
	err := e.ApplyTo(stream)
	if err != nil {
		return err
	}

	(*stream).Events().Record(e)

	return nil
}

// New instantiates a new T with all Events applied to it.
// At the same time the stream's internal Events queue is initialised,
// so it can record new upcoming events.
func New[T Stream[T]](eventsSlice []VersionedEvent[T]) (*T, error) {
	var t T

	events, err := newEvents(eventsSlice)
	if err != nil {
		return nil, err
	}

	eventsSlice = events.PopEvents()

	target := t.NewFromEvents(events)
	for _, e := range eventsSlice {
		err := e.ApplyTo(target)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}
