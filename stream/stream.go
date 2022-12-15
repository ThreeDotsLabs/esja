package stream

// Stream represents the type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Stream:
//   - Embed pointer to Events queue.
//   - Implement `StreamID` returning a unique identifier (usually the same as your stream's internal ID).
//   - Implement `Events` to expose the Events queue.
//   - Implement `WithEvents` that returns a new instance with the provided Events queue set.
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
//	func (u User) WithEvents(events *stream.Events[User]) User {
//	    u.events = events
//		return u
//	}
type Stream[T any] interface {
	StreamID() ID
	Events() *Events[T]
	WithEvents(events *Events[T]) T
}

// ID is the unique identifier of a stream.
type ID string

func (i ID) String() string {
	return string(i)
}

func Record[T Stream[T]](stream *T, e Event[T]) error {
	err := e.ApplyTo(stream)
	if err != nil {
		return err
	}

	(*stream).Events().record(e)

	return nil
}

func New[T Stream[T]](eventsSlice []VersionedEvent[T]) (T, error) {
	var t T

	events, err := newEvents(eventsSlice)
	if err != nil {
		return t, err
	}

	eventsSlice = events.PopEvents()
	for _, e := range eventsSlice {
		err := e.ApplyTo(&t)
		if err != nil {
			return t, err
		}
	}

	return t.WithEvents(events), nil
}
