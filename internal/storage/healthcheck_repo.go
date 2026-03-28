package storage

import (
	"database/sql"
	"fmt"

	"github.com/felixgeelhaar/heartbeat/internal/domain"
	"github.com/felixgeelhaar/heartbeat/internal/events"
)

func (s *Store) CreateHealthCheck(hc *domain.HealthCheck) error {
	_, err := s.db.Exec(
		`INSERT INTO healthchecks (id, team_id, template_id, name, anonymous, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		hc.ID, hc.TeamID, hc.TemplateID, hc.Name, hc.Anonymous, hc.Status, hc.CreatedAt,
	)
	if err == nil {
		s.publish(events.Event{
			Type:          events.HealthCheckCreated,
			HealthCheckID: hc.ID,
			TeamID:        hc.TeamID,
		})
	}
	return err
}

func (s *Store) FindHealthCheckByID(id string) (*domain.HealthCheck, error) {
	hc := &domain.HealthCheck{}
	var closedAt sql.NullTime
	err := s.db.QueryRow(
		`SELECT id, team_id, template_id, name, anonymous, status, created_at, closed_at
		 FROM healthchecks WHERE id = ?`, id,
	).Scan(&hc.ID, &hc.TeamID, &hc.TemplateID, &hc.Name, &hc.Anonymous, &hc.Status, &hc.CreatedAt, &closedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if closedAt.Valid {
		hc.ClosedAt = &closedAt.Time
	}
	return hc, nil
}

func (s *Store) FindAllHealthChecks(filter domain.HealthCheckFilter) ([]*domain.HealthCheck, error) {
	query := `SELECT id, team_id, template_id, name, anonymous, status, created_at, closed_at FROM healthchecks WHERE 1=1`
	var args []any

	if filter.TeamID != nil {
		query += ` AND team_id = ?`
		args = append(args, *filter.TeamID)
	}
	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, string(*filter.Status))
	}

	query += ` ORDER BY created_at DESC`

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	query += ` LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.HealthCheck
	for rows.Next() {
		hc := &domain.HealthCheck{}
		var closedAt sql.NullTime
		if err := rows.Scan(&hc.ID, &hc.TeamID, &hc.TemplateID, &hc.Name, &hc.Anonymous, &hc.Status, &hc.CreatedAt, &closedAt); err != nil {
			return nil, err
		}
		if closedAt.Valid {
			hc.ClosedAt = &closedAt.Time
		}
		results = append(results, hc)
	}
	return results, rows.Err()
}

func (s *Store) UpdateHealthCheck(hc *domain.HealthCheck) error {
	_, err := s.db.Exec(
		`UPDATE healthchecks SET status = ?, closed_at = ? WHERE id = ?`,
		hc.Status, hc.ClosedAt, hc.ID,
	)
	if err == nil {
		s.publish(events.Event{
			Type:          events.HealthCheckStatusChanged,
			HealthCheckID: hc.ID,
		})
	}
	return err
}

func (s *Store) DeleteHealthCheck(id string) error {
	res, err := s.db.Exec(`DELETE FROM healthchecks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("health check %q not found", id)
	}
	s.publish(events.Event{
		Type:          events.HealthCheckDeleted,
		HealthCheckID: id,
	})
	return nil
}
