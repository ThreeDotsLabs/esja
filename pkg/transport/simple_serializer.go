package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type SimpleSerializer[T any] struct {
	marshaler Marshaler
	events    map[aggregate.EventName]aggregate.Event[T]
}

func NewSimpleSerializer[T any](
	marshaler Marshaler,
	supportedEvents []aggregate.Event[T],
) *SimpleSerializer[T] {
	events := make(map[aggregate.EventName]aggregate.Event[T])
	for _, c := range supportedEvents {
		events[c.EventName()] = c
	}
	return &SimpleSerializer[T]{
		marshaler: marshaler,
		events:    events,
	}
}

func (m *SimpleSerializer[T]) Serialize(aggregateID aggregate.ID, event aggregate.Event[T]) ([]byte, error) {
	_, err := m.eventByName(event.EventName())
	if err != nil {
		return nil, err
	}

	return m.marshaler.Marshal(aggregateID, event)
}

func (m *SimpleSerializer[T]) Deserialize(aggregateID aggregate.ID, name aggregate.EventName, payload []byte) (aggregate.Event[T], error) {
	e, err := m.eventByName(name)
	if err != nil {
		return nil, err
	}

	event := reflect.New(reflect.TypeOf(e)).Interface().(aggregate.Event[T])

	err = m.marshaler.Unmarshal(aggregateID, payload, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (m *SimpleSerializer[T]) eventByName(name aggregate.EventName) (aggregate.Event[T], error) {
	for n, event := range m.events {
		if name == n {
			return event, nil
		}
	}

	return nil, fmt.Errorf("no event for event %s", name)
}
