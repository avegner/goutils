package buffers

import (
	"errors"
	"sync"
)

var ErrClosed = errors.New("closed")

type Option func(rb *RingBuf)

func WithItemRemoveCallback(remove ItemRemoveCallback) Option {
	return func(rb *RingBuf) {
		rb.remove = remove
	}
}

type RingBuf struct {
	items  []interface{}
	remove ItemRemoveCallback
	mu     sync.RWMutex
	donec  chan struct{}
}

type ItemRemoveCallback func(item interface{})

func NewRingBuf(size int, ops ...Option) *RingBuf {
	rb := &RingBuf{
		items: make([]interface{}, 0, size),
		donec: make(chan struct{}),
	}
	for _, op := range ops {
		op(rb)
	}
	return rb
}

func (rb *RingBuf) Close() error {
	select {
	case <-rb.donec:
		return ErrClosed
	default:
		close(rb.donec)
	}

	if rb.remove != nil {
		for _, item := range rb.items {
			rb.remove(item)
		}
	}
	return nil
}

func (rb *RingBuf) PushTail(item interface{}) error {
	select {
	case <-rb.donec:
		return ErrClosed
	default:
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.items) >= cap(rb.items) {
		if rb.remove != nil {
			rb.remove(rb.items[0])
		}
		rb.removeFirst()
	}
	rb.items = append(rb.items, item)
	return nil
}

func (rb *RingBuf) PopFront() (interface{}, error) {
	select {
	case <-rb.donec:
		return nil, ErrClosed
	default:
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	var item interface{}
	if len(rb.items) > 0 {
		item = rb.items[0]
		rb.removeFirst()
	}
	return item, nil
}

func (rb *RingBuf) PopTail() (interface{}, error) {
	select {
	case <-rb.donec:
		return nil, ErrClosed
	default:
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	var item interface{}
	if len(rb.items) > 0 {
		item = rb.items[len(rb.items)-1]
		rb.items = rb.items[:len(rb.items)-1]
	}
	return item, nil
}

func (rb *RingBuf) PopAll() ([]interface{}, error) {
	select {
	case <-rb.donec:
		return nil, ErrClosed
	default:
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	items := rb.items
	rb.items = rb.items[:0]
	return items, nil
}

func (rb *RingBuf) ItemsBuffered() (int, error) {
	select {
	case <-rb.donec:
		return 0, ErrClosed
	default:
	}

	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return len(rb.items), nil
}

func (rb *RingBuf) removeFirst() {
	copy(rb.items, rb.items[1:])
	rb.items = rb.items[:len(rb.items)-1]
}
