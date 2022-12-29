package transport

import (
	"bytes"
	"encoding/gob"
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

type GOBMarshaler struct{}

func (GOBMarshaler) Marshal(data interface{}) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	e := gob.NewEncoder(b)
	err := e.Encode(data)
	return b.Bytes(), err
}

func (GOBMarshaler) Unmarshal(data []byte, target interface{}) error {
	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	err := d.Decode(target)
	return err
}
