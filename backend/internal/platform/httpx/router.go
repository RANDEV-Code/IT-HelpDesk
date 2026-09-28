// Package httpx merakit router Gin dan endpoint infrastruktur.
// Middleware platform lengkap (request ID, log, timeout, error envelope)
// menyusul di TASK-004.
package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ReadyCheck melaporkan apakah dependensi runtime siap melayani traffic.
// Error yang dikembalikan HANYA untuk log server — tidak boleh dikirim ke body
// respons karena bisa membocorkan detail infrastruktur (ARCHITECTURE §10).
type ReadyCheck func(ctx context.Context) error

// NewRouter membuat router dengan endpoint health. readyTimeout membatasi
// pemeriksaan readiness (default 3 detik bila nol).
func NewRouter(ready ReadyCheck, readyTimeout time.Duration) *gin.Engine {
	if readyTimeout <= 0 {
		readyTimeout = 3 * time.Second
	}

	r := gin.New()
	r.Use(gin.Recovery())
	// Tidak memakai gin.Logger default di production tanpa request ID;
	// logging request terstruktur dipasang di TASK-004.

	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), readyTimeout)
		defer cancel()

		if err := ready(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	return r
}
