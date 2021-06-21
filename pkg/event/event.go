package event

// Event generic event type alias
type Event interface {
}

// Listener generic listener type alias
type Listener = func(name string, event Event) error

// ErrorFn handle errors
type ErrorFn = func(err error) (stopPropagation bool)

// Dispatcher abstraction for event dispatchers
type Dispatcher interface {
	Emitter
	AddListener(eventName string, listener Listener) string
	AddListenerWithPriority(eventName string, listener Listener, priority int) string
	AddSubscriber(subscriber Subscriber) string
	RemoveListener(listenerID string)
	RemoveSubscriber(subscriberID string)
}

// Emitter abstraction for event emitters
type Emitter interface {
	Fire(eventName string, event Event, onError ...ErrorFn)
	Dispatch(eventName string, event Event, onError ...ErrorFn) error
}
