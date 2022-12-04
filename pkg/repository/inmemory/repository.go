package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
)

type Repository[T any] struct {
	lock   sync.RWMutex
	events map[aggregate.ID][]aggregate.VersionedEvent[T]
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{
		lock:   sync.RWMutex{},
		events: map[aggregate.ID][]aggregate.VersionedEvent[T]{},
	}
}

func (i *Repository[T]) Load(_ context.Context, id aggregate.ID, target aggregate.Aggregate[T]) error {
	i.lock.RLock()
	defer i.lock.RUnlock()

	events, ok := i.events[id]
	if !ok {
		return repository.ErrAggregateNotFound
	}

	return target.FromEvents(events)
}

func (i *Repository[T]) Save(_ context.Context, a aggregate.Aggregate[T]) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	events := a.PopEvents()
	if len(events) == 0 {
		return errors.New("no events to save")
	}

	if priorEvents, ok := i.events[a.AggregateID()]; !ok {
		i.events[a.AggregateID()] = events
	} else {
		for _, event := range events {
			if len(priorEvents) > 0 {
				if priorEvents[len(priorEvents)-1].AggregateVersion >= event.AggregateVersion {
					return errors.New("aggregate version duplicate")
				}
			}
			i.events[a.AggregateID()] = append(i.events[a.AggregateID()], event)
		}
	}

	return nil
}
