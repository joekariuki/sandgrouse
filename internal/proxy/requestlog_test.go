package proxy

import (
	"sync"
	"testing"
	"time"
)

func TestRequestLogAdd(t *testing.T) {
	rl := NewRequestLog(50)
	rl.Add(RequestEvent{Method: "POST", Path: "/v1/messages", Provider: "anthropic"})
	rl.Add(RequestEvent{Method: "POST", Path: "/v1/chat", Provider: "openai"})

	recent := rl.Recent()
	if len(recent) != 2 {
		t.Fatalf("Recent() len = %d, want 2", len(recent))
	}
	// Newest first
	if recent[0].Provider != "openai" {
		t.Errorf("Recent()[0].Provider = %q, want 'openai'", recent[0].Provider)
	}
	if recent[1].Provider != "anthropic" {
		t.Errorf("Recent()[1].Provider = %q, want 'anthropic'", recent[1].Provider)
	}
}

func TestRequestLogRingBuffer(t *testing.T) {
	rl := NewRequestLog(5)
	for i := 0; i < 10; i++ {
		rl.Add(RequestEvent{RequestBytes: int64(i)})
	}

	recent := rl.Recent()
	if len(recent) != 5 {
		t.Fatalf("Recent() len = %d, want 5", len(recent))
	}
	// Should have events 5-9, newest first
	if recent[0].RequestBytes != 9 {
		t.Errorf("Recent()[0].RequestBytes = %d, want 9", recent[0].RequestBytes)
	}
	if recent[4].RequestBytes != 5 {
		t.Errorf("Recent()[4].RequestBytes = %d, want 5", recent[4].RequestBytes)
	}
}

func TestRequestLogSubscribe(t *testing.T) {
	rl := NewRequestLog(50)
	ch, unsub := rl.Subscribe()
	defer unsub()

	event := RequestEvent{Method: "POST", Provider: "anthropic"}
	rl.Add(event)

	select {
	case got := <-ch:
		if got.Provider != "anthropic" {
			t.Errorf("received event Provider = %q, want 'anthropic'", got.Provider)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on subscriber channel")
	}
}

func TestRequestLogUnsubscribe(t *testing.T) {
	rl := NewRequestLog(50)
	ch, unsub := rl.Subscribe()
	unsub()

	rl.Add(RequestEvent{Method: "POST"})

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("should not receive events after unsubscribe")
		}
		// Channel closed, expected
	case <-time.After(100 * time.Millisecond):
		// No event received, also acceptable (channel closed)
	}
}

func TestRequestLogConcurrency(t *testing.T) {
	rl := NewRequestLog(50)
	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				rl.Add(RequestEvent{RequestBytes: int64(n*20 + j)})
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				rl.Recent()
			}
		}()
	}

	wg.Wait()

	recent := rl.Recent()
	if len(recent) != 50 {
		t.Errorf("Recent() len = %d, want 50 after 200 inserts", len(recent))
	}
}
