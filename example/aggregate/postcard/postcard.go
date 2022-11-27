package postcard

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Postcard struct {
	es aggregate.EventStore[*Postcard]

	id string

	sender    Address
	addressee Address

	content string

	sent bool
}

type Address struct {
	Name  string
	Line1 string
	Line2 string
	Line3 string
}

func NewPostcard(id string) (*Postcard, error) {
	p := &Postcard{}
	p.es = aggregate.NewEventStore(p)

	err := p.es.Record(Created{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Postcard) PopEvents() []aggregate.VersionedEvent[*Postcard] {
	return p.es.PopEvents()
}

func (p *Postcard) FromEvents(events []aggregate.VersionedEvent[*Postcard]) error {
	es, err := aggregate.NewEventStoreFromEvents(p, events)
	if err != nil {
		return err
	}

	p.es = es

	return nil
}

func (p *Postcard) ID() string {
	return p.id
}

func (p *Postcard) AggregateID() aggregate.ID {
	return aggregate.ID(p.id)
}

func (p *Postcard) Write(content string) error {
	return p.es.Record(Written{
		Content: content,
	})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return p.es.Record(Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
}

func (p *Postcard) Send() error {
	if p.sent {
		return fmt.Errorf("postcard already sent")
	}

	return p.es.Record(Sent{})
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
