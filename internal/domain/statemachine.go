package domain

import "fmt"

// LifecycleEvent represents a state transition event for a health check.
type LifecycleEvent string

const (
	EventClose   LifecycleEvent = "CLOSE"
	EventReopen  LifecycleEvent = "REOPEN"
	EventArchive LifecycleEvent = "ARCHIVE"
)

// HealthCheckLifecycle manages state transitions for health check sessions.
// The implementation lives outside the domain (e.g., internal/lifecycle).
type HealthCheckLifecycle interface {
	Transition(hc *HealthCheck, event LifecycleEvent, voteCount int) error
}

// --- Domain services ---

// CastVote validates and creates a Vote on a HealthCheck.
// This consolidates the vote-casting invariants in the aggregate root.
func (hc *HealthCheck) CastVote(metricName, participant, color, comment string, templateMetrics []TemplateMetric) (*Vote, error) {
	if !hc.IsVotable() {
		return nil, NewDomainError("health check %q is not accepting votes (status: %s)", hc.ID, hc.Status)
	}

	// Validate metric exists in template
	valid := false
	for _, m := range templateMetrics {
		if m.Name == metricName {
			valid = true
			break
		}
	}
	if !valid {
		return nil, NewDomainError("metric %q not found in template", metricName)
	}

	// Validate color
	vc, err := ParseVoteColor(color)
	if err != nil {
		return nil, err
	}

	return &Vote{
		HealthCheckID: hc.ID,
		MetricName:    metricName,
		Participant:   participant,
		Color:         vc,
		Comment:       comment,
	}, nil
}

// ComputeOverallScore computes the weighted average score across all metric results.
func ComputeOverallScore(results []MetricResult, votes []*Vote) (avgScore float64, totalVotes int, participantNames []string) {
	var totalScore float64
	participantSet := make(map[string]bool)

	for _, v := range votes {
		if !participantSet[v.Participant] {
			participantSet[v.Participant] = true
			participantNames = append(participantNames, v.Participant)
		}
	}
	for _, r := range results {
		totalScore += r.Score * float64(r.TotalVotes)
		totalVotes += r.TotalVotes
	}
	if totalVotes > 0 {
		avgScore = totalScore / float64(totalVotes)
	}
	return
}

// Alert severity constants.
const (
	AlertSeverityCritical = "critical"
	AlertSeverityWarning  = "warning"

	AlertThresholdCritical = 1.5
	AlertThresholdWarning  = 2.0
)

// GenerateSuggestedActions creates action items for metrics with low scores or high disagreement.
func GenerateSuggestedActions(results []MetricResult, healthCheckID string) []*Action {
	var actions []*Action
	for _, res := range results {
		if res.TotalVotes == 0 {
			continue
		}

		var desc string

		if res.Score < AlertThresholdWarning {
			desc = NewDomainError("Improve %s: currently scoring %.1f/3.0. Discuss root causes and identify one concrete improvement for next sprint.", res.MetricName, res.Score).Error()
		}

		if desc == "" && res.GreenCount > 0 && res.RedCount > 0 {
			spread := float64(res.GreenCount-res.RedCount) / float64(res.TotalVotes)
			if spread < 0.5 && spread > -0.5 {
				desc = NewDomainError("Discuss %s: team has split opinions (%d green, %d yellow, %d red). Understand different perspectives.", res.MetricName, res.GreenCount, res.YellowCount, res.RedCount).Error()
			}
		}

		if desc != "" {
			actions = append(actions, &Action{
				HealthCheckID: healthCheckID,
				MetricName:    res.MetricName,
				Description:   desc,
			})
		}
	}
	return actions
}

// DomainError represents a business rule violation.
type DomainError struct {
	msg string
}

func NewDomainError(format string, args ...any) *DomainError {
	return &DomainError{msg: fmt.Sprintf(format, args...)}
}

func (e *DomainError) Error() string {
	return e.msg
}
