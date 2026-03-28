// Package sdk provides the plugin interface for extending healthcheck-mcp.
//
// Plugins implement the Plugin interface and optionally one or more of:
// ToolProvider, RouteProvider, EventListener, Migrator, UIProvider.
//
// Example:
//
//	type MyPlugin struct{ db *sql.DB }
//
//	func (p *MyPlugin) Name() string        { return "myplugin" }
//	func (p *MyPlugin) Version() string     { return "1.0.0" }
//	func (p *MyPlugin) Description() string { return "My custom plugin" }
//	func (p *MyPlugin) Init(ctx sdk.PluginContext) error {
//	    p.db = ctx.DB
//	    return nil
//	}
//
//	func init() { sdk.Register(&MyPlugin{}) }
package sdk

import (
	"database/sql"
	"net/http"
)

// Plugin is the core interface that all plugins must implement.
type Plugin interface {
	// Name returns the unique identifier for this plugin (used in config).
	Name() string
	// Version returns the semantic version of this plugin.
	Version() string
	// Description returns a human-readable description.
	Description() string
	// Init is called once at startup with access to core services.
	Init(ctx PluginContext) error
}

// ToolProvider is implemented by plugins that register MCP tools.
type ToolProvider interface {
	RegisterTools(reg ToolRegistry)
}

// RouteProvider is implemented by plugins that add REST API endpoints.
type RouteProvider interface {
	RegisterRoutes(reg RouteRegistry)
}

// EventListener is implemented by plugins that react to health check events.
type EventListener interface {
	OnEvent(event Event)
}

// Migrator is implemented by plugins that need database tables.
type Migrator interface {
	Migrate(db *sql.DB) error
}

// UIProvider is implemented by plugins that add pages to the dashboard.
type UIProvider interface {
	UIManifest() []UIEntry
}

// UIEntry describes a plugin page for the frontend.
type UIEntry struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Icon   string `json:"icon"`
	Route  string `json:"route"`
	NavPos string `json:"nav_pos"` // "main" = top nav, "healthcheck" = HC detail page
}

// ToolRegistry allows plugins to register MCP tools.
// The underlying implementation wraps *mcp.Server.
type ToolRegistry interface {
	RegisterPluginTool(name, description string, handler any)
}

// RouteRegistry allows plugins to register HTTP routes.
// The underlying implementation is *http.ServeMux.
type RouteRegistry interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
