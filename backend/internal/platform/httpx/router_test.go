package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthLive(t *testing.T) {
	r := NewRouter(func(context.Context) error { return nil }, time.Second)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health/live status = %d, ingin 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body bukan JSON valid: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, ingin ok", body["status"])
	}
}

func TestHealthReady_OK(t *testing.T) {
	r := NewRouter(func(context.Context) error { return nil }, time.Second)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health/ready status = %d, ingin 200", rec.Code)
	}
}

func TestHealthReady_Unavailable(t *testing.T) {
	// Error internal memuat "connection string" palsu untuk memastikan body
	// respons TIDAK membocorkan detail infrastruktur (ARCHITECTURE §10).
	r := NewRouter(func(context.Context) error {
		return errors.New("dial postgres://user:rahasia@db:5432/randesk: timeout")
	}, time.Second)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready status = %d, ingin 503", rec.Code)
	}
	body := rec.Body.String()
	for _, leak := range []string{"rahasia", "postgres://", "timeout", "dial"} {
		if strings.Contains(body, leak) {
			t.Errorf("body membocorkan %q: %s", leak, body)
		}
	}
}

func TestHealthReady_TimeoutDihormati(t *testing.T) {
	called := make(chan struct{})
	r := NewRouter(func(ctx context.Context) error {
		close(called)
		<-ctx.Done() // simulasi dependensi lambat
		return ctx.Err()
	}, 50*time.Millisecond)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(rec, req)

	<-called
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("ready lambat status = %d, ingin 503", rec.Code)
	}
}
