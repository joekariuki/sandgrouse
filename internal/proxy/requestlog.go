package proxy

import (
	"sync"
	"time"
)

// RequestEvent captures per-request bandwidth data.
type RequestEvent struct {
	Timestamp    time.Time `json:"timestamp"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Provider     string    `json:"provider"`
	RequestBytes int64     `json:"request_bytes"`
	ResponseWire int64     `json:"response_wire"`
	ResponseOrig int64     `json:"response_orig"`
}

// RequestLog maintains a ring buffer of recent request events
// and supports fan-out subscriptions for live streaming.
type RequestLog struct {
	mu          sync.Mutex
	events      []RequestEvent
	capacity    int
	subscribers []chan RequestEvent
}

// NewRequestLog creates a RequestLog with the given capacity.
func NewRequestLog(capacity int) *RequestLog {
	return &RequestLog{
		events:   make([]RequestEvent, 0, capacity),
		capacity: capacity,
	}
}

// Add appends an event to the ring buffer and notifies subscribers.
func (rl *RequestLog) Add(event RequestEvent) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if len(rl.events) >= rl.capacity {
		// Shift: drop oldest, append newest
		copy(rl.events, rl.events[1:])
		rl.events[len(rl.events)-1] = event
	} else {
		rl.events = append(rl.events, event)
	}

	// Non-blocking send to all subscribers
	for _, ch := range rl.subscribers {
		select {
		case ch <- event:
		default:
			// Subscriber is slow, drop the event for them
		}
	}
}

// Recent returns a copy of the buffered events, newest first.
func (rl *RequestLog) Recent() []RequestEvent {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	result := make([]RequestEvent, len(rl.events))
	for i, e := range rl.events {
		result[len(rl.events)-1-i] = e
	}
	return result
}

// Subscribe returns a channel that receives new request events
// and an unsubscribe function to stop receiving.
func (rl *RequestLog) Subscribe() (<-chan RequestEvent, func()) {
	ch := make(chan RequestEvent, 16)
	rl.mu.Lock()
	rl.subscribers = append(rl.subscribers, ch)
	rl.mu.Unlock()

	unsubscribe := func() {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		for i, sub := range rl.subscribers {
			if sub == ch {
				rl.subscribers = append(rl.subscribers[:i], rl.subscribers[i+1:]...)
				close(ch)
				break
			}
		}
	}
	return ch, unsubscribe
}
