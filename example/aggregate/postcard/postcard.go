package postcard

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Postcard struct {
	eq aggregate.EventsQueue[*Postcard]

	id string

	sender    Address
	addressee Address

	content string

	sent bool
}

func NewPostcard(id string) (*Postcard, error) {
	p := &Postcard{}

	err := aggregate.Record[*Postcard](p, &p.eq, &Created{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Postcard) PopEvents() []aggregate.VersionedEvent[*Postcard] {
	return p.eq.PopEvents()
}

func (p *Postcard) FromEventsQueue(eq aggregate.EventsQueue[*Postcard]) error {
	events := eq.PopEvents()

	for _, e := range events {
		err := e.Apply(p)
		if err != nil {
			return err
		}
	}

	p.eq = eq

	return nil
}

func (p *Postcard) ID() string {
	return p.id
}

func (p *Postcard) AggregateID() aggregate.ID {
	return aggregate.ID(p.id)
}

func (p *Postcard) Write(content string) error {
	return aggregate.Record[*Postcard](p, &p.eq, &Written{
		Content: content,
	})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return aggregate.Record[*Postcard](p, &p.eq, &Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
}

func (p *Postcard) Send() error {
	if p.sent {
		return fmt.Errorf("postcard already sent")
	}

	return aggregate.Record[*Postcard](p, &p.eq, &Sent{})
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
