package domain

import "time"

// ActionRepository defines persistence operations for actions.
type ActionRepository interface {
	Create(a *Action) error
	Complete(id string) error
	FindByHealthCheck(healthCheckID string) ([]*Action, error)
}

// Action represents a follow-up item from a health check discussion.
type Action struct {
	ID            string
	HealthCheckID string
	MetricName    string
	Description   string
	Assignee      string
	Completed     bool
	CreatedAt     time.Time
	CompletedAt   *time.Time
}
