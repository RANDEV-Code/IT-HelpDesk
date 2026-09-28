// Package testdb menyediakan harness database integrasi yang dapat diulang:
// membaca DSN dari TEST_DATABASE_URL, menerapkan seluruh migrasi baseline,
// dan membersihkan data antar suite (TEST_PLAN §1, §6).
//
// Kebijakan: bila TEST_DATABASE_URL tidak diset, tes di-SKIP (bukan gagal) agar
// `go test ./...` unit tetap hijau tanpa DB. CI dan dev mengset variabel ini ke
// database terisolasi (randesk_test), BUKAN ke database dev.
package testdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"

	"randesk/backend/migrations"
)

// Connect membuka pool ke TEST_DATABASE_URL dan menerapkan seluruh migrasi.
// Mengembalikan *sql.DB siap pakai. Mendi-skip tes bila env tidak ada.
func Connect(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL tidak diset; lewati tes integrasi PostgreSQL")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open TEST_DATABASE_URL: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		t.Fatalf("ping database test: %v (apakah PostgreSQL tersedia & TEST_DATABASE_URL benar?)", err)
	}

	applyMigrations(t, db)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// applyMigrations menjalankan seluruh .up.sql dari FS tertanam sampai versi
// terbaru. Idempoten: DB yang sudah termigrasi menghasilkan ErrNoChange.
func applyMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		t.Fatalf("sumber migrasi iofs: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{SchemaName: "public", MigrationsTable: "schema_migrations"})
	if err != nil {
		t.Fatalf("driver postgres migrate: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		t.Fatalf("init migrate: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up: %v", err)
	}
}

// TruncateAll menghapus seluruh baris tabel aplikasi publik (CASCADE) dan
// me-reset identity, memberi slate bersih antar tes. Tidak menyentuh
// schema_migrations.
func TruncateAll(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := appTables(t, db)
	if len(tables) == 0 {
		return
	}
	list := make([]string, len(tables))
	for i, tbl := range tables {
		list[i] = "public." + quoteIdent(tbl)
	}
	stmt := "TRUNCATE " + strings.Join(list, ", ") + " RESTART IDENTITY CASCADE"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		t.Fatalf("truncate tabel aplikasi: %v", err)
	}
}

// appTables mendaftar tabel base di schema public kecuali buku migrasi.
func appTables(t *testing.T, db *sql.DB) []string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, `
		SELECT table_name FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
		  AND table_name <> 'schema_migrations'
		ORDER BY table_name`)
	if err != nil {
		t.Fatalf("daftar tabel: %v", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan tabel: %v", err)
		}
		out = append(out, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterasi tabel: %v", err)
	}
	return out
}

func quoteIdent(s string) string {
	return fmt.Sprintf(`%q`, strings.ReplaceAll(s, `"`, `""`))
}
