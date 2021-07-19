package event

import "context"

// Event generic event type alias
type Event interface {
	Context() context.Context
	Payload() interface{}
	StopPropagation()
	PropagationStopped() bool
}

// Listener generic listener type alias
type Listener = func(name string, event Event) error

// ErrorFn handle errors
type ErrorFn = func(err error) (stopPropagation bool)

// Dispatcher abstraction for event dispatchers
type Dispatcher interface {
	Emitter
	Registry
}

// Registry abstraction for event registries
type Registry interface {
	AddListener(eventName string, listener Listener) string
	AddListenerWithPriority(eventName string, listener Listener, priority int) string
	AddSubscriber(subscriber Subscriber) string
	RemoveListener(listenerID string)
	RemoveSubscriber(subscriberID string)
}

// Emitter abstraction for event emitters
type Emitter interface {
	Emit(eventName string, event Event) error
}

// Subscription abstraction for event subscriptions
type Subscription interface {
	GetEventName() string
	GetPriority() int
	GetListener() Listener
}

// Subscriber abstraction for event subscribers
type Subscriber interface {
	GetSubscriptions() []Subscription
}
