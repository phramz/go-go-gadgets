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

var _ Dispatcher = &defaultDispatcher{}

// NewDispatcher returns a new instance of event.Dispatcher
func NewDispatcher(ctx context.Context, logger logger.Logger) Dispatcher {
	return &defaultDispatcher{
		ctx:        ctx,
		logger:     logger,
		handlers:   make([]handler, 0),
		subscriber: make(map[string][]string),
	}
}

// Dispatcher abstraction for event dispatchers
type Dispatcher interface {
	Fire(eventName string, event Event, onError ...func(err error) (stopPropagation bool))
	Dispatch(eventName string, event Event, onError ...func(err error) (stopPropagation bool)) error
	AddListener(eventName string, listener Listener) string
	AddListenerWithPriority(eventName string, listener Listener, priority int) string
	AddSubscriber(subscriber Subscriber) string
	RemoveListener(listenerID string)
	RemoveSubscriber(subscriberID string)
}

type handler struct {
	id        string
	eventName string
	listener  Listener
	priority  int
}

type defaultDispatcher struct {
	sync.Mutex
	ctx        context.Context
	logger     logger.Logger
	handlers   []handler
	subscriber map[string][]string
}

func (d *defaultDispatcher) Fire(eventName string, event Event, onError ...func(err error) (stopPropagation bool)) {
	go func(n string, e Event) {
		_ = d.Dispatch(n, e, onError...)
	}(eventName, event)
}

func (d *defaultDispatcher) Dispatch(eventName string, event Event, onError ...func(err error) (stopPropagation bool)) error {
	var result *multierror.Error

	for _, h := range d.handlers {
		if h.eventName != eventName {
			continue
		}

		if err := h.listener(event); err != nil {
			result = multierror.Append(result, err)

			d.logger.Errorf("error while dispatching event: %v", err)

			if len(onError) < 1 {
				continue
			}

			for _, fn := range onError {
				if fn(err) {
					return result.ErrorOrNil()
				}
			}
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

	return d.add(eventName, listener, priority)
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

	d.remove(listenerID)
}

func (d *defaultDispatcher) remove(listenerID string) {
	d.handlers = funk.Filter(d.handlers, func(e handler) bool {
		return e.id != listenerID
	}).([]handler)
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

	return subscriberID
}

func (d *defaultDispatcher) RemoveSubscriber(subscriberID string) {
	d.Lock()
	defer d.Unlock()

	if _, exists := d.subscriber[subscriberID]; !exists {
		return
	}

	for _, listenerID := range d.subscriber[subscriberID] {
		d.remove(listenerID)
	}
}
