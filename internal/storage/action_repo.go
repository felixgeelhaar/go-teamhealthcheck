package storage

import (
	"fmt"
	"time"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
)

func (s *Store) CreateAction(a *domain.Action) error {
	_, err := s.db.Exec(
		`INSERT INTO actions (id, healthcheck_id, metric_name, description, assignee, completed, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.HealthCheckID, a.MetricName, a.Description, a.Assignee, a.Completed, a.CreatedAt,
	)
	return err
}

func (s *Store) CompleteAction(id string) error {
	now := time.Now()
	res, err := s.db.Exec(
		`UPDATE actions SET completed = 1, completed_at = ? WHERE id = ?`,
		now, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("action %q not found", id)
	}
	return nil
}

func (s *Store) FindActionsByHealthCheck(healthCheckID string) ([]*domain.Action, error) {
	rows, err := s.db.Query(
		`SELECT id, healthcheck_id, metric_name, description, assignee, completed, created_at, completed_at
		 FROM actions WHERE healthcheck_id = ? ORDER BY created_at`,
		healthCheckID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*domain.Action
	for rows.Next() {
		a := &domain.Action{}
		var completedAt *time.Time
		if err := rows.Scan(&a.ID, &a.HealthCheckID, &a.MetricName, &a.Description, &a.Assignee, &a.Completed, &a.CreatedAt, &completedAt); err != nil {
			return nil, err
		}
		a.CompletedAt = completedAt
		actions = append(actions, a)
	}
	if actions == nil {
		actions = []*domain.Action{}
	}
	return actions, rows.Err()
}
