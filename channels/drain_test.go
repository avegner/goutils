package channels_test

import (
	"testing"

	"github.com/avegner/utils/channels"
)

func TestDrainEmptyChannel(t *testing.T) {
	testDrain(t, 1, 0, false, false)
}

func TestDrainEmptyAndClosedChannel(t *testing.T) {
	testDrain(t, 1, 0, true, false)
}

func TestDrainFilledChannel(t *testing.T) {
	testDrain(t, 3, 3, false, false)
}

func TestDrainFilledAndClosedChannel(t *testing.T) {
	testDrain(t, 3, 3, true, false)
}

func TestDrainNilChannel(t *testing.T) {
	defer checkPanic(t, false)
	var c chan struct{}
	channels.Drain(c)
}

func TestDrainNil(t *testing.T) {
	defer checkPanic(t, true)
	channels.Drain(nil)
}

func TestDrainInt(t *testing.T) {
	defer checkPanic(t, true)
	channels.Drain(0)
}

func TestDrainSendOnlyChannel(t *testing.T) {
	defer checkPanic(t, true)
	c := initChan(1, 0)

	func(c chan<- struct{}) {
		channels.Drain(c)
	}(c)
}

func TestDrainReceiveOnlyChannel(t *testing.T) {
	defer checkPanic(t, false)
	c := initChan(1, 1)

	func(c <-chan struct{}) {
		channels.Drain(c)
	}(c)
	checkDrained(t, c)
}

func initChan(size, elems int) chan struct{} {
	c := make(chan struct{}, size)
	for i := 0; i < elems; i++ {
		c <- struct{}{}
	}
	return c
}

func checkPanic(t *testing.T, need bool) {
	if r := recover(); r != nil {
		t.Logf("panic: %v", r)
		if !need {
			t.Fatalf("got panic, want no panic")
		}
		return
	}
	if need {
		t.Fatalf("got no panic, want panic")
	}
}

func checkDrained(t *testing.T, c chan struct{}) {
	select {
	case _, ok := <-c:
		if ok {
			t.Fatalf("channel not drained")
		}
	default:
	}
}

func testDrain(t *testing.T, size, elems int, needClose, needPanic bool) {
	defer checkPanic(t, needPanic)
	c := initChan(size, elems)
	if needClose {
		close(c)
	}
	channels.Drain(c)
	checkDrained(t, c)
}

func BenchmarkDrainReflect(b *testing.B) {
	c := make(chan struct{}, 1)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		channels.Drain(c)
	}
}

func BenchmarkDrainNative(b *testing.B) {
	c := make(chan struct{}, 1)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		select {
		case <-c:
		default:
		}
	}
}
