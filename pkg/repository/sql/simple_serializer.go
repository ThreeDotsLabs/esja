package sql

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type EventConstructor[T any] func() aggregate.Event[T]

type SimpleSerializer[T any] struct {
	marshaler Marshaler
	events    map[aggregate.EventName]EventConstructor[T]
}

func NewSimpleSerializer[T any](
	marshaler Marshaler,
	constructors []EventConstructor[T],
) *SimpleSerializer[T] {
	events := make(map[aggregate.EventName]EventConstructor[T])
	for _, c := range constructors {
		event := c()
		events[event.EventName()] = c
	}
	return &SimpleSerializer[T]{
		marshaler: marshaler,
		events:    events,
	}
}

func (m *SimpleSerializer[T]) Serialize(aggregateID aggregate.ID, event aggregate.Event[T]) ([]byte, error) {
	_, err := m.constructorForEventName(event.EventName())
	if err != nil {
		return nil, err
	}

	return m.marshaler.Marshal(aggregateID, event)
}

func (m *SimpleSerializer[T]) Deserialize(aggregateID aggregate.ID, name aggregate.EventName, payload []byte) (aggregate.Event[T], error) {
	c, err := m.constructorForEventName(name)
	if err != nil {
		return nil, err
	}

	event := c()
	err = m.marshaler.Unmarshal(aggregateID, payload, event)
	if err != nil {
		return nil, err
	}

	return event.(aggregate.Event[T]), nil
}

func (m *SimpleSerializer[T]) constructorForEventName(name aggregate.EventName) (EventConstructor[T], error) {
	for n, constructor := range m.events {
		if name == n {
			return constructor, nil
		}
	}

	return nil, fmt.Errorf("no constructor for event %s", name)
}
