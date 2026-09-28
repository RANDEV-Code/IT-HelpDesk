// Package middleware menyediakan middleware platform HTTP RANDesk
// (ARCHITECTURE §5, §10; SECURITY §7, §8; API_SPEC §1).
//
// Urutan pemasangan (ARCHITECTURE §5.1): RequestID -> Recovery -> BodyLimit ->
// Logging -> Timeout. Middleware ini tidak pernah mencatat body, cookie, token,
// atau header sensitif ke log.
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"randesk/backend/internal/platform/httpx/apierr"
	"randesk/backend/internal/platform/httpx/respond"
)

// Batas dan timeout default (SECURITY §7; ARCHITECTURE §5).
const (
	DefaultJSONBodyLimit      int64 = 64 << 10 // 64 KiB
	DefaultMultipartBodyLimit int64 = 6 << 20  // 6 MiB
	DefaultHandlerTimeout           = 10 * time.Second
	DefaultQueryTimeout             = 3 * time.Second
)

// RequestID menghasilkan UUID v4 per request, menyetel header X-Request-ID,
// dan menyimpannya di gin.Context + request context (API_SPEC §1: request_id
// dibuat server).
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := newUUIDv4()
		c.Set(respond.RequestIDKey, id)
		c.Header("X-Request-ID", id)
		c.Request = c.Request.WithContext(respond.ContextWithRequestID(c.Request.Context(), id))
		c.Next()
	}
}

// BodyLimit membatasi ukuran body: JSON 64 KiB, multipart 6 MiB (SECURITY §7).
// Content-Length yang melebihi batas langsung ditolak 413; body streaming
// dibungkus http.MaxBytesReader sehingga DecodeJSON di handler memetakan
// kelebihan menjadi 413.
func BodyLimit(jsonLimit, multipartLimit int64) gin.HandlerFunc {
	if jsonLimit <= 0 {
		jsonLimit = DefaultJSONBodyLimit
	}
	if multipartLimit <= 0 {
		multipartLimit = DefaultMultipartBodyLimit
	}
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}

		limit := jsonLimit
		if strings.HasPrefix(c.ContentType(), "multipart/") {
			limit = multipartLimit
		}

		if c.Request.ContentLength > limit {
			c.Abort()
			respond.Error(c, apierr.PayloadTooLarge())
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		c.Next()
	}
}

// Logging mencatat satu baris JSON per request: request_id, method, route
// template (bukan path mentah), status, duration_ms, dan actor_id bila ada
// (ARCHITECTURE §10). Tidak ada body/cookie/header yang dicatat.
func Logging(l *slog.Logger) gin.HandlerFunc {
	if l == nil {
		l = slog.Default()
	}
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "(unmatched)"
		}
		attrs := []any{
			"request_id", respond.RequestID(c),
			"method", c.Request.Method,
			"route", route,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if v, ok := c.Get(respond.ActorIDKey); ok {
			attrs = append(attrs, "actor_id", v)
		}
		l.InfoContext(c.Request.Context(), "http request", attrs...)
	}
}

// Recovery menangkap panic, mencatat stack HANYA ke log server, lalu merespons
// 500 INTERNAL_ERROR generik (tanpa stack/detail ke klien — ARCHITECTURE §10).
func Recovery(l *slog.Logger) gin.HandlerFunc {
	if l == nil {
		l = slog.Default()
	}
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l.ErrorContext(c.Request.Context(), "panic recovered",
					"request_id", respond.RequestID(c),
					"panic", r,
					"stack", string(debug.Stack()),
				)
				if c.Writer.Written() {
					c.Abort()
					return
				}
				c.Abort()
				respond.Error(c, apierr.Internal())
			}
		}()
		c.Next()
	}
}

// Timeout menurunkan context request dengan deadline handler (default 10 detik).
// Handler/service/repository wajib menghormati ctx (query memakai timeout lebih
// pendek, DefaultQueryTimeout). Bila handler selesai tanpa menulis respons dan
// deadline terlampaui, middleware menulis 500 terkontrol.
//
// Pendekatan ini sengaja kooperatif (context turunan), bukan membatalkan handler
// lewat goroutine — sesuai ARCHITECTURE §5: hindari goroutine tanpa pembatalan.
func Timeout(d time.Duration) gin.HandlerFunc {
	if d <= 0 {
		d = DefaultHandlerTimeout
	}
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if !c.Writer.Written() && ctx.Err() == context.DeadlineExceeded {
			respond.Error(c, apierr.Internal().WithInternal(ctx.Err()))
		}
	}
}

// newUUIDv4 menghasilkan UUID versi 4 (acak) tanpa dependensi eksternal.
// Bila crypto/rand gagal (sangat jarang), fallback tidak dipakai; fungsi
// mengembalikan string yang tetap unik berbasis waktu+nanodetik.
func newUUIDv4() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// fallback deterministik-unik berbasis nanodetik (bukan rahasia)
		ns := time.Now().UnixNano()
		for i := 0; i < 8; i++ {
			b[i] = byte(ns >> (8 * i))
		}
		for i := 8; i < 16; i++ {
			b[i] = byte(ns >> (8 * (i - 8)))
		}
	}
	// set version (4) dan variant (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	var buf [36]byte
	hex.Encode(buf[0:8], b[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], b[10:16])
	return string(buf[:])
}
