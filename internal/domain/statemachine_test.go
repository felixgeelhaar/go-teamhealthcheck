package domain_test

import (
	"os"
	"testing"

	bolt "github.com/felixgeelhaar/bolt"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
)

func newStateMachine(t *testing.T) *domain.HealthCheckStateMachine {
	t.Helper()
	logger := bolt.New(bolt.NewConsoleHandler(os.Stderr))
	sm, err := domain.NewHealthCheckStateMachine(logger)
	if err != nil {
		t.Fatalf("build state machine: %v", err)
	}
	return sm
}

func TestTransition_OpenToClosedWithVotes(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}

	err := sm.Transition(hc, domain.EventClose, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hc.Status != domain.StatusClosed {
		t.Errorf("expected closed, got %s", hc.Status)
	}
	if hc.ClosedAt == nil {
		t.Error("expected ClosedAt to be set")
	}
}

func TestTransition_OpenToClosedWithoutVotes(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}

	err := sm.Transition(hc, domain.EventClose, 0)
	if err == nil {
		t.Error("expected error when closing with no votes")
	}
	if hc.Status != domain.StatusOpen {
		t.Errorf("expected status to remain open, got %s", hc.Status)
	}
}

func TestTransition_ClosedToArchived(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusClosed}

	err := sm.Transition(hc, domain.EventArchive, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hc.Status != domain.StatusArchived {
		t.Errorf("expected archived, got %s", hc.Status)
	}
}

func TestTransition_ClosedToReopened(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusClosed}

	err := sm.Transition(hc, domain.EventReopen, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hc.Status != domain.StatusOpen {
		t.Errorf("expected open after reopen, got %s", hc.Status)
	}
	if hc.ClosedAt != nil {
		t.Error("expected ClosedAt to be nil after reopen")
	}
}

func TestTransition_ArchivedBlocked(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusArchived}

	err := sm.Transition(hc, domain.EventClose, 0)
	if err == nil {
		t.Error("expected error for transition from archived")
	}
}

func TestTransition_OpenCannotArchive(t *testing.T) {
	sm := newStateMachine(t)
	hc := &domain.HealthCheck{ID: "hc-1", Status: domain.StatusOpen}

	err := sm.Transition(hc, domain.EventArchive, 0)
	if err == nil {
		t.Error("expected error: cannot archive an open health check")
	}
	if hc.Status != domain.StatusOpen {
		t.Errorf("expected status to remain open, got %s", hc.Status)
	}
}
