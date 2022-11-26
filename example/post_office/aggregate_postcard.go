package post_office

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Postcard struct {
	aggregate.Aggregate

	id string

	sender    Address
	addressee Address

	content string

	sent bool
}

func NewPostcard(id string) (*Postcard, error) {
	p := &Postcard{}
	err := p.handle(Created{ID: id})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Postcard) ID() string {
	return p.id
}

func (p *Postcard) AggregateID() aggregate.ID {
	return aggregate.ID(p.id)
}

func (p *Postcard) Write() error {
	return p.handle(Written{})
}

func (p *Postcard) Address(sender Address, addressee Address) error {
	return p.handle(Addressed{
		Sender:    sender,
		Addressee: addressee,
	})
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

func NewPostcardFromEvents(events []event.Event) (*Postcard, error) {
	p := &Postcard{}
	for _, ev := range events {
		err := p.handle(ev)
		if err != nil {
			return nil, fmt.Errorf("error applying event '%s': %w", ev.EventName(), err)
		}
	}
	p.PopEvents()
	return p, nil
}

func (p *Postcard) handle(event event.Event) error {
	switch e := event.(type) {
	case Created:
		p.handleCreated(e)
	case Addressed:
		p.handleAddressed(e)
	case Written:
		p.handleWritten(e)
	case Sent:
		p.handleSent(e)
	default:
		return fmt.Errorf("don't know how to handle event '%s'", event.EventName())
	}

	p.RecordEvent(event)

	return nil
}
