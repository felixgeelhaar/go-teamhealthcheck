package storage

import (
	"github.com/felixgeelhaar/heartbeat/internal/domain"
	"github.com/felixgeelhaar/heartbeat/internal/events"
)

func (s *Store) UpsertVote(vote *domain.Vote) error {
	_, err := s.db.Exec(
		`INSERT INTO votes (id, healthcheck_id, metric_name, participant, color, comment, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(healthcheck_id, metric_name, participant)
		 DO UPDATE SET color = excluded.color, comment = excluded.comment, created_at = excluded.created_at`,
		vote.ID, vote.HealthCheckID, vote.MetricName, vote.Participant, vote.Color, vote.Comment, vote.CreatedAt,
	)
	if err == nil {
		s.publish(events.Event{
			Type:          events.VoteSubmitted,
			HealthCheckID: vote.HealthCheckID,
			Participant:   vote.Participant,
			MetricName:    vote.MetricName,
		})
	}
	return err
}

func (s *Store) FindVotesByHealthCheck(healthCheckID string) ([]*domain.Vote, error) {
	rows, err := s.db.Query(
		`SELECT id, healthcheck_id, metric_name, participant, color, comment, created_at
		 FROM votes WHERE healthcheck_id = ? ORDER BY metric_name, participant`,
		healthCheckID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []*domain.Vote
	for rows.Next() {
		v := &domain.Vote{}
		if err := rows.Scan(&v.ID, &v.HealthCheckID, &v.MetricName, &v.Participant, &v.Color, &v.Comment, &v.CreatedAt); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, rows.Err()
}
