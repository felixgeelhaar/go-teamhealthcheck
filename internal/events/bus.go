package events

import "sync"

// EventType identifies what changed in the store.
type EventType string

const (
	VoteSubmitted            EventType = "vote_submitted"
	HealthCheckCreated       EventType = "healthcheck_created"
	HealthCheckStatusChanged EventType = "healthcheck_status_changed"
	HealthCheckDeleted       EventType = "healthcheck_deleted"
)

// Event is a lightweight envelope carrying only IDs.
// Dashboard clients use these IDs to refetch full state via REST endpoints.
type Event struct {
	Type          EventType `json:"type"`
	HealthCheckID string    `json:"healthcheck_id"`
	TeamID        string    `json:"team_id,omitempty"`
	Participant   string    `json:"participant,omitempty"`
	MetricName    string    `json:"metric_name,omitempty"`
}

// Listener receives events after successful store mutations.
type Listener interface {
	OnEvent(event Event)
}

// ListenerFunc adapts a plain function to the Listener interface.
type ListenerFunc func(Event)

func (f ListenerFunc) OnEvent(e Event) { f(e) }

// Bus is a fan-out event bus. Each listener is called in its own goroutine
// so it does not block the store mutation path.
type Bus struct {
	mu        sync.RWMutex
	listeners []Listener
}

// NewBus creates a new event bus.
func NewBus() *Bus { return &Bus{} }

// Subscribe registers a listener to receive all future events.
func (b *Bus) Subscribe(l Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, l)
}

// Publish sends an event to all registered listeners asynchronously.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, l := range b.listeners {
		go l.OnEvent(e)
	}
}
