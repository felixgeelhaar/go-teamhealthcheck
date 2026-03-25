# healthcheck-mcp

An MCP (Model Context Protocol) server for running Spotify Squad Health Checks with AI agents.

Let AI facilitate team health checks directly in conversations — collecting votes, aggregating results, tracking trends, and surfacing discussion topics. Supports both single-user (stdio) and multi-user (HTTP/SSE) modes.

## Features

- **24 MCP tools** — Full health check lifecycle: teams, templates, sessions, voting, comparison, and analysis
- **Multi-user HTTP/SSE** — Multiple team members connect to a shared server, each voting independently
- **Token-based auth** — Bearer token authentication with auto-filled participant identity
- **State machine lifecycle** — statekit-powered transitions: open → closed → archived, with guards (can't close without votes) and reopen support
- **Built-in Spotify template** — Ships with the original 10 Squad Health Check categories
- **Custom templates** — Create your own health check formats
- **Trend tracking** — Compare results across sprints, flag declining metrics
- **AI-friendly analysis** — Structured summaries, discussion topic generation, disagreement detection
- **Structured logging** — bolt-powered JSON (prod) or colored console (dev) logging
- **SQLite storage** — Zero-config persistence, single-file database
- **Single binary** — No runtime dependencies, cross-platform

## Installation

### From source

```bash
go install github.com/felixgeelhaar/go-teamhealthcheck/cmd/healthcheck-mcp@latest
```

### From release

Download a pre-built binary from [Releases](https://github.com/felixgeelhaar/go-teamhealthcheck/releases).

## Usage

### Single-user (stdio)

```bash
healthcheck-mcp
```

For use with Claude Desktop, add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "healthcheck": {
      "command": "healthcheck-mcp",
      "args": []
    }
  }
}
```

### Multi-user (HTTP/SSE)

```bash
healthcheck-mcp --mode http --addr :8080
```

Multiple team members connect their MCP clients to the same server. Each authenticates with a Bearer token that maps to their identity.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--mode` | `stdio` | Transport: `stdio` or `http` |
| `--addr` | `:8080` | HTTP listen address |
| `--db` | `~/.healthcheck-mcp/data.db` | SQLite database path |
| `--auth` | `~/.healthcheck-mcp/auth.json` | Auth config file (HTTP mode) |
| `--dev` | `false` | Development mode (colored console logging) |

### Auth Configuration

Create `~/.healthcheck-mcp/auth.json`:

```json
{
  "tokens": {
    "alice-secret-token": {"name": "Alice", "team_id": "your-team-uuid"},
    "bob-secret-token": {"name": "Bob", "team_id": "your-team-uuid"}
  }
}
```

When authenticated, `submit_vote` auto-fills the participant name from the token identity. The `my_pending_healthchecks` tool shows which metrics the user hasn't voted on yet.

## MCP Tools

### Team Management
| Tool | Description |
|------|-------------|
| `create_team` | Create a new team with optional initial members |
| `list_teams` | List all teams |
| `get_team` | Get team details including members |
| `delete_team` | Delete a team |
| `add_team_member` | Add a member to a team |
| `remove_team_member` | Remove a member from a team |

### Templates
| Tool | Description |
|------|-------------|
| `list_templates` | List all templates (includes built-in Spotify template) |
| `get_template` | Get template with all metric definitions |
| `create_template` | Create a custom template with metrics |
| `delete_template` | Delete a custom template |

### Health Check Sessions
| Tool | Description |
|------|-------------|
| `create_healthcheck` | Start a new health check session |
| `list_healthchecks` | List sessions (filter by team/status) |
| `get_healthcheck` | Get session with current results |
| `close_healthcheck` | Close session (requires at least one vote) |
| `reopen_healthcheck` | Reopen a closed session for more votes |
| `archive_healthcheck` | Archive a closed session (terminal) |
| `delete_healthcheck` | Delete session and all votes |
| `my_pending_healthchecks` | List metrics the authenticated user hasn't voted on |

### Voting
| Tool | Description |
|------|-------------|
| `submit_vote` | Vote green/yellow/red on a metric (participant auto-filled from auth) |
| `get_results` | Get aggregated results with scores and stats |

### Analysis
| Tool | Description |
|------|-------------|
| `compare_sessions` | Compare results across sprints with trend detection |
| `analyze_healthcheck` | AI-friendly summary with strengths/concerns |
| `get_trends` | Historical trend analysis for a team |
| `get_discussion_topics` | Suggested topics based on disagreement, low scores, and declining trends |

## Health Check Lifecycle

The session lifecycle is managed by a [statekit](https://github.com/felixgeelhaar/statekit) state machine:

```
         create              close (requires votes)        archive
         ──────→  open  ─────────────────────────→  closed  ──────→  archived
                   ↑                                  │
                   └──────────── reopen ──────────────┘
```

- **open** — accepting votes
- **closed** — voting complete, results available, can reopen or archive
- **archived** — terminal state, read-only

Guards enforce business rules (e.g., can't close without votes). Actions execute side effects (set timestamps, log transitions).

## Multi-User Workflow

```
# Team lead starts the server
healthcheck-mcp --mode http --addr :8080

# Each team member connects their AI client with their token
# Alice's session:
Alice: "Do I have any pending health checks?"
Agent: [calls my_pending_healthchecks] You have Sprint 42 — 10 metrics pending.
Alice: "I vote green on Fun, red on Tech Quality..."
Agent: [calls submit_vote x10] All votes recorded!

# Bob's session (same server):
Bob: "What health checks are open?"
Agent: [calls my_pending_healthchecks] Sprint 42 — 10 metrics pending.
Bob: "Green on everything except Speed, that's yellow"
Agent: [calls submit_vote x10] Done!

# Team lead reviews:
Lead: "Show me the Sprint 42 results"
Agent: [calls get_results] Here's the breakdown...
Lead: "What should we discuss?"
Agent: [calls get_discussion_topics] Top topics: Tech Quality (disagreement)...
```

## Built-in Spotify Template

Ships with the original [Spotify Squad Health Check](https://labs.spotify.com/2014/09/16/squad-health-check-model/) categories:

1. Easy to Release
2. Suitable Process
3. Tech Quality
4. Value
5. Speed
6. Mission
7. Fun
8. Learning
9. Support
10. Pawns or Players

## Built With

- [mcp-go](https://github.com/felixgeelhaar/mcp-go) — MCP server framework
- [bolt](https://github.com/felixgeelhaar/bolt) — Structured logging
- [statekit](https://github.com/felixgeelhaar/statekit) — State machine engine
- [fortify](https://github.com/felixgeelhaar/fortify) — Resilience middleware (via mcp-go)

## Commit Convention

Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `perf:`, `test:`.

## License

MIT
