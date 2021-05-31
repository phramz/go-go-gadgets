package event

// Event generic event type alias
type Event interface {
}

// Listener generic listener type alias
type Listener = func(name string, event Event) error

// ErrorFn handle errors
type ErrorFn = func(err error) (stopPropagation bool)
