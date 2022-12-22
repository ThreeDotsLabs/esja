package transport

import (
	"bytes"
	"encoding/gob"
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

type GOBMarshaler struct{}

func (GOBMarshaler) Marshal(_ stream.ID, data interface{}) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	e := gob.NewEncoder(b)
	err := e.Encode(data)
	return b.Bytes(), err
}

func (GOBMarshaler) Unmarshal(_ stream.ID, data []byte, target interface{}) error {
	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	err := d.Decode(target)
	return err
}
