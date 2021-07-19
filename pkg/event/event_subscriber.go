package event

import (
	"context"
)

var _ Event = (*subscriberEvent)(nil)

// internal subscription events
const (
	SubscriberAdded   = "subscriber_added"
	SubscriberRemoved = "subscriber_removed"
)

// SubscriberEvent return a new event.Event
func SubscriberEvent(ctx context.Context, id string, subscriber Subscriber) Event {
	return &subscriberEvent{
		defaultEvent: defaultEvent{
			ctx: ctx,
		},
		ID:         id,
		Subscriber: subscriber,
	}
}

type subscriberEvent struct {
	defaultEvent
	ID         string
	Subscriber Subscriber
}
