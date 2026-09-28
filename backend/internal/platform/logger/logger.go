// Package logger menyediakan logger JSON terstruktur untuk platform RANDesk.
//
// Aturan SECURITY §8 / ARCHITECTURE §10: log TIDAK boleh memuat body request,
// cookie, token, atau nilai header sensitif. Pemanggil hanya boleh meneruskan
// field yang diizinkan (request_id, route template, method, status,
// duration_ms, actor_id, dan pesan error internal).
package logger

import (
	"io"
	"log/slog"
)

// New membuat logger slog JSON yang menulis ke w.
// level: "debug", "info", "warn", "error" (default info).
func New(w io.Writer, level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: lvl}))
}
