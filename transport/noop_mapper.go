package transport

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

type NoOpMapper[T any] struct {
}

func NewNoOpMapper[T any]() *NoOpMapper[T] {
	return &NoOpMapper[T]{}
}

func (*NoOpMapper[T]) ToStorage(
	_ stream.ID,
	event stream.Event[T],
) (interface{}, error) {
	return event, nil
}

func (*NoOpMapper[T]) FromStorage(
	_ stream.ID,
	_ stream.EventName,
	payload interface{},
) (stream.Event[T], error) {
	event, ok := payload.(stream.Event[T])
	if !ok {
		return nil, fmt.Errorf("payload does not implement the stream.Event[T] interface")
	}

	return event, nil
}
