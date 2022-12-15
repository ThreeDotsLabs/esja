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

func (i *InMemoryStore[T]) Load(_ context.Context, id stream.ID) (T, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var t T

	eventsSlice, ok := i.events[id]
	if !ok {
		return t, ErrStreamNotFound
	}

	events, err := stream.NewEvents(eventsSlice)
	if err != nil {
		return t, err
	}

	return stream.New(events)
}

func (i *InMemoryStore[T]) Save(_ context.Context, a T) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	events := a.Events().PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	if priorEvents, ok := i.events[a.StreamID()]; !ok {
		i.events[a.StreamID()] = events
	} else {
		for _, event := range events {
			if len(priorEvents) > 0 {
				if priorEvents[len(priorEvents)-1].StreamVersion >= event.StreamVersion {
					return errors.New("stream version duplicate")
				}
			}
			i.events[a.StreamID()] = append(i.events[a.StreamID()], event)
		}
	}

	return nil
}
