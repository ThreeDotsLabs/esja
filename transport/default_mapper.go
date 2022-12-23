package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

type Event[T any] interface {
	SupportedEvent() stream.Event[T]
	FromEvent(event stream.Event[T]) Event[T]
	ToEvent() stream.Event[T]
}

type DefaultMapper[T any] struct {
	supported map[stream.EventName]Event[T]
}

func NewDefaultMapper[T any](
	supportedEvents []Event[T],
) DefaultMapper[T] {
	supported := map[stream.EventName]Event[T]{}
	for _, e := range supportedEvents {
		supported[e.SupportedEvent().EventName()] = e
	}

	return DefaultMapper[T]{
		supported: supported,
	}
}

func (m DefaultMapper[T]) RegisterEvent(e Event[T]) {
	m.supported[e.SupportedEvent().EventName()] = e
}

func (m DefaultMapper[T]) New(name stream.EventName) (any, error) {
	e, err := m.eventForEventName(name)
	if err != nil {
		return nil, err
	}

	return newInstance(e), nil
}

func (m DefaultMapper[T]) ToStorage(
	_ stream.ID,
	event stream.Event[T],
) (any, error) {
	e, err := m.eventForEventName(event.EventName())
	if err != nil {
		return nil, err
	}

	newEvent := reflect.New(reflect.TypeOf(e)).Interface().(Event[T])
	newEvent = newEvent.FromEvent(event)

	return newEvent, nil
}

func (m DefaultMapper[T]) FromStorage(
	_ stream.ID,
	i any,
) (stream.Event[T], error) {
	e, ok := i.(Event[T])
	if !ok {
		return nil, fmt.Errorf("payload does not implement the Event[T] interface")
	}

	return e.ToEvent(), nil
}

func (m DefaultMapper[T]) eventForEventName(name stream.EventName) (Event[T], error) {
	e, ok := m.supported[name]
	if !ok {
		return nil, fmt.Errorf("unsupported event of name '%s'", name)
	}

	return e, nil
}
