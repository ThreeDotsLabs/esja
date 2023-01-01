package esja

type Event[T any] interface {
	// EventName should identify the event and the version of its schema.
	//
	// Example:
	//
	// 	func (e FooCreated) EventName() string {
	// 		return "FooCreated_v1"
	// 	}
	EventName() string

	// ApplyTo applies the event to the entity.
	ApplyTo(*T) error
}

// VersionedEvent is an event with a corresponding stream version.
type VersionedEvent[T any] struct {
	Event[T]
	StreamVersion int
}
