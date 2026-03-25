package domain

import (
	"math"
	"testing"
)

func TestComputeMetricResults(t *testing.T) {
	metrics := []TemplateMetric{
		{Name: "Fun", DescriptionGood: "Great fun", DescriptionBad: "Boring", SortOrder: 1},
		{Name: "Speed", DescriptionGood: "Fast", DescriptionBad: "Slow", SortOrder: 2},
	}

	votes := []*Vote{
		{MetricName: "Fun", Participant: "Alice", Color: VoteGreen},
		{MetricName: "Fun", Participant: "Bob", Color: VoteGreen},
		{MetricName: "Fun", Participant: "Carol", Color: VoteYellow, Comment: "could be better"},
		{MetricName: "Speed", Participant: "Alice", Color: VoteRed, Comment: "too slow"},
		{MetricName: "Speed", Participant: "Bob", Color: VoteRed},
	}

	results := ComputeMetricResults(votes, metrics)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Fun: 2 green + 1 yellow = (3+3+2)/3 = 2.667
	fun := results[0]
	if fun.MetricName != "Fun" {
		t.Errorf("expected first result to be Fun, got %s", fun.MetricName)
	}
	if fun.GreenCount != 2 || fun.YellowCount != 1 || fun.RedCount != 0 {
		t.Errorf("unexpected Fun counts: G=%d Y=%d R=%d", fun.GreenCount, fun.YellowCount, fun.RedCount)
	}
	if fun.TotalVotes != 3 {
		t.Errorf("expected 3 total votes, got %d", fun.TotalVotes)
	}
	expectedScore := (3.0 + 3.0 + 2.0) / 3.0
	if math.Abs(fun.Score-expectedScore) > 0.001 {
		t.Errorf("expected score %.3f, got %.3f", expectedScore, fun.Score)
	}
	if len(fun.Comments) != 1 || fun.Comments[0] != "could be better" {
		t.Errorf("unexpected Fun comments: %v", fun.Comments)
	}

	// Speed: 2 red = (1+1)/2 = 1.0
	speed := results[1]
	if speed.RedCount != 2 || speed.TotalVotes != 2 {
		t.Errorf("unexpected Speed counts: R=%d Total=%d", speed.RedCount, speed.TotalVotes)
	}
	if speed.Score != 1.0 {
		t.Errorf("expected score 1.0, got %.3f", speed.Score)
	}
}

func TestComputeMetricResults_NoVotes(t *testing.T) {
	metrics := []TemplateMetric{
		{Name: "Fun", SortOrder: 1},
	}

	results := ComputeMetricResults(nil, metrics)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].TotalVotes != 0 || results[0].Score != 0 {
		t.Errorf("expected zero votes and zero score for no votes")
	}
	if results[0].Comments == nil {
		t.Error("expected non-nil empty comments slice")
	}
}

func TestComputeTendency(t *testing.T) {
	tests := []struct {
		delta float64
		want  Tendency
	}{
		{0.5, TendencyImproving},
		{0.21, TendencyImproving},
		{0.2, TendencyStable},
		{0.0, TendencyStable},
		{-0.2, TendencyStable},
		{-0.21, TendencyDeclining},
		{-1.0, TendencyDeclining},
	}

	for _, tt := range tests {
		got := ComputeTendency(tt.delta)
		if got != tt.want {
			t.Errorf("ComputeTendency(%.2f) = %s, want %s", tt.delta, got, tt.want)
		}
	}
}
