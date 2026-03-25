package mcp

import (
	"context"
	"fmt"
	"time"

	bolt "github.com/felixgeelhaar/bolt"
	"github.com/felixgeelhaar/mcp-go"
	"github.com/felixgeelhaar/mcp-go/middleware"
	"github.com/google/uuid"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/storage"
)

type submitVoteInput struct {
	HealthCheckID string `json:"healthcheck_id" jsonschema:"required,description=Health check session ID"`
	MetricName    string `json:"metric_name" jsonschema:"required,description=Name of the metric to vote on"`
	Participant   string `json:"participant,omitempty" jsonschema:"description=Name of the person voting (auto-filled from auth when available)"`
	Color         string `json:"color" jsonschema:"required,description=Vote color: green yellow or red"`
	Comment       string `json:"comment,omitempty" jsonschema:"description=Optional comment explaining the vote"`
}

type getResultsInput struct {
	HealthCheckID string `json:"healthcheck_id" jsonschema:"required,description=Health check session ID"`
}

func registerVoteTools(srv *mcp.Server, store *storage.Store, logger *bolt.Logger) {
	srv.Tool("submit_vote").
		Description("Submit a vote for a metric in an open health check. One vote per participant per metric; re-submitting updates the existing vote.").
		Handler(func(ctx context.Context, in submitVoteInput) (any, error) {
			// Auto-fill participant from auth identity when available
			if in.Participant == "" {
				if identity := middleware.IdentityFromContext(ctx); identity != nil {
					in.Participant = identity.Name
				} else {
					return nil, fmt.Errorf("participant is required (no auth identity available)")
				}
			}

			// Validate health check exists and is open
			hc, err := store.FindHealthCheckByID(in.HealthCheckID)
			if err != nil {
				return nil, err
			}
			if hc == nil {
				return nil, fmt.Errorf("health check %q not found", in.HealthCheckID)
			}
			if !hc.IsVotable() {
				return nil, fmt.Errorf("health check %q is not accepting votes (status: %s)", hc.ID, hc.Status)
			}

			// Validate metric exists in template
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil {
				return nil, err
			}
			validMetric := false
			for _, m := range tmpl.Metrics {
				if m.Name == in.MetricName {
					validMetric = true
					break
				}
			}
			if !validMetric {
				return nil, fmt.Errorf("metric %q not found in template %q", in.MetricName, tmpl.Name)
			}

			// Validate color
			color, err := domain.ParseVoteColor(in.Color)
			if err != nil {
				return nil, err
			}

			vote := &domain.Vote{
				ID:            uuid.NewString(),
				HealthCheckID: in.HealthCheckID,
				MetricName:    in.MetricName,
				Participant:   in.Participant,
				Color:         color,
				Comment:       in.Comment,
				CreatedAt:     time.Now(),
			}

			if err := store.UpsertVote(vote); err != nil {
				return nil, fmt.Errorf("submit vote: %w", err)
			}

			return vote, nil
		})

	srv.Tool("get_results").
		Description("Get aggregated results for a health check: per-metric breakdown of green/yellow/red counts, computed score (1-3), and all comments").
		Handler(func(ctx context.Context, in getResultsInput) (any, error) {
			hc, err := store.FindHealthCheckByID(in.HealthCheckID)
			if err != nil {
				return nil, err
			}
			if hc == nil {
				return nil, fmt.Errorf("health check %q not found", in.HealthCheckID)
			}

			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil {
				return nil, err
			}

			votes, err := store.FindVotesByHealthCheck(in.HealthCheckID)
			if err != nil {
				return nil, err
			}

			results := domain.ComputeMetricResults(votes, tmpl.Metrics)

			// Compute overall stats
			var totalScore float64
			var totalVotes int
			participants := make(map[string]bool)
			for _, v := range votes {
				participants[v.Participant] = true
			}
			for _, r := range results {
				totalScore += r.Score * float64(r.TotalVotes)
				totalVotes += r.TotalVotes
			}
			var avgScore float64
			if totalVotes > 0 {
				avgScore = totalScore / float64(totalVotes)
			}

			return map[string]any{
				"healthcheck":   hc,
				"results":       results,
				"average_score": avgScore,
				"participants":  len(participants),
				"total_votes":   totalVotes,
			}, nil
		})
}
