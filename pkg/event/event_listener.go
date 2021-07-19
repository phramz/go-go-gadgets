package event

import (
	"context"
)

// internal listener events
const (
	ListenerAdded   = "listener_added"
	ListenerRemoved = "listener_removed"
)

// ListenerEvent return a new event.Event
func ListenerEvent(ctx context.Context, id, eventName string, priority int, listener Listener) Event {
	return &listerEvent{
		defaultEvent: defaultEvent{
			ctx: ctx,
		},
		ID:        id,
		EventName: eventName,
		Priority:  priority,
		Listener:  listener,
	}
}

type listerEvent struct {
	defaultEvent
	ID        string
	EventName string
	Listener  Listener
	Priority  int
}
