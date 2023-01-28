package eventstore

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/esja"
)

type InMemoryStoreConfig struct {
	// MakeSnapshotEveryNVersions configures a frequency of snapshot creation
	// Once the current event version and last snapshot version difference exceeds the value,
	// a new snapshot will be created with a version of the current event version.
	MakeSnapshotEveryNVersions int
}

type InMemoryStore[T esja.Entity[T]] struct {
	lock      sync.RWMutex
	events    map[string][]esja.VersionedEvent[T]
	snapshots map[string][]esja.VersionedEvent[T]
	config    InMemoryStoreConfig
}

func NewInMemoryStore[T esja.Entity[T]](config InMemoryStoreConfig) *InMemoryStore[T] {
	return &InMemoryStore[T]{
		lock:      sync.RWMutex{},
		events:    map[string][]esja.VersionedEvent[T]{},
		snapshots: map[string][]esja.VersionedEvent[T]{},
		config:    config,
	}
}

func (i *InMemoryStore[T]) Load(_ context.Context, id string) (*T, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	// In the other databases this could be optimized
	// as we do not need to load events of version lower than the lastSnapshot version.
	events, ok := i.events[id]
	if !ok {
		return nil, ErrEntityNotFound
	}

	var eventsToApply []esja.VersionedEvent[T]

	s, found := i.loadLastSnapshot(id)
	if found {
		eventsToApply = append(eventsToApply, s)
	}

	for _, e := range events {
		if e.StreamVersion > s.StreamVersion {
			eventsToApply = append(eventsToApply, e)
		}
	}

	return esja.NewEntity(id, eventsToApply)
}

func (i *InMemoryStore[T]) Save(_ context.Context, t *T) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if t == nil {
		return errors.New("target to save must not be nil")
	}

	entity := *t
	currentVersion, err := i.storeEntityEvents(entity)
	if err != nil {
		return err
	}

	entityWithSnapshots, ok := supportsSnapshots(t)
	if !ok {
		return nil
	}

	err = i.storeEntitySnapshot(entityWithSnapshots, currentVersion)
	if err != nil {
		return err
	}

	return nil
}

func (i *InMemoryStore[T]) loadLastSnapshot(id string) (esja.VersionedEvent[T], bool) {
	snapshots, found := i.snapshots[id]
	var lastSnapshot esja.VersionedEvent[T]
	for _, s := range snapshots {
		if s.StreamVersion >= lastSnapshot.StreamVersion {
			lastSnapshot = s
		}
	}
	return lastSnapshot, found
}

func (i *InMemoryStore[T]) storeEntityEvents(entity T) (int, error) {
	events := entity.Stream().PopEvents()
	if len(events) == 0 {
		return 0, errors.New("no events to save")
	}

	priorEvents, ok := i.events[entity.Stream().ID()]
	if !ok {
		i.events[entity.Stream().ID()] = make([]esja.VersionedEvent[T], 0)
	}

	lastVersion := 0
	for _, e := range priorEvents {
		if e.StreamVersion > lastVersion {
			lastVersion = e.StreamVersion
		}
	}

	currentVersion := lastVersion
	for _, e := range events {
		if e.StreamVersion <= currentVersion {
			return 0, errors.New("stream version duplicate")
		}
		if e.StreamVersion > currentVersion {
			currentVersion = e.StreamVersion
		}

		i.events[entity.Stream().ID()] = append(
			i.events[entity.Stream().ID()],
			e,
		)
	}

	return currentVersion, nil
}

func (i *InMemoryStore[T]) storeEntitySnapshot(
	entity esja.EntityWithSnapshots[T],
	currentVersion int,
) error {
	if i.config.MakeSnapshotEveryNVersions <= 0 {
		return nil
	}

	lastSnapshot, found := i.loadLastSnapshot(entity.Stream().ID())
	lastSnapshotVersion := 0
	if found {
		lastSnapshotVersion = lastSnapshot.StreamVersion
	}

	if currentVersion-lastSnapshotVersion < i.config.MakeSnapshotEveryNVersions {
		return nil
	}

	snapshot := entity.Snapshot()
	snapshotVersioned := esja.VersionedEvent[T]{
		Event:         esja.Event[T](snapshot),
		StreamVersion: currentVersion,
	}

	_, ok := i.snapshots[entity.Stream().ID()]
	if !ok {
		i.snapshots[entity.Stream().ID()] = make([]esja.VersionedEvent[T], 0)
	}

	i.snapshots[entity.Stream().ID()] = append(
		i.snapshots[entity.Stream().ID()],
		snapshotVersioned,
	)

	return nil
}

func supportsSnapshots[T esja.Entity[T]](t *T) (esja.EntityWithSnapshots[T], bool) {
	var entity interface{}
	entity = *t
	entityWithSnapshots, ok := entity.(esja.EntityWithSnapshots[T])
	return entityWithSnapshots, ok
}
