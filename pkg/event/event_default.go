package event

import (
	"context"
	"fmt"
)

var (
	_ Event        = (*defaultEvent)(nil)
	_ fmt.Stringer = (*defaultEvent)(nil)
)

// DefaultEvent return a new Event
func DefaultEvent(ctx context.Context, payload interface{}) Event {
	return &defaultEvent{
		ctx:     ctx,
		payload: payload,
	}
}

type defaultEvent struct {
	ctx                context.Context
	payload            interface{}
	propagationStopped bool
}

func (d *defaultEvent) String() string {
	return fmt.Sprintf("%v", d.Payload())
}

func (d *defaultEvent) Payload() interface{} {
	return d.payload
}

func (d *defaultEvent) StopPropagation() {
	d.propagationStopped = true
}

func (d *defaultEvent) PropagationStopped() bool {
	return d.propagationStopped
}

func (d *defaultEvent) Context() context.Context {
	return d.ctx
}
