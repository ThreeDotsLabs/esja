package postcard

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

type Postcard struct {
	stream *stream.Stream[Postcard]

	id string

	sender    Address
	addressee Address
	content   string
	sent      bool
}
type Address struct {
	Name  string `anonymize:"true"`
	Line1 string
	Line2 string
	Line3 string
}

func NewPostcard(id string) (*Postcard, error) {
	s, err := stream.NewStreamWithType[Postcard](stream.ID(id), "Postcard")
	if err != nil {
		return nil, err
	}

	p := &Postcard{
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

func (p *Postcard) Send() error {
	if p.sent {
		return fmt.Errorf("postcard already sent")
	}

	return p.stream.Record(p, Sent{})
}

func (p Postcard) Stream() *stream.Stream[Postcard] {
	return p.stream
}

func (p Postcard) NewWithStream(stream *stream.Stream[Postcard]) *Postcard {
	return &Postcard{stream: stream}
}

func (p Postcard) ID() string {
	return p.id
}

func (p Postcard) Sender() Address {
	return p.sender
}

func (p Postcard) Addressee() Address {
	return p.addressee
}

func (p Postcard) Content() string {
	return p.content
}

func (p Postcard) Sent() bool {
	return p.sent
}

func (p *Postcard) Write(content string) error {
	return p.stream.Record(p, Written{
		Content: content,
	})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return p.stream.Record(p, Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
}
