// Command api menjalankan HTTP API RANDesk.
//
// Siklus hidup (TASK-002):
//  1. Muat .env opsional → validasi konfigurasi (gagal start bila tidak valid).
//  2. Buat pool DB lazy + router health.
//  3. Serve HTTP; pada SIGINT/SIGTERM lakukan graceful shutdown maks 15 detik.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"randesk/backend/internal/platform/config"
	"randesk/backend/internal/platform/db"
	"randesk/backend/internal/platform/httpx"
	"randesk/backend/internal/platform/httpx/middleware"
)

const shutdownTimeout = 15 * time.Second

func main() {
	if err := run(); err != nil {
		// konfigurasi gagal → pesan ke stderr, exit code 1
		fmt.Fprintf(os.Stderr, "startup gagal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// .env opsional di root repo (satu level di atas backend/ saat go run ./cmd/api)
	for _, p := range []string{".env", filepath.Join("..", ".env")} {
		if err := config.LoadDotEnv(p); err != nil {
			return fmt.Errorf("membaca %s: %w", p, err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := newLogger(cfg)
	slog.SetDefault(logger)

	pool, err := db.Open(cfg.DatabaseURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	if err != nil {
		return err
	}
	defer pool.Close()

	ready := func(ctx context.Context) error {
		if err := db.Ping(ctx, pool); err != nil {
			return err
		}
		return storageWritable(cfg.StorageRoot)
	}

	router := httpx.NewRouter(ready, 3*time.Second,
		httpx.WithLogger(logger),
		httpx.WithHandlerTimeout(middleware.DefaultHandlerTimeout),
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// context batal saat SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api starting", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal diterima")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown melebihi batas; keluar paksa", "error", err)
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.Info("api berhenti bersih")
	return nil
}

// storageWritable memastikan STORAGE_ROOT ada dan dapat ditulis (readiness
// lampiran). File probe dihapus kembali; kegagalan tidak membocorkan path ke
// respons HTTP — hanya ke log.
func storageWritable(root string) error {
	if err := os.MkdirAll(root, 0o750); err != nil {
		return fmt.Errorf("storage root tidak siap: %w", err)
	}
	probe := filepath.Join(root, ".write-probe")
	f, err := os.Create(probe)
	if err != nil {
		return fmt.Errorf("storage root tidak writable: %w", err)
	}
	_ = f.Close()
	_ = os.Remove(probe)
	return nil
}

func newLogger(cfg config.Config) *slog.Logger {
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	// development: teks mudah dibaca; lingkungan lain: JSON terstruktur
	var h slog.Handler
	opts := &slog.HandlerOptions{Level: level}
	if cfg.AppEnv == "development" {
		h = slog.NewTextHandler(os.Stdout, opts)
	} else {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(h)
}
