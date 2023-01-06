package postcard

type Created struct {
	ID string
}

func (Created) EventName() string {
	return "Created_v1"
}

func (e Created) ApplyTo(p *Postcard) error {
	p.id = e.ID
	return nil
}

type Addressed struct {
	Sender    Address
	Addressee Address
}

func (Addressed) EventName() string {
	return "Addressed_v1"
}

func (e Addressed) ApplyTo(p *Postcard) error {
	p.sender = e.Sender
	p.addressee = e.Addressee
	return nil
}

type Written struct {
	Content string
}

func (Written) EventName() string {
	return "Written_v1"
}

func (e Written) ApplyTo(p *Postcard) error {
	p.content = e.Content
	return nil
}

type Sent struct{}

func (Sent) EventName() string {
	return "Sent_v1"
}

func (e Sent) ApplyTo(p *Postcard) error {
	p.sent = true
	return nil
}
