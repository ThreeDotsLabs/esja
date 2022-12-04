package sql

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Marshaler interface {
	Marshal(aggregateID aggregate.ID, data interface{}) ([]byte, error)
	Unmarshal(aggregateID aggregate.ID, bytes []byte, target interface{}) error
}

type JSONMarshaler struct{}

func (JSONMarshaler) Marshal(_ aggregate.ID, data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONMarshaler) Unmarshal(_ aggregate.ID, bytes []byte, target interface{}) error {
	return json.Unmarshal(bytes, target)
}
