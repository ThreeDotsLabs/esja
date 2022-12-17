package postcard

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/stream"
)

type Postcard struct {
	events *stream.Events[Postcard]

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
	p := &Postcard{
		events: new(stream.Events[Postcard]),
	}

	err := stream.Record[Postcard](p, Created{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p Postcard) StreamID() stream.ID {
	return stream.ID(p.id)
}

func (p Postcard) Events() *stream.Events[Postcard] {
	return p.events
}

func (p Postcard) NewFromEvents(events *stream.Events[Postcard]) *Postcard {
	return &Postcard{events: events}
}

func (p *Postcard) ID() string {
	return p.id
}

func (p *Postcard) Write(content string) error {
	return stream.Record[Postcard](p, Written{
		Content: content,
	})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return stream.Record[Postcard](p, Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
}

func (p *Postcard) Send() error {
	if p.sent {
		return fmt.Errorf("postcard already sent")
	}

	return stream.Record[Postcard](p, Sent{})
}

func (p *Postcard) Sender() Address {
	return p.sender
}

func (p *Postcard) Addressee() Address {
	return p.addressee
}

func (p *Postcard) Content() string {
	return p.content
}

func (p *Postcard) Sent() bool {
	return p.sent
}
