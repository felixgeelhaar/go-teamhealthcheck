package seed_test

import (
	"testing"

	"github.com/felixgeelhaar/heartbeat/internal/seed"
)

func TestSpotifyTemplate(t *testing.T) {
	tmpl := seed.SpotifyTemplate()

	if tmpl.Name != "Spotify Squad Health Check" {
		t.Errorf("expected Spotify Squad Health Check, got %s", tmpl.Name)
	}
	if !tmpl.BuiltIn {
		t.Error("expected BuiltIn=true")
	}
	if len(tmpl.Metrics) != 10 {
		t.Errorf("expected 10 metrics, got %d", len(tmpl.Metrics))
	}

	// Verify all metrics have descriptions
	for _, m := range tmpl.Metrics {
		if m.Name == "" {
			t.Error("metric has empty name")
		}
		if m.DescriptionGood == "" {
			t.Errorf("metric %q has empty DescriptionGood", m.Name)
		}
		if m.DescriptionBad == "" {
			t.Errorf("metric %q has empty DescriptionBad", m.Name)
		}
	}

	// Verify sort order is sequential
	for i, m := range tmpl.Metrics {
		if m.SortOrder != i+1 {
			t.Errorf("metric %q has sort order %d, expected %d", m.Name, m.SortOrder, i+1)
		}
	}
}
