package testdb

import (
	"context"
	"testing"
	"time"
)

// TestIntegrationMigrateAndCRUD adalah contoh satu tes DB integrasi lewat helper
// testdb (TASK-007 acceptance): menerapkan migrasi, menulis, membaca kembali,
// dan menegakkan CHECK constraint pada PostgreSQL nyata. Di-skip bila
// TEST_DATABASE_URL tidak diset (lihat Connect).
func TestIntegrationMigrateAndCRUD(t *testing.T) {
	db := Connect(t)
	TruncateAll(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Baseline harus membuat tabel departments.
	var ada bool
	if err := db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables
		  WHERE table_schema='public' AND table_name='departments')`).Scan(&ada); err != nil {
		t.Fatalf("cek tabel departments: %v", err)
	}
	if !ada {
		t.Fatalf("tabel departments tidak ada setelah migrasi")
	}

	// Insert valid + baca kembali.
	var id string
	if err := db.QueryRowContext(ctx,
		`INSERT INTO departments (id, name) VALUES (gen_random_uuid(), 'Fondasi CI') RETURNING id`).Scan(&id); err != nil {
		t.Fatalf("insert department: %v", err)
	}
	if id == "" {
		t.Fatalf("id kosong setelah insert")
	}

	var jumlah int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM departments`).Scan(&jumlah); err != nil {
		t.Fatalf("count departments: %v", err)
	}
	if jumlah != 1 {
		t.Fatalf("jumlah departments = %d, ingin 1", jumlah)
	}

	// CHECK constraint SCHEMA ditegakkan: name 1 karakter wajib ditolak.
	_, err := db.ExecContext(ctx,
		`INSERT INTO departments (id, name) VALUES (gen_random_uuid(), 'x')`)
	if err == nil {
		t.Fatalf("CHECK char_length(name) tidak ditegakkan")
	}
}
