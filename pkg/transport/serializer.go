package transport

import "github.com/ThreeDotsLabs/esja/pkg/aggregate"

// EventSerializer translates the event into bytes and back.
type EventSerializer[T any] interface {
	Serialize(aggregate.ID, aggregate.Event[T]) ([]byte, error)
	Deserialize(aggregate.ID, aggregate.EventName, []byte) (aggregate.Event[T], error)
}
