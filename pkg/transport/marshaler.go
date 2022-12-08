package transport

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/esja/stream"
)

type Marshaler interface {
	Marshal(streamID stream.ID, data interface{}) ([]byte, error)
	Unmarshal(streamID stream.ID, bytes []byte, target interface{}) error
}

type JSONMarshaler struct{}

func (JSONMarshaler) Marshal(_ stream.ID, data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONMarshaler) Unmarshal(_ stream.ID, bytes []byte, target interface{}) error {
	return json.Unmarshal(bytes, target)
}
