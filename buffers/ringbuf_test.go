package buffers

import "testing"

func TestPushTailPopFront(t *testing.T) {
	const itemToPush = 0
	rb := NewRingBuf(10)

	if err := rb.PushTail(itemToPush); err != nil {
		t.Fatalf("push failed: %v", err)
	}
	it, err := rb.PopFront()
	if err != nil {
		t.Fatalf("pop failed: %v", err)
	}
	if v := it.(int); v != itemToPush {
		t.Fatalf("got %v value, want %v value", v, itemToPush)
	}
}

func TestPushWithItemRemoveCallback(t *testing.T) {
	itemsToPush := []int{0, 1, 2, 3}
	bufSize := 2

	rb := NewRingBuf(bufSize, WithItemRemoveCallback(func(item interface{}) {
		for _, it := range itemsToPush[:bufSize] {
			if item.(int) == it {
				t.Logf("remove: %v", item)
				return
			}
		}
		t.Fatalf("remove callback shouldn't be called for item %v", item)
	}))

	for _, it := range itemsToPush {
		if err := rb.PushTail(it); err != nil {
			t.Fatalf("push failed: %v", err)
		}
	}
}

func TestCloseWithItemRemoveCallback(t *testing.T) {
	itemsToPush := []int{0, 1, 2, 3}
	bufSize := len(itemsToPush)
	needRemove := false

	rb := NewRingBuf(bufSize, WithItemRemoveCallback(func(item interface{}) {
		if !needRemove {
			t.Fatal("rmove callback shouldn't be called now")
		}
		for _, it := range itemsToPush[:bufSize] {
			if item.(int) == it {
				t.Logf("remove: %v", item)
				return
			}
		}
		t.Fatalf("remove callback shouldn't be called for item %v", item)
	}))
	defer rb.Close()

	for _, it := range itemsToPush {
		if err := rb.PushTail(it); err != nil {
			t.Fatalf("push failed: %v", err)
		}
	}
	needRemove = true
}
