package repository

import (
	"fmt"
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

func (i InMemoryRepository[T]) Get(id aggregate.ID) (*aggregate.Aggregate[T], error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var (
		t   T
		a   *aggregate.Aggregate[T]
		err error
	)

	events, ok := i.events[id]
	if !ok {
		return a, ErrAggregateNotFound
	}

	for _, event := range events {
		err := t.Handle(event)
		if err != nil {
			return a, fmt.Errorf("error handling event '%s': %w", event.EventName(), err)
		}
	}

	a, err = aggregate.NewAggregate(id, t)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (i *InMemoryRepository[T]) Save(a *aggregate.Aggregate[T]) error {
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
