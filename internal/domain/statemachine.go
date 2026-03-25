package domain

import (
	"fmt"
	"time"

	bolt "github.com/felixgeelhaar/bolt"
	"github.com/felixgeelhaar/statekit"
)

// HealthCheckContext is the statekit context for health check state transitions.
type HealthCheckContext struct {
	HealthCheck *HealthCheck
	VoteCount   int
}

// Event types for health check state machine.
const (
	EventClose   statekit.EventType = "CLOSE"
	EventReopen  statekit.EventType = "REOPEN"
	EventArchive statekit.EventType = "ARCHIVE"
)

// HealthCheckStateMachine manages the lifecycle of health check sessions.
type HealthCheckStateMachine struct {
	openMachine   *statekit.MachineConfig[HealthCheckContext]
	closedMachine *statekit.MachineConfig[HealthCheckContext]
	logger        *bolt.Logger
}

// NewHealthCheckStateMachine builds the state machine configs for health check lifecycle.
// Two separate machines handle transitions from different starting states since the
// regular interpreter always starts at the initial state.
func NewHealthCheckStateMachine(logger *bolt.Logger) (*HealthCheckStateMachine, error) {
	// Machine for open health checks: open → closed
	openMachine, err := statekit.NewMachine[HealthCheckContext]("healthcheck-open").
		WithInitial("open").
		WithGuard("hasVotes", func(ctx HealthCheckContext, e statekit.Event) bool {
			return ctx.VoteCount > 0
		}).
		WithAction("setClosedAt", func(ctx *HealthCheckContext, e statekit.Event) {
			now := time.Now()
			ctx.HealthCheck.Status = StatusClosed
			ctx.HealthCheck.ClosedAt = &now
		}).
		State("open").
		On("CLOSE").Target("closed").Guard("hasVotes").Do("setClosedAt").Done().
		State("closed").Final().Done().
		Build()
	if err != nil {
		return nil, fmt.Errorf("build open machine: %w", err)
	}

	// Machine for closed health checks: closed → archived or → open (reopen)
	closedMachine, err := statekit.NewMachine[HealthCheckContext]("healthcheck-closed").
		WithInitial("closed").
		WithAction("setArchived", func(ctx *HealthCheckContext, e statekit.Event) {
			ctx.HealthCheck.Status = StatusArchived
		}).
		WithAction("clearClosedAt", func(ctx *HealthCheckContext, e statekit.Event) {
			ctx.HealthCheck.Status = StatusOpen
			ctx.HealthCheck.ClosedAt = nil
		}).
		State("closed").
		On("ARCHIVE").Target("archived").Do("setArchived").
		On("REOPEN").Target("reopened").Do("clearClosedAt").Done().
		State("archived").Final().Done().
		State("reopened").Final().Done().
		Build()
	if err != nil {
		return nil, fmt.Errorf("build closed machine: %w", err)
	}

	return &HealthCheckStateMachine{
		openMachine:   openMachine,
		closedMachine: closedMachine,
		logger:        logger,
	}, nil
}

// Transition validates and executes a state transition on a health check.
func (sm *HealthCheckStateMachine) Transition(hc *HealthCheck, event statekit.EventType, voteCount int) error {
	var machine *statekit.MachineConfig[HealthCheckContext]
	switch hc.Status {
	case StatusOpen:
		machine = sm.openMachine
	case StatusClosed:
		machine = sm.closedMachine
	case StatusArchived:
		return fmt.Errorf("health check %q is archived, no transitions allowed", hc.ID)
	default:
		return fmt.Errorf("unknown status %q for health check %q", hc.Status, hc.ID)
	}

	interp := statekit.NewInterpreter(machine)
	interp.UpdateContext(func(ctx *HealthCheckContext) {
		ctx.HealthCheck = hc
		ctx.VoteCount = voteCount
	})
	interp.Start()

	prevStatus := hc.Status
	interp.Send(statekit.Event{Type: event})

	// If status didn't change, the transition was rejected (guard failed or no matching event)
	if hc.Status == prevStatus {
		return fmt.Errorf("transition %q not allowed from state %q (guard condition not met)", event, prevStatus)
	}

	sm.logger.Info().
		Str("healthcheck_id", hc.ID).
		Str("event", string(event)).
		Str("from", string(prevStatus)).
		Str("to", string(hc.Status)).
		Msg("health check state transition")

	return nil
}
