package rangedbtest

// ThingWasDone is an event used for testing.
type ThingWasDone struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

// AggregateID returns the aggregate id.
func (t ThingWasDone) AggregateID() string {
	return t.ID
}

// AggregateType returns the aggregate type.
func (t ThingWasDone) AggregateType() string {
	return "thing"
}

// EventType returns the event type. This will always be the struct name.
func (t ThingWasDone) EventType() string {
	return "ThingWasDone"
}

// AnotherWasComplete is an event used for testing.
type AnotherWasComplete struct {
	ID string `json:"id"`
}

// AggregateID returns the aggregate id.
func (t AnotherWasComplete) AggregateID() string {
	return t.ID
}

// AggregateType returns the aggregate type.
func (t AnotherWasComplete) AggregateType() string {
	return "another"
}

// EventType returns the event type. This will always be the struct name.
func (t AnotherWasComplete) EventType() string {
	return "AnotherWasComplete"
}
