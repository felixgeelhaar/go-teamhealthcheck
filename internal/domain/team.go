package domain

import "time"

// Team is an aggregate root representing a group of people who run health checks together.
type Team struct {
	ID        string
	Name      string
	Members   []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TeamRepository defines persistence operations for the Team aggregate.
type TeamRepository interface {
	Create(team *Team) error
	FindByID(id string) (*Team, error)
	FindAll() ([]*Team, error)
	Delete(id string) error
	AddMember(teamID, name string) error
	RemoveMember(teamID, name string) error
}
