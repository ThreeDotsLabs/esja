package eventstore

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/esja/stream"
)

type InMemoryStore[T stream.Stream[T]] struct {
	lock   sync.RWMutex
	events map[stream.ID][]stream.VersionedEvent[T]
}

func NewInMemoryStore[T stream.Stream[T]]() *InMemoryStore[T] {
	return &InMemoryStore[T]{
		lock:   sync.RWMutex{},
		events: map[stream.ID][]stream.VersionedEvent[T]{},
	}
}

func (i *InMemoryStore[T]) Load(_ context.Context, id stream.ID) (*T, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	events, ok := i.events[id]
	if !ok {
		return nil, ErrStreamNotFound
	}

	return stream.New(events)
}

func (i *InMemoryStore[T]) Save(_ context.Context, t *T) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if t == nil {
		return errors.New("target to save must not be nil")
	}

	stm := *t

	events := stm.Events().PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	if priorEvents, ok := i.events[stm.StreamID()]; !ok {
		i.events[stm.StreamID()] = events
	} else {
		for _, event := range events {
			if len(priorEvents) > 0 {
				if priorEvents[len(priorEvents)-1].StreamVersion >= event.StreamVersion {
					return errors.New("stream version duplicate")
				}
			}
			i.events[stm.StreamID()] = append(i.events[stm.StreamID()], event)
		}
	}

	return nil
}
