package sql

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type EventMapper[T any] interface {
	SupportedEvent() aggregate.Event[T]
	StorageEvent() any
	ToStorage(event aggregate.Event[T]) any
	FromStorage(event any) aggregate.Event[T]
}

type MappingSerializer[T any] struct {
	marshaler Marshaler
	mappers   map[aggregate.EventName]EventMapper[T]
}

func NewMappingSerializer[T any](
	marshaler Marshaler,
	mappers []EventMapper[T],
) *MappingSerializer[T] {
	mappersMap := map[aggregate.EventName]EventMapper[T]{}

	for _, m := range mappers {
		mappersMap[m.SupportedEvent().EventName()] = m
	}

	return &MappingSerializer[T]{
		marshaler: marshaler,
		mappers:   mappersMap,
	}
}

func (m *MappingSerializer[T]) RegisterMapper(mapper EventMapper[T]) {
	m.mappers[mapper.SupportedEvent().EventName()] = mapper
}

func (m *MappingSerializer[T]) Serialize(aggregateID aggregate.ID, event aggregate.Event[T]) ([]byte, error) {
	mapper, err := m.mapperForEventName(event.EventName())
	if err != nil {
		return nil, err
	}

	mappedEvent := mapper.ToStorage(event)

	return m.marshaler.Marshal(aggregateID, mappedEvent)
}

func (m *MappingSerializer[T]) Deserialize(aggregateID aggregate.ID, name aggregate.EventName, payload []byte) (aggregate.Event[T], error) {
	mapper, err := m.mapperForEventName(name)
	if err != nil {
		return nil, err
	}

	event := mapper.StorageEvent()

	err = m.marshaler.Unmarshal(aggregateID, payload, event)
	if err != nil {
		return nil, err
	}

	mappedEvent := mapper.FromStorage(event)

	return mappedEvent, nil
}

func (m *MappingSerializer[T]) mapperForEventName(name aggregate.EventName) (EventMapper[T], error) {
	for n, mapper := range m.mappers {
		if name == n {
			return mapper, nil
		}
	}

	return nil, fmt.Errorf("no mapper for event %s", name)
}
