package repository

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type InMemoryRepository[T aggregate.EventSourced] struct {
	lock        sync.RWMutex
	events      map[aggregate.ID][]event.Event
	constructor func(events []event.Event) (T, error)
}

func NewInMemoryRepository[T aggregate.EventSourced](constructor func(events []event.Event) (T, error)) *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		lock:        sync.RWMutex{},
		events:      map[aggregate.ID][]event.Event{},
		constructor: constructor,
	}
}

func (i *InMemoryRepository[T]) Load(_ context.Context, id aggregate.ID) (T, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var a T

	events, ok := i.events[id]
	if !ok {
		return a, ErrAggregateNotFound
	}

	return i.constructor(events)
}

func (i *InMemoryRepository[T]) Save(_ context.Context, a T) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	events := a.PopEvents()
	if priorEvents, ok := i.events[a.AggregateID()]; !ok {
		i.events[a.AggregateID()] = events
	} else {
		i.events[a.AggregateID()] = append(priorEvents, events...)
	}

	return nil
}
