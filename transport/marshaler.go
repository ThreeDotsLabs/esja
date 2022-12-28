package transport

import (
	"encoding/json"
)

type Marshaler interface {
	Marshal(data interface{}) ([]byte, error)
	Unmarshal(bytes []byte, target interface{}) error
}

type JSONMarshaler struct{}

func (JSONMarshaler) Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONMarshaler) Unmarshal(bytes []byte, target interface{}) error {
	return json.Unmarshal(bytes, target)
}
