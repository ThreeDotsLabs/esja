package eventstore

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/esja"
)

type InMemoryStore[T esja.Entity[T]] struct {
	lock        sync.RWMutex
	events      map[string][]esja.VersionedEvent[T]
	streamTypes map[string]string
}

func NewInMemoryStore[T esja.Entity[T]]() *InMemoryStore[T] {
	return &InMemoryStore[T]{
		lock:        sync.RWMutex{},
		events:      map[string][]esja.VersionedEvent[T]{},
		streamTypes: map[string]string{},
	}
}

func (i *InMemoryStore[T]) Load(_ context.Context, id string) (*T, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	events, ok := i.events[id]
	if !ok {
		return nil, ErrEntityNotFound
	}

	streamType, _ := i.streamTypes[id]

	return esja.NewEntity(id, streamType, events)
}

func (i *InMemoryStore[T]) Save(_ context.Context, t *T) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if t == nil {
		return errors.New("target to save must not be nil")
	}

	stm := *t

	i.streamTypes[stm.Stream().ID()] = stm.Stream().Type()

	events := stm.Stream().PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	if priorEvents, ok := i.events[stm.Stream().ID()]; !ok {
		i.events[stm.Stream().ID()] = events
	} else {
		for _, event := range events {
			if len(priorEvents) > 0 {
				if priorEvents[len(priorEvents)-1].StreamVersion >= event.StreamVersion {
					return errors.New("stream version duplicate")
				}
			}
			i.events[stm.Stream().ID()] = append(i.events[stm.Stream().ID()], event)
		}
	}

	return nil
}
