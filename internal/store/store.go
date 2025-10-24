package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/phoneinfoga-desktop/internal/cfg"
)

// DefaultPath resolves the default on-disk location for the SQLite database file.
func DefaultPath() (string, error) {
	dir, err := cfg.AppConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "app.db"), nil
}

// Client is a placeholder for the future ent-backed store client.
type Client struct{}

// Open currently returns a not-implemented error until the ent store is wired.
func Open(_ context.Context) (*Client, error) {
	return nil, fmt.Errorf("store Open not yet implemented")
}

// EnsurePath guarantees the directory for the database exists.
func EnsurePath(path string) error {
	if path == "" {
		return fmt.Errorf("empty database path")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}
	return nil
}
