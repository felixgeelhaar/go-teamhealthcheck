package sdk

import (
	"database/sql"
	"time"
)

// PluginContext provides access to core services during plugin initialization.
type PluginContext struct {
	// Store provides read-only access to core health check data.
	Store StoreReader
	// DB provides direct database access for plugin-owned tables.
	DB *sql.DB
	// Logger provides structured logging.
	Logger Logger
	// Bus provides event subscribe/publish.
	Bus EventBus
}

// Logger is a minimal logging interface compatible with bolt.
type Logger interface {
	Info() LogEvent
	Debug() LogEvent
	Error() LogEvent
	Warn() LogEvent
}

// LogEvent is a fluent log event builder.
type LogEvent interface {
	Str(key, val string) LogEvent
	Int(key string, val int) LogEvent
	Err(err error) LogEvent
	Msg(msg string)
}

// EventBus allows plugins to subscribe to and publish events.
type EventBus interface {
	Subscribe(listener EventListener)
	Publish(event Event)
}

// StoreReader provides read-only access to core health check data.
// Plugins should NOT write to core tables — use PluginContext.DB for plugin tables.
type StoreReader interface {
	FindTeamByID(id string) (*Team, error)
	FindAllTeams() ([]*Team, error)
	FindHealthCheckByID(id string) (*HealthCheck, error)
	FindAllHealthChecks(filter HealthCheckFilter) ([]*HealthCheck, error)
	FindTemplateByID(id string) (*Template, error)
	FindVotesByHealthCheck(id string) ([]*Vote, error)
}

// --- Core domain types (re-exported so plugins don't import internal/) ---

// Team represents a group of people who run health checks together.
type Team struct {
	ID        string
	Name      string
	Members   []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HealthCheck represents a single health check session.
type HealthCheck struct {
	ID         string
	TeamID     string
	TemplateID string
	Name       string
	Anonymous  bool
	Status     string
	CreatedAt  time.Time
	ClosedAt   *time.Time
}

// HealthCheckFilter specifies criteria for listing health checks.
type HealthCheckFilter struct {
	TeamID *string
	Status *string
	Limit  int
}

// Template is a reusable set of metrics for health checks.
type Template struct {
	ID          string
	Name        string
	Description string
	BuiltIn     bool
	Metrics     []TemplateMetric
	CreatedAt   time.Time
}

// TemplateMetric describes one dimension of a health check.
type TemplateMetric struct {
	ID              string
	TemplateID      string
	Name            string
	DescriptionGood string
	DescriptionBad  string
	SortOrder       int
}

// VoteColor represents the traffic-light vote choice.
type VoteColor string

const (
	VoteGreen  VoteColor = "green"
	VoteYellow VoteColor = "yellow"
	VoteRed    VoteColor = "red"
)

// Vote represents a single participant's assessment of one metric.
type Vote struct {
	ID            string
	HealthCheckID string
	MetricName    string
	Participant   string
	Color         VoteColor
	Comment       string
	CreatedAt     time.Time
}
