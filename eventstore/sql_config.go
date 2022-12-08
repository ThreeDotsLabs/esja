package eventstore

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/esja/transport"
)

type SQLConfig[T any] struct {
	SchemaAdapter schemaAdapter[T]
	Serializer    transport.EventSerializer[T]
}

func (c SQLConfig[T]) validate() error {
	if c.SchemaAdapter == nil {
		return fmt.Errorf("schema adapter is nil")
	}
	if c.Serializer == nil {
		return fmt.Errorf("serializer is nil")
	}
	return nil
}

func NewPostgresSQLConfig[T any](
	supportedEvents []stream.Event[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Serializer:    transport.NewSimpleSerializer(transport.JSONMarshaler{}, supportedEvents),
	}
}

func NewMappingPostgresSQLConfig[T any](
	eventMappers []transport.EventMapper[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Serializer:    transport.NewMappingSerializer[T](transport.JSONMarshaler{}, eventMappers),
	}
}

func NewSQLiteConfig[T any](
	supportedEvents []stream.Event[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewSQLiteSchemaAdapter[T](""),
		Serializer:    transport.NewSimpleSerializer(transport.JSONMarshaler{}, supportedEvents),
	}
}

func NewMappingSQLiteConfig[T any](
	eventMappers []transport.EventMapper[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewSQLiteSchemaAdapter[T](""),
		Serializer:    transport.NewMappingSerializer[T](transport.JSONMarshaler{}, eventMappers),
	}
}
