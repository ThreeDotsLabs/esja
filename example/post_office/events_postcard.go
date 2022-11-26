package post_office

import (
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Created struct {
	ID string
}

func (Created) EventName() event.Name {
	return "Created_v1"
}

func (Created) New() event.Event {
	return Created{}
}

func (p *Postcard) handleCreated(e Created) {
	p.id = e.ID
}

type Addressed struct {
	Sender    Address
	Addressee Address
}

func (Addressed) EventName() event.Name {
	return "Addressed_v1"
}

func (Addressed) New() event.Event {
	return Addressed{}
}

func (p *Postcard) handleAddressed(e Addressed) {
	p.sender = e.Sender
	p.addressee = e.Addressee
}

type Written struct {
	content string
}

func (Written) EventName() event.Name {
	return "Written_v1"
}

func (Written) New() event.Event {
	return Written{}
}

func (p *Postcard) handleWritten(e Written) {
	p.content = e.content
}

type Sent struct{}

func (Sent) EventName() event.Name {
	return "Sent_v1"
}

func (Sent) New() event.Event {
	return Sent{}
}

func (p *Postcard) handleSent(e Sent) {
	p.sent = true
}
