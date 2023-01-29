package esja

// Event is a simple Entity event model
type Event[T any] interface {
	// ApplyTo applies the event to the entity.
	ApplyTo(*T) error

	// EventName should identify the event and the version of its schema.
	//
	// Example:
	//
	// 	func (e FooCreated) EventName() string {
	// 		return "FooCreated_v1"
	// 	}
	EventName() string
}

// VersionedEvent is an event with a corresponding stream version.
type VersionedEvent[T any] struct {
	Event[T]
	StreamVersion int
}
