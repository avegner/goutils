package subs

type Subscriber interface {
	Subscribe(event Event, cb EventCallback) Subscription
	Unsubscribe(sub Subscription)
}

type Subscription interface{}

type Event interface{}

type EventCallback func()
