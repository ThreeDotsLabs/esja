package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

type NoOpMapper[T any] struct {
	supported map[stream.EventName]stream.Event[T]
}

func NewNoOpMapper[T any](
	supportedEvents []stream.Event[T],
) NoOpMapper[T] {
	supported := make(map[stream.EventName]stream.Event[T])
	for _, e := range supportedEvents {
		supported[e.EventName()] = e
	}

	return NoOpMapper[T]{
		supported: supported,
	}
}

func (m NoOpMapper[T]) New(name stream.EventName) (any, error) {
	e, ok := m.supported[name]
	if !ok {
		return nil, fmt.Errorf("unsupported event of name '%s'", name)
	}

	return newInstance(e), nil
}

func (NoOpMapper[T]) ToStorage(
	_ stream.ID,
	event stream.Event[T],
) (any, error) {
	return event, nil
}

func (NoOpMapper[T]) FromStorage(
	_ stream.ID,
	payload any,
) (stream.Event[T], error) {
	event, ok := payload.(stream.Event[T])
	if !ok {
		return nil, fmt.Errorf("payload does not implement the stream.Event[T] interface")
	}

	return event, nil
}

func newInstance(e any) any {
	return reflect.New(reflect.TypeOf(e)).Interface()
}
