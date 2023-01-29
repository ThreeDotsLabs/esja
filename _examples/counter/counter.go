package counter

import (
	"github.com/ThreeDotsLabs/esja"
)

type Counter struct {
	stream *esja.Stream[Counter]

	id           string
	currentValue int
}

func NewCounter(id string) (*Counter, error) {
	s, err := esja.NewStreamWithType[Counter](id, "Counter")
	if err != nil {
		return nil, err
	}

	p := &Counter{
		stream: s,
	}

	err = p.stream.Record(p, Created{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c Counter) Stream() *esja.Stream[Counter] {
	return c.stream
}

func (c Counter) NewWithStream(s *esja.Stream[Counter]) *Counter {
	return &Counter{stream: s}
}

func (c Counter) Snapshot() esja.Snapshot[Counter] {
	return Snapshot{
		ID:           c.id,
		CurrentValue: c.currentValue,
	}
}

func (c Counter) ID() string {
	return c.id
}

func (c *Counter) CurrentValue() int {
	return c.currentValue
}

func (c *Counter) IncrementBy(v int) error {
	return c.stream.Record(c, IncrementedBy{
		Value: v,
	})
}
