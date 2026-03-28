// Package retro provides a retrospective plugin for healthcheck-mcp.
//
// When enabled, it allows teams to run retrospective discussions linked to
// health check sessions, with three categories: went well, to improve, and action items.
package retro

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felixgeelhaar/go-teamhealthcheck/sdk"
	"github.com/google/uuid"
)

// RetroPlugin implements the retrospective feature as an SDK plugin.
type RetroPlugin struct {
	db     *sql.DB
	store  sdk.StoreReader
	logger sdk.Logger
}

func (p *RetroPlugin) Name() string    { return "retro" }
func (p *RetroPlugin) Version() string { return "1.0.0" }
func (p *RetroPlugin) Description() string {
	return "Retrospective discussions linked to health checks"
}

func (p *RetroPlugin) Init(ctx sdk.PluginContext) error {
	p.db = ctx.DB
	p.store = ctx.Store
	p.logger = ctx.Logger
	return nil
}

func (p *RetroPlugin) Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS retro_sessions (
			id              TEXT PRIMARY KEY,
			healthcheck_id  TEXT NOT NULL,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_retro_sessions_hc ON retro_sessions(healthcheck_id);

		CREATE TABLE IF NOT EXISTS retro_items (
			id              TEXT PRIMARY KEY,
			session_id      TEXT NOT NULL REFERENCES retro_sessions(id) ON DELETE CASCADE,
			category        TEXT NOT NULL CHECK(category IN ('went_well', 'to_improve', 'action_item')),
			text            TEXT NOT NULL,
			author          TEXT NOT NULL DEFAULT '',
			votes           INTEGER NOT NULL DEFAULT 0,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_retro_items_session ON retro_items(session_id);
	`)
	return err
}

func (p *RetroPlugin) RegisterRoutes(reg sdk.RouteRegistry) {
	reg.HandleFunc("GET /api/retro/{hcId}", p.handleGetRetro)
	reg.HandleFunc("POST /api/retro/{hcId}", p.handleCreateRetro)
	reg.HandleFunc("POST /api/retro/{hcId}/items", p.handleAddItem)
	reg.HandleFunc("POST /api/retro/items/{id}/vote", p.handleVoteItem)
}

func (p *RetroPlugin) UIManifest() []sdk.UIEntry {
	return []sdk.UIEntry{{
		Name:   "retro",
		Label:  "Retrospective",
		Icon:   "\U0001f4ac",
		Route:  "/retro/:hcId",
		NavPos: "healthcheck",
	}}
}

func (p *RetroPlugin) OnEvent(e sdk.Event) {
	// Could auto-create retro when HC closes
}

// --- Handlers ---

func (p *RetroPlugin) handleGetRetro(w http.ResponseWriter, r *http.Request) {
	hcID := r.PathValue("hcId")

	var session struct {
		ID            string `json:"id"`
		HealthCheckID string `json:"healthcheck_id"`
		CreatedAt     string `json:"created_at"`
	}

	err := p.db.QueryRow(
		`SELECT id, healthcheck_id, created_at FROM retro_sessions WHERE healthcheck_id = ? LIMIT 1`, hcID,
	).Scan(&session.ID, &session.HealthCheckID, &session.CreatedAt)

	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"session": nil, "items": []any{}})
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, err := p.db.Query(
		`SELECT id, session_id, category, text, author, votes, created_at
		 FROM retro_items WHERE session_id = ? ORDER BY votes DESC, created_at`, session.ID,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type item struct {
		ID        string `json:"id"`
		SessionID string `json:"session_id"`
		Category  string `json:"category"`
		Text      string `json:"text"`
		Author    string `json:"author"`
		Votes     int    `json:"votes"`
		CreatedAt string `json:"created_at"`
	}

	var items []item
	for rows.Next() {
		var it item
		rows.Scan(&it.ID, &it.SessionID, &it.Category, &it.Text, &it.Author, &it.Votes, &it.CreatedAt)
		items = append(items, it)
	}
	if items == nil {
		items = []item{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"session": session, "items": items})
}

func (p *RetroPlugin) handleCreateRetro(w http.ResponseWriter, r *http.Request) {
	hcID := r.PathValue("hcId")

	// Check if retro already exists
	var existingID string
	err := p.db.QueryRow(`SELECT id FROM retro_sessions WHERE healthcheck_id = ?`, hcID).Scan(&existingID)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": existingID, "status": "exists"})
		return
	}

	id := uuid.NewString()
	_, err = p.db.Exec(
		`INSERT INTO retro_sessions (id, healthcheck_id, created_at) VALUES (?, ?, ?)`,
		id, hcID, time.Now(),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "created"})
}

func (p *RetroPlugin) handleAddItem(w http.ResponseWriter, r *http.Request) {
	hcID := r.PathValue("hcId")

	var req struct {
		Category string `json:"category"`
		Text     string `json:"text"`
		Author   string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}
	if req.Category == "" || req.Text == "" {
		http.Error(w, "category and text required", 400)
		return
	}

	// Find or create session
	var sessionID string
	err := p.db.QueryRow(`SELECT id FROM retro_sessions WHERE healthcheck_id = ?`, hcID).Scan(&sessionID)
	if err == sql.ErrNoRows {
		sessionID = uuid.NewString()
		p.db.Exec(`INSERT INTO retro_sessions (id, healthcheck_id, created_at) VALUES (?, ?, ?)`,
			sessionID, hcID, time.Now())
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id := uuid.NewString()
	_, err = p.db.Exec(
		`INSERT INTO retro_items (id, session_id, category, text, author, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		id, sessionID, req.Category, req.Text, req.Author, time.Now(),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "created"})
}

func (p *RetroPlugin) handleVoteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res, err := p.db.Exec(`UPDATE retro_items SET votes = votes + 1 WHERE id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, fmt.Sprintf("item %q not found", id), 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "voted"})
}
