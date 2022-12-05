package sql

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Config[T any] struct {
	SchemaAdapter schemaAdapter[T]
	Serializer    EventSerializer[T]
}

func (c Config[T]) validate() error {
	if c.SchemaAdapter == nil {
		return fmt.Errorf("schema adapter is nil")
	}
	if c.Serializer == nil {
		return fmt.Errorf("serializer is nil")
	}
	return nil
}

func NewPostgresConfig[T any](
	supportedEvents []aggregate.Event[T],
) Config[T] {
	return Config[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Serializer:    NewSimpleSerializer(JSONMarshaler{}, supportedEvents),
	}
}

func NewMappingPostgresConfig[T any](
	eventMappers []EventMapper[T],
) Config[T] {
	return Config[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Serializer:    NewMappingSerializer[T](JSONMarshaler{}, eventMappers),
	}
}
