package counter

type Created struct {
	ID string
}

func (Created) EventName() string {
	return "Created_v1"
}

func (e Created) ApplyTo(c *Counter) error {
	c.id = e.ID
	return nil
}

type IncrementedBy struct {
	Value int
}

func (IncrementedBy) EventName() string {
	return "Created_v1"
}

func (e IncrementedBy) ApplyTo(c *Counter) error {
	c.currentValue += e.Value
	return nil
}
