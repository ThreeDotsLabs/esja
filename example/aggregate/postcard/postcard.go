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
	p.eq = aggregate.NewEventsQueue(p)

	err := p.eq.PushAndApply(&Created{
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

func (p *Postcard) FromEvents(events []aggregate.VersionedEvent[*Postcard]) error {
	es, err := aggregate.NewEventsQueueFromEvents(p, events)
	if err != nil {
		return err
	}

	p.eq = es

	return nil
}

func (p *Postcard) ID() string {
	return p.id
}

func (p *Postcard) AggregateID() aggregate.ID {
	return aggregate.ID(p.id)
}

func (p *Postcard) Write(content string) error {
	return p.eq.PushAndApply(&Written{
		Content: content,
	})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return p.eq.PushAndApply(&Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
}

func (p *Postcard) Send() error {
	if p.sent {
		return fmt.Errorf("postcard already sent")
	}

	return p.eq.PushAndApply(&Sent{})
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
