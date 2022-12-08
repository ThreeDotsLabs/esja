package postcard

type Created struct {
	ID string
}

func (Created) EventName() stream.EventName {
	return "Created_v1"
}

func (e Created) Apply(p *Postcard) error {
	p.id = e.ID
	return nil
}

type Addressed struct {
	Sender    Address
	Addressee Address
}

func (Addressed) EventName() stream.EventName {
	return "Addressed_v1"
}

func (e Addressed) Apply(p *Postcard) error {
	p.sender = e.Sender
	p.addressee = e.Addressee
	return nil
}

type Written struct {
	Content string
}

func (Written) EventName() stream.EventName {
	return "Written_v1"
}

func (e Written) Apply(p *Postcard) error {
	p.content = e.Content
	return nil
}

type Sent struct{}

func (Sent) EventName() stream.EventName {
	return "Sent_v1"
}

func (e Sent) Apply(p *Postcard) error {
	p.sent = true
	return nil
}
