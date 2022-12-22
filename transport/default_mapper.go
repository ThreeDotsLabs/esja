package transport

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

type EventMapper[T any] interface {
	SupportedEvent() stream.Event[T]
	ToStorage(event stream.Event[T]) any
	FromStorage(event any) stream.Event[T]
}

type DefaultMapper[T any] struct {
	mappers map[stream.EventName]EventMapper[T]
}

func NewDefaultMapper[T any](
	mappers []EventMapper[T],
) *DefaultMapper[T] {
	mappersMap := map[stream.EventName]EventMapper[T]{}

	for _, m := range mappers {
		mappersMap[m.SupportedEvent().EventName()] = m
	}

	return &DefaultMapper[T]{
		mappers: mappersMap,
	}
}

func (m *DefaultMapper[T]) RegisterEvent(mapper EventMapper[T]) {
	m.mappers[mapper.SupportedEvent().EventName()] = mapper
}

func (m *DefaultMapper[T]) FromStorage(
	_ stream.ID,
	name stream.EventName,
	i interface{},
) (stream.Event[T], error) {
	mapper, err := m.mapperForEventName(name)
	if err != nil {
		return nil, err
	}

	return mapper.FromStorage(i), nil
}

func (m *DefaultMapper[T]) ToStorage(
	_ stream.ID,
	event stream.Event[T],
) (interface{}, error) {
	mapper, err := m.mapperForEventName(event.EventName())
	if err != nil {
		return nil, err
	}

	return mapper.ToStorage(event), nil
}

func (m *DefaultMapper[T]) mapperForEventName(name stream.EventName) (EventMapper[T], error) {
	for n, mapper := range m.mappers {
		if name == n {
			return mapper, nil
		}
	}

	return nil, fmt.Errorf("no mapper for event %s", name)
}
