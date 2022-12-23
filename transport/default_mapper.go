package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

// Event is a transport model which defines which stream model
// it corresponds to and implements the mapping from- and to- the stream model.
type Event[T any] interface {
	SupportedEvent() stream.Event[T]
	FromEvent(event stream.Event[T]) Event[T]
	ToEvent() stream.Event[T]
}

// DefaultMapper implements an interface of transport.Mapper
// The mapper keeps a list of registered transport.Event defining
// mapping between stream- and transport- layer models.
type DefaultMapper[T any] struct {
	supported map[stream.EventName]Event[T]
}

// NewDefaultMapper returns a new instance of a DefaultMapper.
func NewDefaultMapper[T any](
	supportedEvents []Event[T],
) *DefaultMapper[T] {
	supported := map[stream.EventName]Event[T]{}
	for _, e := range supportedEvents {
		supported[e.SupportedEvent().EventName()] = e
	}

	return &DefaultMapper[T]{
		supported: supported,
	}
}

// RegisterEvent registers new transport model.
// If there was a model corresponding to the same stream event already,
// the old model will be replaced.
func (m *DefaultMapper[T]) RegisterEvent(e Event[T]) {
	m.supported[e.SupportedEvent().EventName()] = e
}

func (m *DefaultMapper[T]) New(name stream.EventName) (any, error) {
	e, err := m.eventForEventName(name)
	if err != nil {
		return nil, err
	}

	return newInstance(e), nil
}

func (m *DefaultMapper[T]) ToStorage(
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

func (m *DefaultMapper[T]) FromStorage(
	_ stream.ID,
	i any,
) (stream.Event[T], error) {
	e, ok := i.(Event[T])
	if !ok {
		return nil, fmt.Errorf("payload does not implement the Event[T] interface")
	}

	return e.ToEvent(), nil
}

func (m *DefaultMapper[T]) eventForEventName(name stream.EventName) (Event[T], error) {
	e, ok := m.supported[name]
	if !ok {
		return nil, fmt.Errorf("unsupported event of name '%s'", name)
	}

	return e, nil
}
