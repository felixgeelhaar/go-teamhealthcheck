package auth_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/felixgeelhaar/heartbeat/internal/auth"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.json")

	content := `{
		"tokens": {
			"token-alice": {"name": "Alice", "team_id": "team-1"},
			"token-bob": {"name": "Bob", "team_id": "team-1"}
		}
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	validator, err := auth.LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	// Valid token
	identity := validator("token-alice")
	if identity == nil {
		t.Fatal("expected identity for token-alice")
	}
	if identity.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", identity.Name)
	}
	teamID, _ := identity.Metadata["team_id"].(string)
	if teamID != "team-1" {
		t.Errorf("expected team_id team-1, got %s", teamID)
	}

	// Another valid token
	identity = validator("token-bob")
	if identity == nil {
		t.Fatal("expected identity for token-bob")
	}
	if identity.Name != "Bob" {
		t.Errorf("expected name Bob, got %s", identity.Name)
	}

	// Invalid token
	identity = validator("invalid-token")
	if identity != nil {
		t.Error("expected nil for invalid token")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := auth.LoadConfig("/nonexistent/auth.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.json")
	os.WriteFile(path, []byte("not json"), 0o644)

	_, err := auth.LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
