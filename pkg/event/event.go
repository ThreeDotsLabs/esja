package event

type Event interface {
	// EventName should identify the event and the version of its schema,
	// e.g. Created_v1.
	EventName() string
}
