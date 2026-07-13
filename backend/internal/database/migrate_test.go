package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverMigrationsSortsSQLFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{"000002_second.sql", "notes.txt", "000001_first.sql"}

	for _, file := range files {
		if err := os.WriteFile(filepath.Join(dir, file), []byte("-- test"), 0o600); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
	}

	migrations, err := DiscoverMigrations(dir)
	if err != nil {
		t.Fatalf("expected migrations, got error: %v", err)
	}

	if len(migrations) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migrations))
	}

	if migrations[0].Version != "000001_first" {
		t.Fatalf("expected first migration sorted first, got %q", migrations[0].Version)
	}

	if migrations[1].Version != "000002_second" {
		t.Fatalf("expected second migration sorted second, got %q", migrations[1].Version)
	}
}
