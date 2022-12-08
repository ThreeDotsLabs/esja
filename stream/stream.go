package stream

// Stream represents the type saved and loaded by the event store.
// In DDD terms, it is the "aggregate root".
//
// In order for your domain type to implement Stream:
//   * Embed Events
//   * Implement `StreamID` returning a unique identifier (usually the same as your stream's internal ID).
//   * Implement `PopEvents` that returns the events on the Events.
//   * Implement `FromEvents` to apply events to your stream.
//
// Then an EventStore will be able to store and load it.
//
// Example:
//
//     type User struct {
//         events stream.Events[*User]
//         id string
//     }
//
//     func (u *User) StreamID() stream.ID {
//         return stream.ID(u.id)
//     }
//
//     func (u *User) PopEvents() []stream.VersionedEvent[*User] {
//         return u.events.PopEvents()
//     }
//
//     func (u *User) FromEvents(events stream.Events[*User]) error {
//         p.events = events
//         return stream.ApplyAll(p)
//     }
type Stream[T any] interface {
	StreamID() ID
	PopEvents() []VersionedEvent[T]
	FromEvents(eq Events[T]) error
}

// ID is the unique identifier of a stream.
type ID string

func (i ID) String() string {
	return string(i)
}

func Record[T any](stream T, events *Events[T], e Event[T]) error {
	err := e.ApplyTo(stream)
	if err != nil {
		return err
	}

	events.Record(e)

	return nil
}

func ApplyAll[T Stream[T]](stream T) error {
	for _, e := range stream.PopEvents() {
		err := e.ApplyTo(stream)
		if err != nil {
			return err
		}
	}

	return nil
}
