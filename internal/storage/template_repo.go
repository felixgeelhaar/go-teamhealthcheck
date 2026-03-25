package storage

import (
	"database/sql"
	"fmt"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
)

func (s *Store) CreateTemplate(tmpl *domain.Template) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO templates (id, name, description, built_in, created_at) VALUES (?, ?, ?, ?, ?)`,
		tmpl.ID, tmpl.Name, tmpl.Description, tmpl.BuiltIn, tmpl.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert template: %w", err)
	}

	for _, m := range tmpl.Metrics {
		_, err = tx.Exec(
			`INSERT INTO template_metrics (id, template_id, name, description_good, description_bad, sort_order)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			m.ID, tmpl.ID, m.Name, m.DescriptionGood, m.DescriptionBad, m.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert template metric %q: %w", m.Name, err)
		}
	}

	return tx.Commit()
}

func (s *Store) FindTemplateByID(id string) (*domain.Template, error) {
	tmpl := &domain.Template{}
	err := s.db.QueryRow(
		`SELECT id, name, description, built_in, created_at FROM templates WHERE id = ?`, id,
	).Scan(&tmpl.ID, &tmpl.Name, &tmpl.Description, &tmpl.BuiltIn, &tmpl.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	metrics, err := s.findTemplateMetrics(id)
	if err != nil {
		return nil, err
	}
	tmpl.Metrics = metrics
	return tmpl, nil
}

func (s *Store) FindTemplateByName(name string) (*domain.Template, error) {
	tmpl := &domain.Template{}
	err := s.db.QueryRow(
		`SELECT id, name, description, built_in, created_at FROM templates WHERE name = ?`, name,
	).Scan(&tmpl.ID, &tmpl.Name, &tmpl.Description, &tmpl.BuiltIn, &tmpl.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	metrics, err := s.findTemplateMetrics(tmpl.ID)
	if err != nil {
		return nil, err
	}
	tmpl.Metrics = metrics
	return tmpl, nil
}

func (s *Store) FindAllTemplates() ([]*domain.Template, error) {
	rows, err := s.db.Query(
		`SELECT id, name, description, built_in, created_at FROM templates ORDER BY built_in DESC, name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.Template
	for rows.Next() {
		t := &domain.Template{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.BuiltIn, &t.CreatedAt); err != nil {
			return nil, err
		}
		metrics, err := s.findTemplateMetrics(t.ID)
		if err != nil {
			return nil, err
		}
		t.Metrics = metrics
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (s *Store) DeleteTemplate(id string) error {
	// Check built-in protection
	var builtIn bool
	err := s.db.QueryRow(`SELECT built_in FROM templates WHERE id = ?`, id).Scan(&builtIn)
	if err == sql.ErrNoRows {
		return fmt.Errorf("template %q not found", id)
	}
	if err != nil {
		return err
	}
	if builtIn {
		return fmt.Errorf("cannot delete built-in template")
	}

	_, err = s.db.Exec(`DELETE FROM templates WHERE id = ?`, id)
	return err
}

func (s *Store) findTemplateMetrics(templateID string) ([]domain.TemplateMetric, error) {
	rows, err := s.db.Query(
		`SELECT id, template_id, name, description_good, description_bad, sort_order
		 FROM template_metrics WHERE template_id = ? ORDER BY sort_order`,
		templateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []domain.TemplateMetric
	for rows.Next() {
		var m domain.TemplateMetric
		if err := rows.Scan(&m.ID, &m.TemplateID, &m.Name, &m.DescriptionGood, &m.DescriptionBad, &m.SortOrder); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	if metrics == nil {
		metrics = []domain.TemplateMetric{}
	}
	return metrics, rows.Err()
}
