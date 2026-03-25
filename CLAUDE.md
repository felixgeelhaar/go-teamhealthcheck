# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MCP server for running Spotify Squad Health Checks with AI agents. Supports single-user (stdio) and multi-user (HTTP/SSE) modes with token-based auth.

Built with:
- `github.com/felixgeelhaar/mcp-go` — MCP framework
- `github.com/felixgeelhaar/bolt` — Structured logging
- `github.com/felixgeelhaar/statekit` — State machine for health check lifecycle
- `modernc.org/sqlite` — Pure Go SQLite

## Commands

```bash
# Build the binary
go build ./cmd/healthcheck-mcp/

# Run all tests
go test -race ./...

# Run a single test
go test -race -run TestTransition_OpenToClosedWithVotes ./internal/domain/

# Vet
go vet ./...

# Run stdio mode
go run ./cmd/healthcheck-mcp/ --db /tmp/test.db

# Run HTTP mode
go run ./cmd/healthcheck-mcp/ --mode http --addr :8080 --dev
```

## Architecture (DDD)

```
cmd/healthcheck-mcp/main.go       → Entry point: bolt logger, dual transport, auth middleware
internal/auth/config.go            → Auth config loader (token → user identity mapping)
internal/domain/                   → Pure domain: entities, value objects, repository interfaces, state machine
internal/storage/                  → SQLite repository implementations
internal/mcp/                      → MCP tool handlers (application/adapter layer)
internal/seed/                     → Built-in Spotify template data
```

**Domain layer** (`internal/domain/`) has zero infrastructure imports except bolt and statekit. Repository interfaces defined here; `storage` implements them.

**Key domain types:**
- `Team` — aggregate root, owns health checks
- `Template` / `TemplateMetric` — reusable health check formats
- `HealthCheck` — session aggregate with open/closed/archived lifecycle via statekit
- `Vote` — participant's green/yellow/red choice (upsert on re-vote)
- `MetricResult` / `MetricTrend` — computed views, not stored
- `HealthCheckStateMachine` — statekit-powered lifecycle transitions with guards and actions

**State machine:** Two machine configs (one for open→closed, one for closed→archived/reopened) since statekit's regular interpreter always starts at the initial state. The `Transition()` method selects the right machine based on current status.

**Auth flow (HTTP mode):** Bearer token in header → mcp-go auth middleware → `middleware.IdentityFromContext(ctx)` in tool handlers → auto-fills participant name on votes.

**Storage:** SQLite via `modernc.org/sqlite` (pure Go, no CGO). In-memory via `:memory:` with shared cache for tests. WAL mode + busy_timeout for concurrent HTTP access.

**MCP tools (24):** Registered in `internal/mcp/server.go`, grouped by file: team, template, healthcheck, vote, compare, analyze.

## Commit Convention

Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `perf:`, `test:`. Subject max 50 chars, imperative mood.
