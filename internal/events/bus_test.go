package events_test

import (
	"sync"
	"testing"
	"time"

	"github.com/felixgeelhaar/heartbeat/internal/events"
)

func TestBus_PublishToSubscribers(t *testing.T) {
	bus := events.NewBus()

	var mu sync.Mutex
	var received []events.Event

	bus.Subscribe(events.ListenerFunc(func(e events.Event) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e)
	}))

	bus.Publish(events.Event{
		Type:          events.VoteSubmitted,
		HealthCheckID: "hc-1",
		Participant:   "Alice",
		MetricName:    "Fun",
	})

	// Wait for async delivery
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received))
	}
	if received[0].Type != events.VoteSubmitted {
		t.Errorf("expected VoteSubmitted, got %s", received[0].Type)
	}
	if received[0].Participant != "Alice" {
		t.Errorf("expected Alice, got %s", received[0].Participant)
	}
}

func TestBus_MultipleListeners(t *testing.T) {
	bus := events.NewBus()

	var wg sync.WaitGroup
	count := 0
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		wg.Add(1)
		bus.Subscribe(events.ListenerFunc(func(e events.Event) {
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
		}))
	}

	bus.Publish(events.Event{Type: events.HealthCheckCreated, HealthCheckID: "hc-1"})
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if count != 3 {
		t.Errorf("expected 3 listener calls, got %d", count)
	}
}

func TestBus_NoListeners(t *testing.T) {
	bus := events.NewBus()
	// Should not panic
	bus.Publish(events.Event{Type: events.VoteSubmitted, HealthCheckID: "hc-1"})
}
