package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

var v7TestCtx = context.Background()

// This builds a database frozen at migration v6: the papers table has
// no location/category columns and one legacy paper exists. It then
// runs ONLY v7 and checks the old row survived with empty values.
func TestV7LegacyRows(t *testing.T) {
	db := sqlx.MustConnect("sqlite", "file::memory:")
	defer db.Close()

	db.MustExec(`
    CREATE TABLE migrations (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      version BIGINT UNIQUE NOT NULL,
      applied_at BIGINT NOT NULL
    );
  `)
	db.MustExec(`INSERT INTO migrations (version, applied_at) VALUES (6, 0);`)

	// Papers schema as created by v2 (no location/category yet).
	db.MustExec(`
    CREATE TABLE papers (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      uuid TEXT NOT NULL UNIQUE CHECK (uuid != ''),
      title TEXT NOT NULL CHECK (title != ''),
      number INTEGER,
      filename TEXT NOT NULL UNIQUE CHECK (filename != ''),
      volume INTEGER CHECK (volume > 0),
      issue INTEGER CHECK (issue > 0),

      owner_uuid TEXT REFERENCES users(uuid) ON DELETE SET NULL,

      UNIQUE(title, owner_uuid),
      UNIQUE(volume, number, issue)
    );
  `)
	db.MustExec(`
    INSERT INTO papers (uuid, title, filename)
    VALUES ('11111111-1111-1111-1111-111111111111', 'Legacy', 'legacy.pdf');
  `)

	tx, err := db.BeginTxx(v7TestCtx, nil)
	if err != nil {
		t.Fatalf("Could not begin tx: %v", err)
	}
	if err := (v7{}).Apply(v7TestCtx, tx); err != nil {
		t.Fatalf("v7 failed on legacy database: %v", err)
	}

	var location, category string
	if err := db.GetContext(v7TestCtx, &location,
		`SELECT location FROM papers WHERE uuid = '11111111-1111-1111-1111-111111111111';`,
	); err != nil {
		t.Fatalf("location column missing after v7: %v", err)
	}
	if err := db.GetContext(v7TestCtx, &category,
		`SELECT category FROM papers WHERE uuid = '11111111-1111-1111-1111-111111111111';`,
	); err != nil {
		t.Fatalf("category column missing after v7: %v", err)
	}
	if location != "" || category != "" {
		t.Fatalf("Legacy row should keep empty values, got %q, %q", location, category)
	}

	// The v7 length cap must hold at the database level too.
	if _, err := db.ExecContext(v7TestCtx,
		`INSERT INTO papers (uuid, title, filename, location, category)
		 VALUES ('22222222-2222-2222-2222-222222222222', 'Too long', 'toolong.pdf', ?, 'x');`,
		strings.Repeat("a", 121),
	); err == nil {
		t.Fatal("Expected CHECK rejection for 121-character location")
	}

	var version uint64
	if err := db.GetContext(v7TestCtx, &version,
		`SELECT MAX(version) FROM migrations;`,
	); err != nil || version != 7 {
		t.Fatalf("Expected migrations to record version 7, got %v, %v", version, err)
	}
}
