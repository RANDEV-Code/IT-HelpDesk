// Package db menyediakan pool koneksi PostgreSQL (pgx melalui database/sql).
package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // driver pgx untuk database/sql
)

// Open membuat pool sesuai konfigurasi. Koneksi bersifat lazy: Open tidak
// gagal saat database mati — kegagalan baru terlihat saat Ping/query, sehingga
// API tetap bisa start dan /health/ready melaporkan 503 (ARCHITECTURE §10).
func Open(databaseURL string, maxOpen, maxIdle int) (*sql.DB, error) {
	pool, err := sql.Open("pgx", databaseURL)
	if err != nil {
		// jangan membungkus err: bisa memuat connection string
		return nil, fmt.Errorf("db: gagal membuka pool koneksi")
	}
	pool.SetMaxOpenConns(maxOpen)
	pool.SetMaxIdleConns(maxIdle)
	pool.SetConnMaxIdleTime(5 * time.Minute)
	return pool, nil
}

// Ping memeriksa konektivitas dengan batas waktu (default query 3 detik,
// ARCHITECTURE §5). Error yang dikembalikan tidak memuat connection string.
func Ping(ctx context.Context, pool *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.PingContext(ctx); err != nil {
		return fmt.Errorf("db: ping gagal: %w", err)
	}
	return nil
}
