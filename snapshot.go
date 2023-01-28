package esja

// Snapshot is an Event that stores and applies the current state back to the Entity.
type Snapshot[T any] interface {
	Event[T]
}

// EntityWithSnapshots is an optional extension to the Entity interface.
// When implemented it informs that Entity supports snapshots
// and those should be created in the event store when applicable.
type EntityWithSnapshots[T any] interface {
	Entity[T]

	// Snapshot returns a Snapshot representing current the state of the Entity.
	Snapshot() Snapshot[T]
}
