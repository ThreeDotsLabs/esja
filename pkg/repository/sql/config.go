package sql

import "fmt"

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

func NewConfig[T any](schemaAdapter schemaAdapter[T], serializer EventSerializer[T]) Config[T] {
	return Config[T]{
		SchemaAdapter: schemaAdapter,
		Serializer:    serializer,
	}
}
