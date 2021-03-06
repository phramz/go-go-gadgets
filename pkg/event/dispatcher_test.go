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
		NewSubscriptionWithPriority("foo", func(_ string, event Event) error {
			t.handle("foo", "!")
			return nil
		}, 1),
		NewSubscription("foo", func(_ string, event Event) error {
			t.handle("foo", "l")
			return nil
		}),
		NewSubscriptionWithPriority("foo", func(_ string, event Event) error {
			t.handle("foo", "o")
			return nil
		}, -1),
		NewSubscription("foo", func(_ string, event Event) error {
			t.handle("foo", "a")
			return nil
		}),
		NewSubscriptionWithPriority("foo", func(_ string, event Event) error {
			t.handle("foo", "h")
			return nil
		}, -123),
	}
}

func (t *testSubscriberBar) GetSubscriptions() []Subscription {
	return []Subscription{
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "c")
			return nil
		}, -100),
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "i")
			return nil
		}, -1),
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "h")
			return nil
		}, -100),
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "c")
			return nil
		}, -1),
		NewSubscription("bar", func(_ string, event Event) error {
			t.handle("bar", "h")
			return nil
		}),
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "a")
			return nil
		}, 55),
		NewSubscriptionWithPriority("bar", func(_ string, event Event) error {
			t.handle("bar", "!")
			return nil
		}, 99),
	}
}

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	assert.NotNil(t, dispatcher)
	assert.Implements(t, (*Dispatcher)(nil), dispatcher)
}

func TestDefaultDispatcher_AddListener(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	var result string
	dispatcher.AddListener("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%v", event)
		return nil
	})

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.Equal(t, "bar", result)
}

func TestDefaultDispatcher_RemoveListener(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	var result string
	listenerID := dispatcher.AddListener("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%v%v", result, event)
		return nil
	})

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bazz")))
	assert.Equal(t, "barbazz", result)

	dispatcher.RemoveListener(listenerID)
	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bazz")))

	assert.Equal(t, "barbazz", result)
}

func TestDefaultDispatcher_AddListenerWithPriority(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	var result string

	dispatcher.AddListenerWithPriority("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%va", result)
		return nil
	}, 100)

	dispatcher.AddListenerWithPriority("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%vo", result)
		return nil
	}, 0)

	dispatcher.AddListenerWithPriority("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%v!", result)
		return nil
	}, 255)

	dispatcher.AddListenerWithPriority("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%vh", result)
		return nil
	}, -100)

	dispatcher.AddListenerWithPriority("foo", func(_ string, event Event) error {
		result = fmt.Sprintf("%vl", result)
		return nil
	}, 0)

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.Equal(t, "hola!", result)
}

func TestDefaultDispatcher_AddSubscriber(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	subscriberFoo := &testSubscriberFoo{}
	dispatcher.AddSubscriber(subscriberFoo)

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))

	assert.Equal(t, "hola!", subscriberFoo.result["foo"])

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.Equal(t, "hola!hola!", subscriberFoo.result["foo"])

	subscriberBar := &testSubscriberBar{}
	dispatcher.AddSubscriber(subscriberBar)

	assert.NoError(t, dispatcher.Emit("bar", DefaultEvent(context.TODO(), "foo")))
	assert.Equal(t, "chicha!", subscriberBar.result["bar"])
	assert.NoError(t, dispatcher.Emit("bar", DefaultEvent(context.TODO(), "foo")))
	assert.Equal(t, "chicha!chicha!", subscriberBar.result["bar"])

	assert.Equal(t, "hola!hola!", subscriberFoo.result["foo"])
}

func TestDefaultDispatcher_RemoveSubscriber(t *testing.T) {
	dispatcher := NewDispatcher(logger.NewNullLogger())

	subscriberFoo := &testSubscriberFoo{}
	subscriberID := dispatcher.AddSubscriber(subscriberFoo)

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))

	assert.Equal(t, "hola!", subscriberFoo.result["foo"])

	dispatcher.RemoveSubscriber(subscriberID)

	assert.NoError(t, dispatcher.Emit("foo", DefaultEvent(context.TODO(), "bar")))
	assert.Equal(t, "hola!", subscriberFoo.result["foo"])
}
