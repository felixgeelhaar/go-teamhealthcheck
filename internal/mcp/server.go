package mcp

import (
	"context"
	"fmt"

	bolt "github.com/felixgeelhaar/bolt"
	"github.com/felixgeelhaar/mcp-go"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/mcpui"
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/storage"
)

// NewServer creates a fully configured MCP server with all health check tools registered.
func NewServer(store *storage.Store, logger *bolt.Logger) *mcp.Server {
	sm, err := domain.NewHealthCheckStateMachine(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to build health check state machine")
	}

	srv := mcp.NewServer(mcp.ServerInfo{
		Name:    "healthcheck-mcp",
		Version: "1.0.0",
	})

	registerTeamTools(srv, store, logger)
	registerTemplateTools(srv, store, logger)
	registerHealthCheckTools(srv, store, logger, sm)
	registerVoteTools(srv, store, logger)
	registerCompareTools(srv, store, logger)
	registerAnalyzeTools(srv, store, logger)
	registerUIResources(srv, store)

	return srv
}

func registerUIResources(srv *mcp.Server, store *storage.Store) {
	// Voting form UI — renders all metrics for a health check
	srv.Resource("ui://healthcheck/{id}/vote").
		Name("Health Check Voting Form").
		Description("Interactive voting form for a health check session").
		MimeType("text/html;profile=mcp-app").
		Handler(func(ctx context.Context, uri string, params map[string]string) (*mcp.ResourceContent, error) {
			hcID := params["id"]
			hc, err := store.FindHealthCheckByID(hcID)
			if err != nil || hc == nil {
				return nil, fmt.Errorf("health check not found")
			}
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil || tmpl == nil {
				return nil, fmt.Errorf("template not found")
			}
			return &mcp.ResourceContent{
				URI:      uri,
				MimeType: "text/html;profile=mcp-app",
				Text:     mcpui.VotingFormHTML(hcID, tmpl.Metrics),
			}, nil
		})

	// Results view UI — renders traffic-light heatmap
	srv.Resource("ui://healthcheck/{id}/results").
		Name("Health Check Results").
		Description("Visual results heatmap for a health check session").
		MimeType("text/html;profile=mcp-app").
		Handler(func(ctx context.Context, uri string, params map[string]string) (*mcp.ResourceContent, error) {
			hcID := params["id"]
			hc, err := store.FindHealthCheckByID(hcID)
			if err != nil || hc == nil {
				return nil, fmt.Errorf("health check not found")
			}
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil || tmpl == nil {
				return nil, fmt.Errorf("template not found")
			}
			votes, err := store.FindVotesByHealthCheck(hcID)
			if err != nil {
				return nil, err
			}
			results := domain.ComputeMetricResults(votes, tmpl.Metrics)

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

			return &mcp.ResourceContent{
				URI:      uri,
				MimeType: "text/html;profile=mcp-app",
				Text:     mcpui.ResultsViewHTML(results, avgScore, len(participants)),
			}, nil
		})
}
