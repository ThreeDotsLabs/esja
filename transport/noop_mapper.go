package transport

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja"
)

// NoOpMapper implements an interface of transport.Mapper
// The mapper will use provided original stream events as transport models.
type NoOpMapper[T any] struct {
	supported map[string]esja.Event[T]
}

// NewNoOpMapper returns a new instance of NoOpMapper.
func NewNoOpMapper[T any](
	supportedEvents []esja.Event[T],
) NoOpMapper[T] {
	supported := make(map[string]esja.Event[T])
	for _, e := range supportedEvents {
		supported[e.EventName()] = e
	}

	return NoOpMapper[T]{
		supported: supported,
	}
}

func (m NoOpMapper[T]) New(eventName string) (any, error) {
	e, ok := m.supported[eventName]
	if !ok {
		return nil, fmt.Errorf("unsupported event of name '%s'", eventName)
	}

	return newInstance(e), nil
}

func (m NoOpMapper[T]) ToTransport(
	_ context.Context,
	_ string,
	event esja.Event[T],
) (any, error) {
	return event, nil
}

func (m NoOpMapper[T]) FromTransport(
	_ context.Context,
	_ string,
	transportEvent any,
) (esja.Event[T], error) {
	event, ok := transportEvent.(esja.Event[T])
	if !ok {
		return nil, fmt.Errorf("transport event does not implement the esja.Event[T] interface")
	}

	return event, nil
}

func newInstance(e any) any {
	return reflect.New(reflect.TypeOf(e)).Interface()
}
