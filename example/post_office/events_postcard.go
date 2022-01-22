package post_office

import (
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Created struct {
	id aggregate.ID
}

func (Created) EventName() string {
	return "Created_v1"
}

func (l *Postcard) handleCreated(_ Created) {}

type Addressed struct {
	Sender    Address
	Addressee Address
}

func (Addressed) EventName() string {
	return "Addressed_v1"
}

func (l *Postcard) handleAddressed(e Addressed) {
	l.sender = e.Sender
	l.addressee = e.Addressee
}

type Written struct {
	content string
}

func (Written) EventName() string {
	return "Written_v1"
}

func (l *Postcard) handleWritten(e Written) {
	l.content = e.content
}

type Stamped struct {
	StampValue int
}

func (Stamped) EventName() string {
	return "Stamped_v1"
}

func (l *Postcard) handleStamped(e Stamped) {
	l.stamped = true
}

type Sent struct{}

func (Sent) EventName() string {
	return "Sent_v1"
}

func (l *Postcard) handleSent(e Sent) {
	l.sent = true
}

type Received struct{}

func (Received) EventName() string {
	return "Received_v1"
}

func (l *Postcard) handleReceived(e Received) {
	l.received = true
}
