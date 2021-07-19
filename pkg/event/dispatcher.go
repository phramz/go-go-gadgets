package event

import (
	"context"
	"sort"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/phramz/go-go-gadgets/pkg/logger"
	"github.com/satori/go.uuid"
	"github.com/thoas/go-funk"
)

var _ Dispatcher = (*defaultDispatcher)(nil)

// NewDispatcher returns a new instance of event.Dispatcher
func NewDispatcher(logger logger.Logger) Dispatcher {
	return &defaultDispatcher{
		logger:     logger,
		handlers:   make([]handler, 0),
		subscriber: make(map[string][]string),
	}
}

type handler struct {
	id        string
	eventName string
	listener  Listener
	priority  int
}

type defaultDispatcher struct {
	sync.Mutex
	logger     logger.Logger
	handlers   []handler
	subscriber map[string][]string
}

func (d *defaultDispatcher) Emit(eventName string, event Event) error {
	var result *multierror.Error

	for _, h := range d.handlers {
		if h.eventName != eventName {
			continue
		}

		if err := h.listener(eventName, event); err != nil {
			result = multierror.Append(result, err)

			d.logger.Errorf("error while dispatching event: %v", err)
		}

		if event.PropagationStopped() {
			break
		}
	}

	if result == nil {
		return nil
	}

	return result.ErrorOrNil()
}

func (d *defaultDispatcher) AddListener(eventName string, listener Listener) string {
	return d.AddListenerWithPriority(eventName, listener, 0)
}

func (d *defaultDispatcher) AddListenerWithPriority(eventName string, listener Listener, priority int) string {
	d.Lock()
	defer d.Unlock()

	listenerID := d.add(eventName, listener, priority)
	defer func() {
		if err := d.Emit(ListenerAdded, ListenerEvent(context.Background(), listenerID, eventName, priority, listener)); err != nil {
			d.logger.Fatalf("unable to add listener %q: %v", listenerID, err)
		}
	}()

	return listenerID
}

func (d *defaultDispatcher) add(eventName string, listener Listener, priority int) string {
	add := handler{id: uuid.NewV4().String(), eventName: eventName, listener: listener, priority: priority}

	d.handlers = append(d.handlers, add)

	sort.SliceStable(d.handlers, func(i, j int) bool {
		return d.handlers[i].priority < d.handlers[j].priority
	})

	return add.id
}

func (d *defaultDispatcher) RemoveListener(listenerID string) {
	d.Lock()
	defer d.Unlock()

	if d.remove(listenerID) {
		defer func() {
			if err := d.Emit(ListenerRemoved, ListenerEvent(context.Background(), listenerID, "", 0, nil)); err != nil {
				d.logger.Fatalf("unable to remove listener %q: %v", listenerID, err)
			}
		}()
	}
}

func (d *defaultDispatcher) remove(listenerID string) bool {
	cnt := len(d.handlers)
	d.handlers = funk.Filter(d.handlers, func(e handler) bool {
		return e.id != listenerID
	}).([]handler)

	return cnt > len(d.handlers)
}

func (d *defaultDispatcher) AddSubscriber(subscriber Subscriber) string {
	d.Lock()
	defer d.Unlock()

	subscriberID := uuid.NewV4().String()
	listenerIDs := make([]string, 0)

	for _, v := range subscriber.GetSubscriptions() {
		listenerIDs = append(listenerIDs, d.add(v.GetEventName(), v.GetListener(), v.GetPriority()))
	}

	d.subscriber[subscriberID] = listenerIDs

	if err := d.Emit(SubscriberAdded, SubscriberEvent(context.Background(), subscriberID, subscriber)); err != nil {
		d.logger.Fatalf("unable to add subscriber %q: %v", subscriberID, err)
		return ""
	}

	return subscriberID
}

func (d *defaultDispatcher) RemoveSubscriber(subscriberID string) {
	d.Lock()
	defer d.Unlock()

	if _, exists := d.subscriber[subscriberID]; !exists {
		return
	}

	if err := d.Emit(SubscriberRemoved, SubscriberEvent(context.Background(), subscriberID, nil)); err != nil {
		d.logger.Fatalf("unable to remove subscriber %q: %v", subscriberID, err)
		return
	}

	for _, listenerID := range d.subscriber[subscriberID] {
		d.remove(listenerID)
	}
}
