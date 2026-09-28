package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// envValid menyetel environment minimum yang valid; tes memakai t.Setenv agar
// otomatis dipulihkan.
func envValid(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_ORIGIN", "http://localhost:5173")
	t.Setenv("HTTP_ADDR", "127.0.0.1:8080")
	t.Setenv("DATABASE_URL", "postgres://randesk:secret@localhost:5432/randesk_dev?sslmode=disable")
	t.Setenv("STORAGE_ROOT", "./storage")
	t.Setenv("LOG_LEVEL", "debug")
}

func TestLoad_ValidMinimal(t *testing.T) {
	envValid(t)
	// pastikan nilai lama tidak bocor antar-tes
	for _, k := range []string{"SESSION_TTL_HOURS", "TRUSTED_PROXIES", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS"} {
		t.Setenv(k, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, ingin nil", err)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("AppEnv = %q, ingin development", cfg.AppEnv)
	}
	if cfg.AppOrigin != "http://localhost:5173" {
		t.Errorf("AppOrigin = %q, ingin http://localhost:5173", cfg.AppOrigin)
	}
	if cfg.SessionTTLHour != 8 {
		t.Errorf("SessionTTLHour = %d, ingin default 8", cfg.SessionTTLHour)
	}
	if cfg.DBMaxOpenConns != 10 || cfg.DBMaxIdleConns != 5 {
		t.Errorf("pool default = %d/%d, ingin 10/5", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	}
	if len(cfg.TrustedProxies) != 0 {
		t.Errorf("TrustedProxies = %v, ingin kosong", cfg.TrustedProxies)
	}
}

func TestLoad_Failures(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(t *testing.T)
		wantErr string // substring yang harus muncul di error
	}{
		{
			name:    "APP_ENV hilang",
			mutate:  func(t *testing.T) { t.Setenv("APP_ENV", "") },
			wantErr: "APP_ENV",
		},
		{
			name:    "APP_ENV tidak dikenal",
			mutate:  func(t *testing.T) { t.Setenv("APP_ENV", "banana") },
			wantErr: "APP_ENV",
		},
		{
			name:    "APP_ORIGIN wildcard",
			mutate:  func(t *testing.T) { t.Setenv("APP_ORIGIN", "https://*.example.com") },
			wantErr: "wildcard",
		},
		{
			name:    "APP_ORIGIN tanpa scheme",
			mutate:  func(t *testing.T) { t.Setenv("APP_ORIGIN", "localhost:5173") },
			wantErr: "APP_ORIGIN",
		},
		{
			name:    "HTTP_ADDR hilang",
			mutate:  func(t *testing.T) { t.Setenv("HTTP_ADDR", "") },
			wantErr: "HTTP_ADDR",
		},
		{
			name:    "DATABASE_URL hilang",
			mutate:  func(t *testing.T) { t.Setenv("DATABASE_URL", "") },
			wantErr: "DATABASE_URL",
		},
		{
			name:    "DATABASE_URL bukan postgres",
			mutate:  func(t *testing.T) { t.Setenv("DATABASE_URL", "mysql://user:pass@localhost/db") },
			wantErr: "postgres",
		},
		{
			name:    "STORAGE_ROOT hilang",
			mutate:  func(t *testing.T) { t.Setenv("STORAGE_ROOT", "") },
			wantErr: "STORAGE_ROOT",
		},
		{
			name:    "SESSION_TTL_HOURS bukan angka",
			mutate:  func(t *testing.T) { t.Setenv("SESSION_TTL_HOURS", "delapan") },
			wantErr: "SESSION_TTL_HOURS",
		},
		{
			name:    "SESSION_TTL_HOURS nol",
			mutate:  func(t *testing.T) { t.Setenv("SESSION_TTL_HOURS", "0") },
			wantErr: "SESSION_TTL_HOURS",
		},
		{
			name:    "LOG_LEVEL tidak dikenal",
			mutate:  func(t *testing.T) { t.Setenv("LOG_LEVEL", "verbose") },
			wantErr: "LOG_LEVEL",
		},
		{
			name: "production origin http",
			mutate: func(t *testing.T) {
				t.Setenv("APP_ENV", "production")
				t.Setenv("APP_ORIGIN", "http://helpdesk.example.com")
				t.Setenv("LOG_LEVEL", "info")
			},
			wantErr: "https",
		},
		{
			name: "production log debug",
			mutate: func(t *testing.T) {
				t.Setenv("APP_ENV", "production")
				t.Setenv("APP_ORIGIN", "https://helpdesk.example.com")
				t.Setenv("LOG_LEVEL", "debug")
			},
			wantErr: "LOG_LEVEL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envValid(t)
			tt.mutate(t)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() = nil error, ingin error memuat %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q tidak memuat %q", err.Error(), tt.wantErr)
			}
			// secret tidak boleh bocor ke pesan error
			if strings.Contains(err.Error(), "secret") {
				t.Errorf("pesan error membocorkan isi DATABASE_URL: %q", err.Error())
			}
		})
	}
}

func TestLoad_TrustedProxiesList(t *testing.T) {
	envValid(t)
	t.Setenv("TRUSTED_PROXIES", "10.0.0.2, 10.0.0.3")
	t.Setenv("SESSION_TTL_HOURS", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_MAX_IDLE_CONNS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := []string{"10.0.0.2", "10.0.0.3"}
	if strings.Join(cfg.TrustedProxies, ",") != strings.Join(want, ",") {
		t.Errorf("TrustedProxies = %v, ingin %v", cfg.TrustedProxies, want)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "# komentar\nAPP_ENV=development\n\nAPP_ORIGIN=\"http://localhost:5173\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// Pastikan variabel belum disetel di environment proses. t.Setenv("", ...) saja
	// tidak cukup karena LookupEnv tetap menganggap key "ada"; Unsetenv membuat
	// LoadDotEnv benar-benar mengisi. Cleanup t.Setenv tetap memulihkan kondisi awal.
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_ORIGIN", "")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_ORIGIN")
	if err := LoadDotEnv(path); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != "development" {
		t.Errorf("APP_ENV = %q, ingin development", got)
	}
	if got := os.Getenv("APP_ORIGIN"); got != "http://localhost:5173" {
		t.Errorf("APP_ORIGIN = %q (kutip harus di-trim)", got)
	}

	// file tidak ada → bukan error
	if err := LoadDotEnv(filepath.Join(dir, "missing.env")); err != nil {
		t.Errorf("LoadDotEnv(file hilang) error = %v, ingin nil", err)
	}

	// baris rusak → error menyebut nomor baris
	bad := filepath.Join(dir, "bad.env")
	if err := os.WriteFile(bad, []byte("A=1\nRUSAK\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := LoadDotEnv(bad); err == nil || !strings.Contains(err.Error(), "2") {
		t.Errorf("LoadDotEnv(bad) error = %v, ingin error baris 2", err)
	}
}
