package event

// Event generic event type alias
type Event = interface{}

// Listener generic listener type alias
type Listener = func(event Event)
