package mysql

import (
	"context"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type SchemaAdapter interface {
	InitializeSchema(ctx context.Context, db beginner) error
	InsertEvents(ctx context.Context, db beginner, events ...event.Event) error
	EventsForAggregate(ctx context.Context, db beginner, id aggregate.ID) ([]event.Event, error)
}
