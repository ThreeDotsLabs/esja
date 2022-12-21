package transport

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/esja/stream"
)

type GOBSerializer[T any] struct {
	events map[stream.EventName]stream.Event[T]
}

func NewGOBSerializer[T any](
	supportedEvents []stream.Event[T],
) GOBSerializer[T] {
	events := make(map[stream.EventName]stream.Event[T])
	for _, c := range supportedEvents {
		events[c.EventName()] = c
	}

	return GOBSerializer[T]{
		events: events,
	}
}

func (s GOBSerializer[T]) Serialize(
	_ stream.ID,
	event stream.Event[T],
) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	e := gob.NewEncoder(b)
	err := e.Encode(event)
	return b.Bytes(), err
}

func (s GOBSerializer[T]) Deserialize(
	_ stream.ID,
	eventName stream.EventName,
	eventBytes []byte,
) (stream.Event[T], error) {
	event, err := s.eventByName(eventName)
	if err != nil {
		return nil, err
	}

	e := reflect.New(reflect.TypeOf(event)).Interface().(stream.Event[T])

	b := bytes.NewBuffer(eventBytes)
	d := gob.NewDecoder(b)
	err = d.Decode(e)
	return e, err
}

func (s GOBSerializer[T]) eventByName(name stream.EventName) (stream.Event[T], error) {
	for n, event := range s.events {
		if name == n {
			return event, nil
		}
	}

	return nil, fmt.Errorf("no event for event %s", name)
}
