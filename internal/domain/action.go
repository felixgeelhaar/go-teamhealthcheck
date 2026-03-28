package domain

import "time"

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
