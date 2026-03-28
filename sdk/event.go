package sdk

// EventType identifies what changed in the store.
type EventType string

const (
	VoteSubmitted            EventType = "vote_submitted"
	HealthCheckCreated       EventType = "healthcheck_created"
	HealthCheckStatusChanged EventType = "healthcheck_status_changed"
	HealthCheckDeleted       EventType = "healthcheck_deleted"
)

// Event is a lightweight envelope carrying IDs of what changed.
type Event struct {
	Type          EventType `json:"type"`
	HealthCheckID string    `json:"healthcheck_id"`
	TeamID        string    `json:"team_id,omitempty"`
	Participant   string    `json:"participant,omitempty"`
	MetricName    string    `json:"metric_name,omitempty"`
}
