package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

type SimpleSerializer[T any] struct {
	marshaler Marshaler
	events    map[stream.EventName]stream.Event[T]
}

func NewSimpleSerializer[T any](
	marshaler Marshaler,
	supportedEvents []stream.Event[T],
) *SimpleSerializer[T] {
	events := make(map[stream.EventName]stream.Event[T])
	for _, c := range supportedEvents {
		events[c.EventName()] = c
	}
	return &SimpleSerializer[T]{
		marshaler: marshaler,
		events:    events,
	}
}

func (m *SimpleSerializer[T]) Serialize(streamID stream.ID, event stream.Event[T]) ([]byte, error) {
	_, err := m.eventByName(event.EventName())
	if err != nil {
		return nil, err
	}

	return m.marshaler.Marshal(streamID, event)
}

func (m *SimpleSerializer[T]) Deserialize(streamID stream.ID, name stream.EventName, payload []byte) (stream.Event[T], error) {
	e, err := m.eventByName(name)
	if err != nil {
		return nil, err
	}

	event := reflect.New(reflect.TypeOf(e)).Interface().(stream.Event[T])

	err = m.marshaler.Unmarshal(streamID, payload, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (m *SimpleSerializer[T]) eventByName(name stream.EventName) (stream.Event[T], error) {
	for n, event := range m.events {
		if name == n {
			return event, nil
		}
	}

	return nil, fmt.Errorf("no event for event %s", name)
}
