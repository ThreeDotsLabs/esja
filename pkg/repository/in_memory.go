package repository

import (
	"context"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
	"sync"
)

type InMemoryRepository[T aggregate.EventSourced] struct {
	lock   sync.RWMutex
	events map[aggregate.ID][]event.Event
}

func NewInMemoryRepository[T aggregate.EventSourced]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		lock:   sync.RWMutex{},
		events: map[aggregate.ID][]event.Event{},
	}
}

func (i InMemoryRepository[T]) Load(_ context.Context, id aggregate.ID, a *aggregate.Aggregate[T]) error {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var (
		t   T
		err error
	)

	events, ok := i.events[id]
	if !ok {
		return ErrAggregateNotFound
	}

	err = a.ApplyEvents(events...)
	if err != nil {
		return err
	}

	a, err = aggregate.NewAggregate(id, t)
	if err != nil {
		return err
	}

	return nil
}

func (i *InMemoryRepository[T]) Save(_ context.Context, a *aggregate.Aggregate[T]) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	events := a.PopEvents()
	if priorEvents, ok := i.events[a.ID()]; !ok {
		i.events[a.ID()] = events
	} else {
		i.events[a.ID()] = append(priorEvents, events...)
	}

	return nil
}
