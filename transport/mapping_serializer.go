package transport

import (
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

type EventMapper[T any] interface {
	SupportedEvent() stream.Event[T]
	StorageEvent() any
	ToStorage(event stream.Event[T]) any
	FromStorage(event any) stream.Event[T]
}

type MappingSerializer[T any] struct {
	marshaler Marshaler
	mappers   map[stream.EventName]EventMapper[T]
}

func NewMappingSerializer[T any](
	marshaler Marshaler,
	mappers []EventMapper[T],
) *MappingSerializer[T] {
	mappersMap := map[stream.EventName]EventMapper[T]{}

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

func (m *MappingSerializer[T]) Serialize(streamID stream.ID, event stream.Event[T]) ([]byte, error) {
	mapper, err := m.mapperForEventName(event.EventName())
	if err != nil {
		return nil, err
	}

	mappedEvent := mapper.ToStorage(event)

	return m.marshaler.Marshal(streamID, mappedEvent)
}

func (m *MappingSerializer[T]) Deserialize(streamID stream.ID, name stream.EventName, payload []byte) (stream.Event[T], error) {
	mapper, err := m.mapperForEventName(name)
	if err != nil {
		return nil, err
	}

	event := mapper.StorageEvent()
	newEvent := reflect.New(reflect.TypeOf(event)).Interface()

	err = m.marshaler.Unmarshal(streamID, payload, &newEvent)
	if err != nil {
		return nil, err
	}

	// Convert from pointer to value
	value := reflect.ValueOf(newEvent).Elem().Interface()
	mappedEvent := mapper.FromStorage(value)

	return mappedEvent, nil
}

func (m *MappingSerializer[T]) mapperForEventName(name stream.EventName) (EventMapper[T], error) {
	for n, mapper := range m.mappers {
		if name == n {
			return mapper, nil
		}
	}

	return nil, fmt.Errorf("no mapper for event %s", name)
}
