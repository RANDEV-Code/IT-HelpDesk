package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"randesk/backend/internal/platform/httpx/apierr"
	"randesk/backend/internal/platform/httpx/respond"
	"randesk/backend/internal/platform/logger"
)

// testRouter merakit rantai middleware yang sama dengan httpx.NewRouter plus
// beberapa rute contoh untuk menguji envelope, panic, dan timeout.
func testRouter(buf *bytes.Buffer, jsonLimit, multipartLimit int64, timeout time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	l := logger.New(buf, "info")

	r := gin.New()
	r.Use(
		RequestID(),
		Recovery(l),
		BodyLimit(jsonLimit, multipartLimit),
		Logging(l),
		Timeout(timeout),
	)

	r.GET("/api/v1/ok", func(c *gin.Context) {
		respond.JSON(c, http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/v1/echo", func(c *gin.Context) {
		var p struct {
			Title string `json:"title"`
		}
		if ae := apierr.DecodeJSON(c.Request, &p); ae != nil {
			respond.Error(c, ae)
			return
		}
		respond.JSON(c, http.StatusOK, gin.H{"title": p.Title})
	})

	r.GET("/api/v1/boom", func(c *gin.Context) {
		panic("ledakan internal rahasia")
	})

	// kooperatif: menunggu deadline lalu mengembalikan ctx.Err() -> 500
	r.GET("/api/v1/slow", func(c *gin.Context) {
		<-c.Request.Context().Done()
		respond.Error(c, c.Request.Context().Err())
	})

	// backstop: menunggu deadline tanpa menulis apa pun -> middleware menulis 500
	r.GET("/api/v1/hang", func(c *gin.Context) {
		<-c.Request.Context().Done()
	})

	return r
}

func TestRequestID_HeaderMatchesEnvelopeMeta(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, time.Second)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/ok", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ingin 200", rec.Code)
	}
	headerID := rec.Header().Get("X-Request-ID")
	if len(headerID) != 36 {
		t.Fatalf("X-Request-ID bukan UUID (len=%d): %q", len(headerID), headerID)
	}

	var env struct {
		Meta struct {
			RequestID string `json:"request_id"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("body bukan JSON: %v (%s)", err, rec.Body.String())
	}
	if env.Meta.RequestID != headerID {
		t.Errorf("meta.request_id = %q, header = %q; harus sama", env.Meta.RequestID, headerID)
	}
}

func TestRequestID_UnikAntarRequest(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, time.Second)

	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, httptest.NewRequest(http.MethodGet, "/api/v1/ok", nil))
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/api/v1/ok", nil))

	id1 := rec1.Header().Get("X-Request-ID")
	id2 := rec2.Header().Get("X-Request-ID")
	if id1 == id2 {
		t.Errorf("request_id harus unik antar request, keduanya %q", id1)
	}
}

func TestBodyLimit_OversizedRejected413(t *testing.T) {
	var buf bytes.Buffer
	// batas kecil agar tes cepat; perilaku identik dengan 64 KiB produksi
	r := testRouter(&buf, 64, DefaultMultipartBodyLimit, time.Second)

	big := `{"title":"` + strings.Repeat("a", 200) + `"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(big))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, ingin 413", rec.Code)
	}
	var env apierr.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	if env.Error.Code != apierr.CodePayloadTooLarge {
		t.Errorf("code = %q, ingin PAYLOAD_TOO_LARGE", env.Error.Code)
	}
	if env.Meta.RequestID == "" {
		t.Errorf("error envelope harus memuat meta.request_id")
	}
}

func TestDecodeJSON_UnknownFieldRejected400(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, time.Second)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(`{"title":"ok","asing":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, ingin 400", rec.Code)
	}
	var env apierr.ErrorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if env.Error.Code != apierr.CodeMalformedRequest {
		t.Errorf("code = %q, ingin MALFORMED_REQUEST", env.Error.Code)
	}
}

func TestRecovery_PanicBecomes500WithoutStack(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, time.Second)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/boom", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, ingin 500", rec.Code)
	}
	body := rec.Body.String()
	var env apierr.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	if env.Error.Code != apierr.CodeInternalError {
		t.Errorf("code = %q, ingin INTERNAL_ERROR", env.Error.Code)
	}
	// body klien TIDAK boleh memuat pesan panic atau stack
	for _, leak := range []string{"ledakan internal rahasia", "goroutine", "panic("} {
		if strings.Contains(body, leak) {
			t.Errorf("body membocorkan %q: %s", leak, body)
		}
	}
	// log server HARUS memuat stack untuk diagnostik
	logs := buf.String()
	if !strings.Contains(logs, "panic recovered") {
		t.Errorf("log harus memuat 'panic recovered': %s", logs)
	}
	if !strings.Contains(logs, "goroutine") {
		t.Errorf("log harus memuat stack (goroutine): %s", logs)
	}
}

func TestLogging_RequiredFieldsAndNoBody(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, time.Second)

	secretBody := "RAHASIA_BODY_XYZ"
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(`{"title":"`+secretBody+`"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ingin 200", rec.Code)
	}

	logs := buf.String()
	if !strings.Contains(logs, "http request") {
		t.Fatalf("baris log request tidak ditemukan: %s", logs)
	}

	var entry map[string]any
	// ambil baris JSON log (satu baris per request)
	for _, line := range strings.Split(strings.TrimSpace(logs), "\n") {
		if strings.Contains(line, "http request") {
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				t.Fatalf("log bukan JSON valid: %v (%s)", err, line)
			}
			break
		}
	}

	for _, field := range []string{"request_id", "method", "route", "status", "duration_ms", "time", "level"} {
		if _, ok := entry[field]; !ok {
			t.Errorf("log kehilangan field wajib %q: %v", field, entry)
		}
	}
	// route harus template, bukan path mentah
	if entry["route"] != "/api/v1/echo" {
		t.Errorf("route = %v, ingin template /api/v1/echo", entry["route"])
	}
	if entry["method"] != http.MethodPost {
		t.Errorf("method = %v, ingin POST", entry["method"])
	}
	// log TIDAK boleh memuat nilai body
	if strings.Contains(logs, secretBody) {
		t.Errorf("log membocorkan nilai body %q: %s", secretBody, logs)
	}
}

func TestTimeout_CooperativeHandler500(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, 30*time.Millisecond)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/slow", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, ingin 500 (handler kooperatif timeout)", rec.Code)
	}
	var env apierr.ErrorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if env.Error.Code != apierr.CodeInternalError {
		t.Errorf("code = %q, ingin INTERNAL_ERROR", env.Error.Code)
	}
}

func TestTimeout_BackstopWrites500(t *testing.T) {
	var buf bytes.Buffer
	r := testRouter(&buf, DefaultJSONBodyLimit, DefaultMultipartBodyLimit, 30*time.Millisecond)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/hang", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, ingin 500 (backstop timeout)", rec.Code)
	}
}

func TestNewUUIDv4_Format(t *testing.T) {
	id := newUUIDv4()
	if len(id) != 36 {
		t.Fatalf("panjang UUID = %d, ingin 36: %q", len(id), id)
	}
	if id[14] != '4' {
		t.Errorf("versi UUID harus '4' di posisi 14: %q", id)
	}
	if !strings.ContainsAny(string(id[19]), "89ab") {
		t.Errorf("variant RFC4122 harus 8/9/a/b di posisi 19: %q", id)
	}
}
