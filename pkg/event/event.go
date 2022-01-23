package event

// Name identifies the type of the event and the version of its schema, e.g. "FooCreated_v1".
type Name string

type Event interface {
	// EventName should identify the event and the version of its schema.
	//
	// Example implementation:
	// 	func (e FooCreated) EventName() Name {
	// 		return "FooCreated_v1"
	// 	}
	EventName() Name

	// New creates an empty instance of itself.
	//
	// Example implementation:
	// 	func (e FooCreated) New() Event {
	// 		return FooCreated{}
	// 	}
	New() Event
}
