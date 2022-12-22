package eventstore

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/esja/transport"
)

type SQLConfig[T any] struct {
	SchemaAdapter schemaAdapter[T]
	Mapper        transport.Mapper[T]
	Marshaler     transport.Marshaler
}

func (c SQLConfig[T]) validate() error {
	if c.SchemaAdapter == nil {
		return fmt.Errorf("schema adapter is nil")
	}
	if c.Mapper == nil {
		return fmt.Errorf("mapper is nil")
	}
	if c.Marshaler == nil {
		return fmt.Errorf("marshaler is nil")
	}
	return nil
}

func NewPostgresSQLConfig[T any]() SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Mapper:        transport.NewNoOpMapper[T](),
		Marshaler:     transport.JSONMarshaler{},
	}
}

func NewMappingPostgresSQLConfig[T any](
	eventMappers []transport.EventMapper[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewPostgresSchemaAdapter[T](""),
		Mapper:        transport.NewDefaultMapper[T](eventMappers),
		Marshaler:     transport.JSONMarshaler{},
	}
}

func NewSQLiteConfig[T any](
	supportedEvents []stream.Event[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewSQLiteSchemaAdapter[T](""),
		Mapper:        transport.NewNoOpMapper[T](),
		Marshaler:     transport.JSONMarshaler{},
	}
}

func NewMappingSQLiteConfig[T any](
	eventMappers []transport.EventMapper[T],
) SQLConfig[T] {
	return SQLConfig[T]{
		SchemaAdapter: NewSQLiteSchemaAdapter[T](""),
		Mapper:        transport.NewDefaultMapper[T](eventMappers),
		Marshaler:     transport.JSONMarshaler{},
	}
}
