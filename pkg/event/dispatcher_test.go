package event

import (
	"context"
	"fmt"
	"testing"

	"github.com/phramz/go-go-gadgets/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var (
	_ Subscriber = &testSubscriberFoo{}
	_ Subscriber = &testSubscriberBar{}
)

type testSubscriber struct {
	result map[string]string
}

func (t *testSubscriber) handle(eventName, value string) {
	if t.result == nil {
		t.result = make(map[string]string)
	}

	if _, exists := t.result[eventName]; !exists {
		t.result[eventName] = ""
	}

	t.result[eventName] = fmt.Sprintf("%v%s", t.result[eventName], value)
}

type testSubscriberFoo struct {
	testSubscriber
}

type testSubscriberBar struct {
	testSubscriber
}

func (t *testSubscriberFoo) GetSubscriptions() []Subscription {
	return []Subscription{
		NewSubscriptionWithPriority("foo", func(event Event) {
			t.handle("foo", "!")
		}, 1),
		NewSubscription("foo", func(event Event) {
			t.handle("foo", "l")
		}),
		NewSubscriptionWithPriority("foo", func(event Event) {
			t.handle("foo", "o")
		}, -1),
		NewSubscription("foo", func(event Event) {
			t.handle("foo", "a")
		}),
		NewSubscriptionWithPriority("foo", func(event Event) {
			t.handle("foo", "h")
		}, -123),
	}
}

func (t *testSubscriberBar) GetSubscriptions() []Subscription {
	return []Subscription{
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "c")
		}, -100),
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "i")
		}, -1),
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "h")
		}, -100),
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "c")
		}, -1),
		NewSubscription("bar", func(event Event) {
			t.handle("bar", "h")
		}),
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "a")
		}, 55),
		NewSubscriptionWithPriority("bar", func(event Event) {
			t.handle("bar", "!")
		}, 99),
	}
}

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	assert.NotNil(t, dispatcher)
	assert.Implements(t, (*Dispatcher)(nil), dispatcher)
}

func TestDefaultDispatcher_AddListener(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	var result string
	dispatcher.AddListener("foo", func(event Event) {
		result = fmt.Sprintf("%v", event)
	})

	dispatcher.Dispatch("foo", "bar")
	assert.Equal(t, "bar", result)
}

func TestDefaultDispatcher_RemoveListener(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	var result string
	listenerID := dispatcher.AddListener("foo", func(event Event) {
		result = fmt.Sprintf("%v%v", result, event)
	})

	dispatcher.Dispatch("foo", "bar")
	dispatcher.Dispatch("foo", "bazz")
	assert.Equal(t, "barbazz", result)

	dispatcher.RemoveListener(listenerID)
	dispatcher.Dispatch("foo", "bar")
	dispatcher.Dispatch("foo", "bazz")

	assert.Equal(t, "barbazz", result)
}

func TestDefaultDispatcher_AddListenerWithPriority(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	var result string

	dispatcher.AddListenerWithPriority("foo", func(event Event) {
		result = fmt.Sprintf("%va", result)
	}, 100)

	dispatcher.AddListenerWithPriority("foo", func(event Event) {
		result = fmt.Sprintf("%vo", result)
	}, 0)

	dispatcher.AddListenerWithPriority("foo", func(event Event) {
		result = fmt.Sprintf("%v!", result)
	}, 255)

	dispatcher.AddListenerWithPriority("foo", func(event Event) {
		result = fmt.Sprintf("%vh", result)
	}, -100)

	dispatcher.AddListenerWithPriority("foo", func(event Event) {
		result = fmt.Sprintf("%vl", result)
	}, 0)

	dispatcher.Dispatch("foo", "bar")
	assert.Equal(t, "hola!", result)
}

func TestDefaultDispatcher_AddSubscriber(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	subscriberFoo := &testSubscriberFoo{}
	dispatcher.AddSubscriber(subscriberFoo)

	dispatcher.Dispatch("foo", "bar")

	assert.Equal(t, "hola!", subscriberFoo.result["foo"])

	dispatcher.Dispatch("foo", "bar")
	assert.Equal(t, "hola!hola!", subscriberFoo.result["foo"])

	subscriberBar := &testSubscriberBar{}
	dispatcher.AddSubscriber(subscriberBar)

	dispatcher.Dispatch("bar", "foo")
	assert.Equal(t, "chicha!", subscriberBar.result["bar"])
	dispatcher.Dispatch("bar", "foo")
	assert.Equal(t, "chicha!chicha!", subscriberBar.result["bar"])

	assert.Equal(t, "hola!hola!", subscriberFoo.result["foo"])
}

func TestDefaultDispatcher_RemoveSubscriber(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), logger.NewNullLogger())

	subscriberFoo := &testSubscriberFoo{}
	subscriberID := dispatcher.AddSubscriber(subscriberFoo)

	dispatcher.Dispatch("foo", "bar")

	assert.Equal(t, "hola!", subscriberFoo.result["foo"])

	dispatcher.RemoveSubscriber(subscriberID)

	dispatcher.Dispatch("foo", "bar")
	assert.Equal(t, "hola!", subscriberFoo.result["foo"])
}
