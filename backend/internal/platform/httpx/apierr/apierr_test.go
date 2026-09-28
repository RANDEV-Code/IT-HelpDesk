package apierr

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestErrorMappingTableDriven memastikan setiap kelas error API_SPEC §3
// dipetakan ke status + code yang tepat.
func TestErrorMappingTableDriven(t *testing.T) {
	cases := []struct {
		name       string
		err        *APIError
		wantStatus int
		wantCode   string
	}{
		{"malformed", MalformedRequest(), http.StatusBadRequest, CodeMalformedRequest},
		{"unauthenticated", Unauthenticated(), http.StatusUnauthorized, CodeUnauthenticated},
		{"forbidden", Forbidden(), http.StatusForbidden, CodeForbidden},
		{"csrf", CSRFFailed(), http.StatusForbidden, CodeCSRFCailed},
		{"password_change", PasswordChangeRequired(), http.StatusForbidden, CodePasswordChangeReq},
		{"not_found", NotFound(), http.StatusNotFound, CodeNotFound},
		{"version_conflict", VersionConflict(0), http.StatusConflict, CodeVersionConflict},
		{"already_assigned", AlreadyAssigned(), http.StatusConflict, CodeAlreadyAssigned},
		{"invalid_state", InvalidState(), http.StatusConflict, CodeInvalidState},
		{"active_tickets", ActiveTicketsExist(), http.StatusConflict, CodeActiveTicketsExist},
		{"last_admin", LastAdmin(), http.StatusConflict, CodeLastAdmin},
		{"email_exists", EmailExists(), http.StatusConflict, CodeEmailExists},
		{"name_exists", NameExists(), http.StatusConflict, CodeNameExists},
		{"payload_too_large", PayloadTooLarge(), http.StatusRequestEntityTooLarge, CodePayloadTooLarge},
		{"unsupported_media", UnsupportedMediaType(), http.StatusUnsupportedMediaType, CodeUnsupportedMediaType},
		{"validation", ValidationError(map[string]string{"title": "x"}), http.StatusUnprocessableEntity, CodeValidationError},
		{"master_inactive", MasterInactive(), http.StatusUnprocessableEntity, CodeMasterInactive},
		{"no_change", NoChange(), http.StatusUnprocessableEntity, CodeNoChange},
		{"attachment_limit", AttachmentLimit(), http.StatusUnprocessableEntity, CodeAttachmentLimit},
		{"rate_limited", RateLimited(0), http.StatusTooManyRequests, CodeRateLimited},
		{"internal", Internal(), http.StatusInternalServerError, CodeInternalError},
		{"service_unavailable", ServiceUnavailable(), http.StatusServiceUnavailable, CodeServiceUnavailable},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Status() != tc.wantStatus {
				t.Errorf("status = %d, ingin %d", tc.err.Status(), tc.wantStatus)
			}
			if tc.err.Code() != tc.wantCode {
				t.Errorf("code = %q, ingin %q", tc.err.Code(), tc.wantCode)
			}
			if tc.err.Error() == "" {
				t.Errorf("pesan default kosong untuk code %q", tc.wantCode)
			}
		})
	}
}

func TestValidationErrorFields(t *testing.T) {
	e := ValidationError(map[string]string{"title": "Judul harus 5-150 karakter."})
	if e.Fields()["title"] == "" {
		t.Fatalf("fields.title tidak terbawa: %#v", e.Fields())
	}
}

func TestVersionConflictCarriesCurrentVersion(t *testing.T) {
	e := VersionConflict(7)
	cv, ok := e.CurrentVersion()
	if !ok || cv != 7 {
		t.Fatalf("current_version = %d,%v; ingin 7,true", cv, ok)
	}
	if _, ok := VersionConflict(0).CurrentVersion(); ok {
		t.Fatalf("current_version harus tidak diset saat 0")
	}
}

func TestFromMapping(t *testing.T) {
	if got := From(nil); got != nil {
		t.Errorf("From(nil) = %v, ingin nil", got)
	}

	// *APIError dikembalikan apa adanya
	orig := NotFound()
	if got := From(orig); got != orig {
		t.Errorf("From(*APIError) harus mengembalikan instance yang sama")
	}

	// error biasa -> 500 INTERNAL_ERROR generik + internal tersimpan
	secret := errors.New("dial postgres://user:rahasia@db:5432 timeout")
	got := From(secret)
	if got.Status() != http.StatusInternalServerError || got.Code() != CodeInternalError {
		t.Errorf("error biasa -> %d/%s, ingin 500/INTERNAL_ERROR", got.Status(), got.Code())
	}
	if got.Internal() == nil || !strings.Contains(got.Internal().Error(), "rahasia") {
		t.Errorf("detail internal harus tersimpan untuk log, bukan di pesan klien")
	}
	if strings.Contains(got.Error(), "rahasia") {
		t.Errorf("pesan klien membocorkan detail internal: %q", got.Error())
	}

	// context.DeadlineExceeded -> 500 terkontrol
	if d := From(context.DeadlineExceeded); d.Status() != http.StatusInternalServerError {
		t.Errorf("DeadlineExceeded -> %d, ingin 500", d.Status())
	}

	// *http.MaxBytesError -> 413
	mbe := From(&http.MaxBytesError{Limit: 10})
	if mbe.Status() != http.StatusRequestEntityTooLarge || mbe.Code() != CodePayloadTooLarge {
		t.Errorf("MaxBytesError -> %d/%s, ingin 413/PAYLOAD_TOO_LARGE", mbe.Status(), mbe.Code())
	}
}

func TestDecodeJSON_Strict(t *testing.T) {
	type payload struct {
		Title string `json:"title"`
	}

	t.Run("valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"ok"}`))
		var p payload
		if ae := DecodeJSON(req, &p); ae != nil {
			t.Fatalf("decode valid gagal: %v", ae.Error())
		}
		if p.Title != "ok" {
			t.Errorf("title = %q", p.Title)
		}
	})

	t.Run("field_asing_ditolak", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"ok","extra":"boom"}`))
		var p payload
		ae := DecodeJSON(req, &p)
		if ae == nil || ae.Code() != CodeMalformedRequest {
			t.Fatalf("field asing harus 400 MALFORMED_REQUEST, got %v", ae)
		}
	})

	t.Run("trailing_data_ditolak", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"ok"} {"title":"lagi"}`))
		var p payload
		ae := DecodeJSON(req, &p)
		if ae == nil || ae.Code() != CodeMalformedRequest {
			t.Fatalf("trailing data harus 400 MALFORMED_REQUEST, got %v", ae)
		}
	})

	t.Run("json_rusak_ditolak", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bukan json`))
		var p payload
		ae := DecodeJSON(req, &p)
		if ae == nil || ae.Code() != CodeMalformedRequest {
			t.Fatalf("json rusak harus 400 MALFORMED_REQUEST, got %v", ae)
		}
	})

	t.Run("body_berlebih_413", func(t *testing.T) {
		big := `{"title":"` + strings.Repeat("a", 100) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(big))
		req.Body = http.MaxBytesReader(httptest.NewRecorder(), req.Body, 20)
		var p payload
		ae := DecodeJSON(req, &p)
		if ae == nil || ae.Code() != CodePayloadTooLarge {
			t.Fatalf("body berlebih harus 413 PAYLOAD_TOO_LARGE, got %v", ae)
		}
	})
}
