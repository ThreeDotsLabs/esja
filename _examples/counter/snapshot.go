package counter

type Snapshot struct {
	ID           string
	CurrentValue int
}

func (s Snapshot) EventName() string {
	return "CounterSnapshot_v1"
}

func (s Snapshot) ApplyTo(c *Counter) error {
	c.id = s.ID
	c.currentValue = s.CurrentValue
	return nil
}
