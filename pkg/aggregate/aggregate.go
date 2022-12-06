package aggregate

// Aggregate represents the type saved and loaded by the event store.
//
// In order for your domain type to implement Aggregate:
//   * Embed Events
//   * Implement `AggregateID` returning a unique identifier (usually the same as your aggregate's internal ID).
//   * Implement `FromEvents` to apply events to your aggregate.
//   * Implement `PopEvents` that returns the events on the Events.
//
// Example:
//
//     type MyAggregate struct {
//         es aggregate.Events[*MyAggregate]
//         id string
//     }
//
//     func (a *MyAggregate) AggregateID() aggregate.ID {
//         return aggregate.ID(a.id)
//     }
//
//     func (a *MyAggregate) PopEvents() []aggregate.VersionedEvent[*MyAggregate] {
//         return p.es.PopEvents()
//     }
//
//     func (a *MyAggregate) FromEvents(events []aggregate.VersionedEvent[*MyAggregate]) error {
//         es, err := aggregate.LoadEvents(a, events)
//         if err != nil {
//             return err
//         }
//
//         p.es = es
//
//         return nil
//     }
//
// Then an EventStore will be able to store and load it.
type Aggregate[A any] interface {
	AggregateID() ID
	PopEvents() []VersionedEvent[A]
	FromEvents(eq Events[A]) error
}

// ID is the unique identifier of an aggregate.
type ID string

func (i ID) String() string {
	return string(i)
}

func Record[A any](agg A, eq *Events[A], e Event[A]) error {
	err := e.Apply(agg)
	if err != nil {
		return err
	}

	eq.Record(e)

	return nil
}
