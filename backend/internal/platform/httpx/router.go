// Package httpx merakit router Gin, middleware platform, dan endpoint
// infrastruktur. Middleware lengkap (request ID, log, recovery, body limit,
// timeout, error envelope) dipasang di sini (TASK-004).
package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"randesk/backend/internal/platform/httpx/middleware"
)

// ReadyCheck melaporkan apakah dependensi runtime siap melayani traffic.
// Error yang dikembalikan HANYA untuk log server — tidak boleh dikirim ke body
// respons karena bisa membocorkan detail infrastruktur (ARCHITECTURE §10).
type ReadyCheck func(ctx context.Context) error

// Options mengonfigurasi middleware platform pada router.
type Options struct {
	Logger         *slog.Logger
	JSONBodyLimit  int64
	MultipartLimit int64
	HandlerTimeout time.Duration
}

// Option adalah functional option untuk NewRouter.
type Option func(*Options)

// WithLogger mengeset logger terstruktur untuk middleware logging/recovery.
func WithLogger(l *slog.Logger) Option {
	return func(o *Options) { o.Logger = l }
}

// WithBodyLimits mengeset batas body JSON dan multipart (byte).
func WithBodyLimits(jsonLimit, multipartLimit int64) Option {
	return func(o *Options) {
		o.JSONBodyLimit = jsonLimit
		o.MultipartLimit = multipartLimit
	}
}

// WithHandlerTimeout mengeset timeout handler JSON (context turunan).
func WithHandlerTimeout(d time.Duration) Option {
	return func(o *Options) { o.HandlerTimeout = d }
}

// NewRouter membuat router dengan middleware platform + endpoint health.
// readyTimeout membatasi pemeriksaan readiness (default 3 detik bila nol).
// Urutan middleware mengikuti ARCHITECTURE §5.1.
func NewRouter(ready ReadyCheck, readyTimeout time.Duration, opts ...Option) *gin.Engine {
	if readyTimeout <= 0 {
		readyTimeout = 3 * time.Second
	}

	o := Options{
		Logger:         slog.Default(),
		JSONBodyLimit:  middleware.DefaultJSONBodyLimit,
		MultipartLimit: middleware.DefaultMultipartBodyLimit,
		HandlerTimeout: middleware.DefaultHandlerTimeout,
	}
	for _, fn := range opts {
		fn(&o)
	}
	if o.Logger == nil {
		o.Logger = slog.Default()
	}

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recovery(o.Logger),
		middleware.BodyLimit(o.JSONBodyLimit, o.MultipartLimit),
		middleware.Logging(o.Logger),
		middleware.Timeout(o.HandlerTimeout),
	)

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
