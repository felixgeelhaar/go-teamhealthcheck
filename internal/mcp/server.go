package mcp

import (
	bolt "github.com/felixgeelhaar/bolt"
	"github.com/felixgeelhaar/mcp-go"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
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

	return srv
}
