package mcpui_test

import (
	"strings"
	"testing"

	"github.com/felixgeelhaar/heartbeat/internal/domain"
	"github.com/felixgeelhaar/heartbeat/internal/mcpui"
)

func TestVotingFormHTML(t *testing.T) {
	metrics := []domain.TemplateMetric{
		{Name: "Fun", DescriptionGood: "Great fun", DescriptionBad: "Boring", SortOrder: 1},
		{Name: "Speed", DescriptionGood: "Fast", DescriptionBad: "Slow", SortOrder: 2},
	}

	html := mcpui.VotingFormHTML("hc-123", metrics)

	if !strings.Contains(html, "Fun") {
		t.Error("expected HTML to contain metric name 'Fun'")
	}
	if !strings.Contains(html, "Great fun") {
		t.Error("expected HTML to contain good description")
	}
	if !strings.Contains(html, "Boring") {
		t.Error("expected HTML to contain bad description")
	}
	if !strings.Contains(html, "hc-123") {
		t.Error("expected HTML to contain health check ID")
	}
	if !strings.Contains(html, `value="green"`) {
		t.Error("expected HTML to contain green vote option")
	}
	if !strings.Contains(html, "submitVotes") {
		t.Error("expected HTML to contain submit function")
	}
}

func TestResultsViewHTML(t *testing.T) {
	results := []domain.MetricResult{
		{MetricName: "Fun", GreenCount: 2, YellowCount: 1, RedCount: 0, TotalVotes: 3, Score: 2.67, Comments: []string{}},
		{MetricName: "Speed", GreenCount: 0, YellowCount: 0, RedCount: 2, TotalVotes: 2, Score: 1.0, Comments: []string{}},
	}

	html := mcpui.ResultsViewHTML(results, 1.83, 3)

	if !strings.Contains(html, "Fun") {
		t.Error("expected HTML to contain 'Fun'")
	}
	if !strings.Contains(html, "Speed") {
		t.Error("expected HTML to contain 'Speed'")
	}
	if !strings.Contains(html, "3 participants") {
		t.Error("expected HTML to contain participant count")
	}
	if !strings.Contains(html, "1.8") {
		t.Error("expected HTML to contain average score")
	}
}
