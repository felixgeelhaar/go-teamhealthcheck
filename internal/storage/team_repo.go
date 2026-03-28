package storage

import (
	"database/sql"
	"fmt"

	"github.com/felixgeelhaar/heartbeat/internal/domain"
)

func (s *Store) CreateTeam(team *domain.Team) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO teams (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		team.ID, team.Name, team.CreatedAt, team.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	for _, m := range team.Members {
		_, err = tx.Exec(
			`INSERT INTO team_members (id, team_id, name) VALUES (?, ?, ?)`,
			newID(), team.ID, m,
		)
		if err != nil {
			return fmt.Errorf("insert member %q: %w", m, err)
		}
	}

	return tx.Commit()
}

func (s *Store) FindTeamByID(id string) (*domain.Team, error) {
	team := &domain.Team{}
	err := s.db.QueryRow(
		`SELECT id, name, created_at, updated_at FROM teams WHERE id = ?`, id,
	).Scan(&team.ID, &team.Name, &team.CreatedAt, &team.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	members, err := s.findTeamMembers(id)
	if err != nil {
		return nil, err
	}
	team.Members = members
	return team, nil
}

func (s *Store) FindAllTeams() ([]*domain.Team, error) {
	rows, err := s.db.Query(`SELECT id, name, created_at, updated_at FROM teams ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		t := &domain.Team{}
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		members, err := s.findTeamMembers(t.ID)
		if err != nil {
			return nil, err
		}
		t.Members = members
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (s *Store) DeleteTeam(id string) error {
	res, err := s.db.Exec(`DELETE FROM teams WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("team %q not found", id)
	}
	return nil
}

func (s *Store) AddTeamMember(teamID, name string) error {
	_, err := s.db.Exec(
		`INSERT INTO team_members (id, team_id, name) VALUES (?, ?, ?)`,
		newID(), teamID, name,
	)
	return err
}

func (s *Store) RemoveTeamMember(teamID, name string) error {
	res, err := s.db.Exec(
		`DELETE FROM team_members WHERE team_id = ? AND name = ?`,
		teamID, name,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("member %q not found in team %q", name, teamID)
	}
	return nil
}

func (s *Store) findTeamMembers(teamID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT name FROM team_members WHERE team_id = ? ORDER BY name`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		members = append(members, name)
	}
	if members == nil {
		members = []string{}
	}
	return members, rows.Err()
}
