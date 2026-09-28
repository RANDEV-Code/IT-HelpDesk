// Package config memuat dan memvalidasi konfigurasi aplikasi dari environment
// (OPERATIONS.md §2). Tidak ada secret yang boleh dilog dari package ini.
package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config adalah konfigurasi runtime API yang sudah tervalidasi.
type Config struct {
	AppEnv         string   // development | test | staging | production
	AppOrigin      string   // origin browser sah (scheme+host+port)
	HTTPAddr       string   // bind address API
	DatabaseURL    string   // koneksi PostgreSQL (secret — jangan dilog)
	StorageRoot    string   // direktori lampiran
	SessionTTLHour int      // masa aktif sesi (jam), default 8
	LogLevel       string   // debug | info | warn | error
	TrustedProxies []string // proxy yang dipercaya (boleh kosong)
	DBMaxOpenConns int      // default 10
	DBMaxIdleConns int      // default 5
}

// LoadDotEnv memuat file .env sederhana (KEY=VALUE per baris) ke environment
// proses. Nilai yang sudah ada di environment TIDAK ditimpa. File yang tidak
// ada bukan error (dipakai opsional di development).
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("%s:%d: baris bukan KEY=VALUE", path, lineNo)
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(strings.Trim(val, `"`))
		if key == "" {
			return fmt.Errorf("%s:%d: nama variabel kosong", path, lineNo)
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return sc.Err()
}

// Load membaca environment dan memvalidasinya. Error yang dikembalikan bersifat
// deskriptif tetapi tidak pernah memuat nilai secret (mis. DATABASE_URL).
func Load() (Config, error) {
	var cfg Config

	cfg.AppEnv = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	switch cfg.AppEnv {
	case "development", "test", "staging", "production":
	default:
		return Config{}, fmt.Errorf("APP_ENV wajib salah satu dari development|test|staging|production")
	}

	origin := strings.TrimSpace(os.Getenv("APP_ORIGIN"))
	if origin == "" {
		return Config{}, fmt.Errorf("APP_ORIGIN wajib diisi (scheme + host + port)")
	}
	if strings.Contains(origin, "*") {
		return Config{}, fmt.Errorf("APP_ORIGIN tidak boleh mengandung wildcard")
	}
	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Config{}, fmt.Errorf("APP_ORIGIN harus URL absolut dengan scheme http/https")
	}
	if strings.TrimSuffix(u.Path, "/") != "" {
		return Config{}, fmt.Errorf("APP_ORIGIN tidak boleh memuat path: %q", origin)
	}
	cfg.AppOrigin = u.Scheme + "://" + u.Host

	if cfg.HTTPAddr = strings.TrimSpace(os.Getenv("HTTP_ADDR")); cfg.HTTPAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR wajib diisi (contoh 127.0.0.1:8080)")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(dbURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL wajib diisi (nilai tidak ditampilkan)")
	}
	if !strings.HasPrefix(dbURL, "postgres://") && !strings.HasPrefix(dbURL, "postgresql://") {
		return Config{}, fmt.Errorf("DATABASE_URL harus berawalan postgres:// atau postgresql:// (nilai tidak ditampilkan)")
	}
	cfg.DatabaseURL = dbURL

	if cfg.StorageRoot = strings.TrimSpace(os.Getenv("STORAGE_ROOT")); cfg.StorageRoot == "" {
		return Config{}, fmt.Errorf("STORAGE_ROOT wajib diisi")
	}

	if cfg.SessionTTLHour, err = envInt("SESSION_TTL_HOURS", 8); err != nil {
		return Config{}, err
	}
	if cfg.SessionTTLHour <= 0 {
		return Config{}, fmt.Errorf("SESSION_TTL_HOURS harus > 0")
	}

	cfg.LogLevel = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	switch cfg.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return Config{}, fmt.Errorf("LOG_LEVEL wajib salah satu dari debug|info|warn|error")
	}

	if raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES")); raw != "" {
		for _, p := range strings.Split(raw, ",") {
			if p = strings.TrimSpace(p); p != "" {
				cfg.TrustedProxies = append(cfg.TrustedProxies, p)
			}
		}
	}

	if cfg.DBMaxOpenConns, err = envInt("DB_MAX_OPEN_CONNS", 10); err != nil {
		return Config{}, err
	}
	if cfg.DBMaxIdleConns, err = envInt("DB_MAX_IDLE_CONNS", 5); err != nil {
		return Config{}, err
	}
	if cfg.DBMaxOpenConns <= 0 || cfg.DBMaxIdleConns < 0 {
		return Config{}, fmt.Errorf("DB_MAX_OPEN_CONNS harus > 0 dan DB_MAX_IDLE_CONNS >= 0")
	}
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		cfg.DBMaxIdleConns = cfg.DBMaxOpenConns
	}

	if err := cfg.validateProductionSafety(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// validateProductionSafety menolak start di production bila konfigurasi tidak
// aman (OPERATIONS §2). Pesan error tidak memuat nilai secret.
func (c Config) validateProductionSafety() error {
	if c.AppEnv != "production" {
		return nil
	}
	u, _ := url.Parse(c.AppOrigin)
	if u == nil || u.Scheme != "https" {
		return fmt.Errorf("APP_ENV=production mengharuskan APP_ORIGIN berscheme https")
	}
	if c.LogLevel == "debug" {
		return fmt.Errorf("APP_ENV=production tidak mengizinkan LOG_LEVEL=debug")
	}
	return nil
}

func envInt(name string, def int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s harus berupa angka, didapat nilai non-numerik", name)
	}
	return v, nil
}
