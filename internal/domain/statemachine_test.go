package domain_test

import (
	"testing"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
)

func TestCastVote_OnOpenHealthCheck(t *testing.T) {
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}
	metrics := []domain.TemplateMetric{
		{Name: "Fun", DescriptionGood: "Great", DescriptionBad: "Bad"},
	}

	vote, err := hc.CastVote("Fun", "Alice", "green", "", metrics)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vote.MetricName != "Fun" || vote.Participant != "Alice" || vote.Color != domain.VoteGreen {
		t.Errorf("unexpected vote: %+v", vote)
	}
}

func TestCastVote_OnClosedHealthCheck(t *testing.T) {
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusClosed}
	metrics := []domain.TemplateMetric{{Name: "Fun"}}

	_, err := hc.CastVote("Fun", "Alice", "green", "", metrics)
	if err == nil {
		t.Error("expected error on closed health check")
	}
}

func TestCastVote_InvalidMetric(t *testing.T) {
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}
	metrics := []domain.TemplateMetric{{Name: "Fun"}}

	_, err := hc.CastVote("Nonexistent", "Alice", "green", "", metrics)
	if err == nil {
		t.Error("expected error for invalid metric")
	}
}

func TestCastVote_InvalidColor(t *testing.T) {
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}
	metrics := []domain.TemplateMetric{{Name: "Fun"}}

	_, err := hc.CastVote("Fun", "Alice", "blue", "", metrics)
	if err == nil {
		t.Error("expected error for invalid color")
	}
}

func TestComputeOverallScore(t *testing.T) {
	results := []domain.MetricResult{
		{MetricName: "Fun", Score: 3.0, TotalVotes: 3},
		{MetricName: "Speed", Score: 1.0, TotalVotes: 3},
	}
	votes := []*domain.Vote{
		{Participant: "Alice"},
		{Participant: "Bob"},
		{Participant: "Alice"},
	}

	avg, total, names := domain.ComputeOverallScore(results, votes)
	if total != 6 {
		t.Errorf("expected 6 total votes, got %d", total)
	}
	if avg != 2.0 {
		t.Errorf("expected avg 2.0, got %.2f", avg)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 participants, got %d", len(names))
	}
}

func TestGenerateSuggestedActions(t *testing.T) {
	results := []domain.MetricResult{
		{MetricName: "Fun", Score: 3.0, TotalVotes: 3, GreenCount: 3},
		{MetricName: "Speed", Score: 1.5, TotalVotes: 3, GreenCount: 0, RedCount: 2, YellowCount: 1},
	}

	actions := domain.GenerateSuggestedActions(results, "hc-1")
	if len(actions) == 0 {
		t.Error("expected at least one suggested action for low-scoring metric")
	}
	if actions[0].MetricName != "Speed" {
		t.Errorf("expected action for Speed, got %s", actions[0].MetricName)
	}
}
