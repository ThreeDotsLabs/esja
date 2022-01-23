package repository

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type NoConstructorRegisteredErr struct {
	EventName event.Name
}

func (n NoConstructorRegisteredErr) Error() string {
	return fmt.Sprintf("constructor not registered for event '%s'; use RegisterEvent to provide one", n.EventName)
}

type EventsMarshaler struct {
	constructors map[event.Name]func() event.Event
}

// NewEventsMarshaler should be called with a list of the events that will be registered for unmarshalling.
func NewEventsMarshaler(events ...event.Event) EventsMarshaler {
	constructors := map[event.Name]func() event.Event{}
	for _, evt := range events {
		constructors[evt.EventName()] = evt.New
		gob.Register(evt)
	}

	return EventsMarshaler{
		constructors: constructors,
	}
}

// Marshal serializes the event into a byte array using gob.
func (m EventsMarshaler) Marshal(e event.Event) ([]byte, error) {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(eventTransport{
		Event: e,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling event '%s': %w", e.EventName(), err)
	}

	return buf.Bytes(), nil
}

// Unmarshal deserializes the event from a byte array using gob.
func (m EventsMarshaler) Unmarshal(eventName event.Name, b []byte) (event.Event, error) {
	constructor, ok := m.constructors[eventName]
	if !ok {
		return nil, NoConstructorRegisteredErr{EventName: eventName}
	}

	newEvent := constructor()
	if !ok {
		return nil, fmt.Errorf("constructor for '%s' returned a value that does not implement event.Event", eventName)
	}
	buf := bytes.NewBuffer(b)

	decoded := eventTransport{
		Event: newEvent,
	}
	err := gob.NewDecoder(buf).Decode(&decoded)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling event '%s': %w", eventName, err)
	}

	return decoded.Event, nil
}

// MarshalDebug serializes the event into a human-readable form. It may be used to browse the events in the storage,
// but for proper serialization it's better to use Marshal for more consistent results.
func (m EventsMarshaler) MarshalDebug(e event.Event) ([]byte, error) {
	return json.MarshalIndent(e, "", "\t")
}

type eventTransport struct {
	Event event.Event
}
