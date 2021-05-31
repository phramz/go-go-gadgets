package event

var (
	_ Subscription = &defaultSubscription{}
)

// Subscriber abstraction for event subscribers
type Subscriber interface {
	GetSubscriptions() []Subscription
}

// NewSubscription return a new instance of Subscription
func NewSubscription(eventName string, listener Listener) Subscription {
	return &defaultSubscription{
		eventName: eventName,
		listener:  listener,
	}
}

// NewSubscriptionWithPriority return a new instance of Subscription with priority
func NewSubscriptionWithPriority(eventName string, listener Listener, priority int) Subscription {
	return &defaultSubscription{
		eventName: eventName,
		listener:  listener,
		priority:  priority,
	}
}

// Subscription abstraction for event subscriptions
type Subscription interface {
	GetEventName() string
	GetPriority() int
	GetListener() Listener
}

type defaultSubscription struct {
	eventName string
	listener  Listener
	priority  int
}

func (d defaultSubscription) GetEventName() string {
	return d.eventName
}

func (d defaultSubscription) GetPriority() int {
	return d.priority
}

func (d defaultSubscription) GetListener() Listener {
	return d.listener
}
