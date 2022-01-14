package post_office

import (
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Postcard struct {
	id aggregate.ID

	sender    Address
	addressee Address

	content string

	stamped  bool
	sent     bool
	received bool
}

func (l *Postcard) ID() aggregate.ID {
	return l.id
}

func (l *Postcard) Sender() Address {
	return l.sender
}

func (l *Postcard) Addressee() Address {
	return l.addressee
}

func (l *Postcard) Content() string {
	return l.content
}

func (l *Postcard) Stamped() bool {
	return l.stamped
}

func (l *Postcard) Sent() bool {
	return l.sent
}

func (l *Postcard) Received() bool {
	return l.received
}

func (l *Postcard) Handle(event event.Event) error {
	switch e := event.(type) {
	case Created:
		l.handleCreated(e)
	case Addressed:
		l.handleAddressed(e)
	case Written:
		l.handleWritten(e)
	case Stamped:
		l.handleStamped(e)
	case Sent:
		l.handleSent(e)
	case Received:
		l.handleReceived(e)
	}

	return fmt.Errorf("don't know how to handle event '%s'", event.EventName())
}

func NewPostcardAggregate(id aggregate.ID) (aggregate.Aggregate[Postcard], error) {
	return aggregate.NewAggregate[Postcard](id, &Postcard{})
}
