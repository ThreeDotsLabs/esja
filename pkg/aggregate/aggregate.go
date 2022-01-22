package aggregate

import (
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Aggregate[T EventHandler] struct {
	id          ID
	eh          T
	eventsQueue []event.Event
}

func (a *Aggregate[T]) Handle(e event.Event) error {
	err := a.eh.Handle(e)
	if err != nil {
		return err
	}
	a.eventsQueue = append(a.eventsQueue, e)
	return nil
}

type EventHandler interface {
	Handle(event event.Event) error
}

func NewAggregate[T EventHandler](id ID, eh T) (Aggregate[T], error) {
	var a Aggregate[T]
	if id == "" {
		return a, errors.New("id must not be empty")
	}

	return Aggregate[T]{
		id:          id,
		eh:          eh,
		eventsQueue: []event.Event{},
	}, nil
}

func (a Aggregate[T]) ID() ID {
	return a.id
}

func (a Aggregate[T]) Base() T {
	return a.eh
}
