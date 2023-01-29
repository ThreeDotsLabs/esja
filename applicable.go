package esja

// Applicable defines a model that can be applied to an Entity/
type Applicable[T any] interface {
	// ApplyTo applies the event to the entity.
	ApplyTo(*T) error
}
