package esja

// Snapshot is an Event that stores and applies the current state back to the Entity.
type Snapshot[T any] interface {
	// Applicable interface requires that each snapshot itself implements
	// the logic how the snapshot data is applied back to the Entity.
	Applicable[T]

	// SnapshotName should identify the snapshot and the version of its schema.
	SnapshotName() string
}

// VersionedSnapshot is a snapshot with a corresponding current stream version.
type VersionedSnapshot[T any] struct {
	Snapshot[T]
	StreamVersion int
}

// EntityWithSnapshots is an optional extension to the Entity interface.
// When implemented it informs that Entity supports snapshots
// and those should be created in the event store when applicable.
type EntityWithSnapshots[T any] interface {
	Entity[T]

	// Snapshot returns a Snapshot representing current the state of the Entity.
	Snapshot() Snapshot[T]
}
