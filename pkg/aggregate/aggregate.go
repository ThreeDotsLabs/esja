package aggregate

import (
	"errors"

	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Aggregate[T any, EH EventHandler[T]] struct {
	id ID
	eh EH
}

type EventHandler[T any] interface {
	Handle(event event.Event) error
	*T
}

func NewAggregate[T any, EH EventHandler[T]](id ID, eh EH) (Aggregate[T, EH], error) {
	var a Aggregate[T, EH]
	if id == "" {
		return a, errors.New("id must not be empty")
	}

	return Aggregate[T, EH]{
		id: id,
		eh: eh,
	}, nil
}

func (a Aggregate[T, EH]) ID() ID {
	return a.id
}
