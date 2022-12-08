package transport

import "github.com/ThreeDotsLabs/esja/stream"

// EventSerializer translates the event into bytes and back.
type EventSerializer[T any] interface {
	Serialize(stream.ID, stream.Event[T]) ([]byte, error)
	Deserialize(stream.ID, stream.EventName, []byte) (stream.Event[T], error)
}
