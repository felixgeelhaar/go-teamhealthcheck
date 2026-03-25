package domain

import "testing"

func TestHealthCheck_IsOpen(t *testing.T) {
	hc := &HealthCheck{ID: "hc-1", Status: StatusOpen}
	if !hc.IsOpen() {
		t.Error("expected IsOpen() = true for open health check")
	}
	if !hc.IsVotable() {
		t.Error("expected IsVotable() = true for open health check")
	}

	hc.Status = StatusClosed
	if hc.IsOpen() {
		t.Error("expected IsOpen() = false for closed health check")
	}
	if hc.IsVotable() {
		t.Error("expected IsVotable() = false for closed health check")
	}

	hc.Status = StatusArchived
	if hc.IsOpen() {
		t.Error("expected IsOpen() = false for archived health check")
	}
	if hc.IsVotable() {
		t.Error("expected IsVotable() = false for archived health check")
	}
}
