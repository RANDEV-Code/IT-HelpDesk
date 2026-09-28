// Package apierr memetakan error domain ke respons HTTP sesuai API_SPEC §2–§3
// dan menulis envelope data/error + meta.request_id yang konsisten.
//
// Prinsip SECURITY/ARCHITECTURE §10: pesan yang dikirim ke klien generik dan
// tidak pernah memuat SQL, stack trace, credential, atau keberadaan objek yang
// tidak boleh diketahui actor. Detail internal disimpan di field `internal`
// dan HANYA untuk log server.
package apierr

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Code error yang dikenal API_SPEC §3.
const (
	CodeMalformedRequest     = "MALFORMED_REQUEST"
	CodeUnauthenticated      = "UNAUTHENTICATED"
	CodeForbidden            = "FORBIDDEN"
	CodeCSRFCailed           = "CSRF_FAILED"
	CodePasswordChangeReq    = "PASSWORD_CHANGE_REQUIRED"
	CodeNotFound             = "NOT_FOUND"
	CodeVersionConflict      = "VERSION_CONFLICT"
	CodeAlreadyAssigned      = "ALREADY_ASSIGNED"
	CodeInvalidState         = "INVALID_STATE"
	CodeActiveTicketsExist   = "ACTIVE_TICKETS_EXIST"
	CodeLastAdmin            = "LAST_ADMIN"
	CodeEmailExists          = "EMAIL_EXISTS"
	CodeNameExists           = "NAME_EXISTS"
	CodePayloadTooLarge      = "PAYLOAD_TOO_LARGE"
	CodeUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	CodeValidationError      = "VALIDATION_ERROR"
	CodeMasterInactive       = "MASTER_INACTIVE"
	CodeNoChange             = "NO_CHANGE"
	CodeAttachmentLimit      = "ATTACHMENT_LIMIT"
	CodeRateLimited          = "RATE_LIMITED"
	CodeInternalError        = "INTERNAL_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
)

// pesan default per code (Bahasa Indonesia, generik, aman diekspos).
var defaultMessages = map[string]string{
	CodeMalformedRequest:     "Permintaan tidak dapat diproses.",
	CodeUnauthenticated:      "Sesi tidak valid atau telah berakhir.",
	CodeForbidden:            "Aksi tidak diizinkan.",
	CodeCSRFCailed:           "Verifikasi keamanan gagal; muat ulang halaman.",
	CodePasswordChangeReq:    "Kata sandi sementara wajib diganti terlebih dahulu.",
	CodeNotFound:             "Data tidak ditemukan.",
	CodeVersionConflict:      "Data telah diubah oleh proses lain; muat ulang.",
	CodeAlreadyAssigned:      "Tiket sudah ditangani teknisi lain.",
	CodeInvalidState:         "Status tiket tidak mengizinkan aksi ini.",
	CodeActiveTicketsExist:   "Teknisi masih memiliki tiket aktif.",
	CodeLastAdmin:            "Tidak dapat menghapus admin aktif terakhir.",
	CodeEmailExists:          "Email sudah terdaftar.",
	CodeNameExists:           "Nama sudah digunakan.",
	CodePayloadTooLarge:      "Ukuran permintaan melebihi batas.",
	CodeUnsupportedMediaType: "Tipe konten tidak didukung.",
	CodeValidationError:      "Periksa kembali data yang dikirim.",
	CodeMasterInactive:       "Referensi master data tidak aktif.",
	CodeNoChange:             "Tidak ada perubahan yang dikirim.",
	CodeAttachmentLimit:      "Jumlah lampiran melebihi batas.",
	CodeRateLimited:          "Terlalu banyak permintaan; coba lagi nanti.",
	CodeInternalError:        "Terjadi kesalahan pada server.",
	CodeServiceUnavailable:   "Layanan belum siap; coba lagi sebentar.",
}

// APIError adalah error domain yang dapat dipetakan ke respons HTTP.
type APIError struct {
	status         int
	code           string
	message        string
	fields         map[string]string
	currentVersion *int64
	retryAfter     time.Duration
	internal       error // hanya untuk log server; tidak pernah dikirim ke klien
}

// Error mengimplementasikan interface error. Pesan yang dikembalikan adalah
// pesan aman untuk klien (bukan detail internal).
func (e *APIError) Error() string {
	if e.message != "" {
		return e.message
	}
	return defaultMessages[e.code]
}

// Code mengembalikan kode error API_SPEC.
func (e *APIError) Code() string { return e.code }

// Status mengembalikan status HTTP.
func (e *APIError) Status() int { return e.status }

// Fields mengembalikan pesan error per field (nil bila bukan error field).
func (e *APIError) Fields() map[string]string { return e.fields }

// Internal mengembalikan error terbungkus untuk logging (boleh nil).
func (e *APIError) Internal() error { return e.internal }

// WithInternal menautkan detail internal (untuk log) tanpa mengubah pesan klien.
func (e *APIError) WithInternal(err error) *APIError {
	e.internal = err
	return e
}

// WithMessage mengesampingkan pesan default dengan pesan generik khusus.
func (e *APIError) WithMessage(msg string) *APIError {
	e.message = msg
	return e
}

// new membangun APIError dengan pesan default code.
func new(status int, code string) *APIError {
	return &APIError{status: status, code: code, message: defaultMessages[code]}
}

// --- Konstruktor error domain ---

func MalformedRequest() *APIError { return new(http.StatusBadRequest, CodeMalformedRequest) }
func Unauthenticated() *APIError  { return new(http.StatusUnauthorized, CodeUnauthenticated) }
func Forbidden() *APIError        { return new(http.StatusForbidden, CodeForbidden) }
func CSRFFailed() *APIError       { return new(http.StatusForbidden, CodeCSRFCailed) }
func PasswordChangeRequired() *APIError {
	return new(http.StatusForbidden, CodePasswordChangeReq)
}
func NotFound() *APIError { return new(http.StatusNotFound, CodeNotFound) }

// VersionConflict menyertakan current_version opsional (API_SPEC §3) bila > 0.
func VersionConflict(currentVersion int64) *APIError {
	e := new(http.StatusConflict, CodeVersionConflict)
	if currentVersion > 0 {
		cv := currentVersion
		e.currentVersion = &cv
	}
	return e
}
func AlreadyAssigned() *APIError    { return new(http.StatusConflict, CodeAlreadyAssigned) }
func InvalidState() *APIError       { return new(http.StatusConflict, CodeInvalidState) }
func ActiveTicketsExist() *APIError { return new(http.StatusConflict, CodeActiveTicketsExist) }
func LastAdmin() *APIError          { return new(http.StatusConflict, CodeLastAdmin) }
func EmailExists() *APIError        { return new(http.StatusConflict, CodeEmailExists) }
func NameExists() *APIError         { return new(http.StatusConflict, CodeNameExists) }
func PayloadTooLarge() *APIError    { return new(http.StatusRequestEntityTooLarge, CodePayloadTooLarge) }
func UnsupportedMediaType() *APIError {
	return new(http.StatusUnsupportedMediaType, CodeUnsupportedMediaType)
}

// ValidationError membawa pesan per field (API_SPEC §2 contoh error field).
func ValidationError(fields map[string]string) *APIError {
	e := new(http.StatusUnprocessableEntity, CodeValidationError)
	e.fields = fields
	return e
}
func MasterInactive() *APIError  { return new(http.StatusUnprocessableEntity, CodeMasterInactive) }
func NoChange() *APIError        { return new(http.StatusUnprocessableEntity, CodeNoChange) }
func AttachmentLimit() *APIError { return new(http.StatusUnprocessableEntity, CodeAttachmentLimit) }

// RateLimited menyertakan Retry-After (API_SPEC §3: 429 berikan Retry-After).
func RateLimited(retryAfter time.Duration) *APIError {
	e := new(http.StatusTooManyRequests, CodeRateLimited)
	e.retryAfter = retryAfter
	return e
}
func Internal() *APIError { return new(http.StatusInternalServerError, CodeInternalError) }
func ServiceUnavailable() *APIError {
	return new(http.StatusServiceUnavailable, CodeServiceUnavailable)
}

// RetryAfter mengembalikan durasi Retry-After (0 bila tidak diset).
func (e *APIError) RetryAfter() time.Duration { return e.retryAfter }

// CurrentVersion mengembalikan current_version (ok=false bila tidak diset).
func (e *APIError) CurrentVersion() (int64, bool) {
	if e.currentVersion == nil {
		return 0, false
	}
	return *e.currentVersion, true
}

// From memetakan error apa pun ke APIError. Bila sudah *APIError, dikembalikan
// apa adanya; context.DeadlineExceeded dipetakan ke 500 terkontrol; selain itu
// dibungkus menjadi 500 INTERNAL_ERROR generik dengan detail asli disimpan
// sebagai internal (untuk log saja).
func From(err error) *APIError {
	if err == nil {
		return nil
	}
	var ae *APIError
	if errors.As(err, &ae) {
		return ae
	}
	// body melebihi batas MaxBytesReader -> 413
	var maxBytes *http.MaxBytesError
	if errors.As(err, &maxBytes) {
		return PayloadTooLarge().WithInternal(err)
	}
	// timeout query/handler (context turunan) -> 500 generik terkontrol
	if errors.Is(err, context.DeadlineExceeded) {
		return Internal().WithInternal(err)
	}
	return Internal().WithInternal(err)
}

// --- Envelope respons (API_SPEC §2) ---

// Meta adalah blok meta respons. request_id selalu ada; field paginasi hanya
// untuk respons daftar.
type Meta struct {
	Page      *int   `json:"page,omitempty"`
	PerPage   *int   `json:"per_page,omitempty"`
	Total     *int64 `json:"total,omitempty"`
	RequestID string `json:"request_id"`
}

// ErrorBody adalah blok error pada respons gagal.
type ErrorBody struct {
	Code           string            `json:"code"`
	Message        string            `json:"message"`
	Fields         map[string]string `json:"fields,omitempty"`
	CurrentVersion *int64            `json:"current_version,omitempty"`
}

// ErrorResponse adalah envelope error lengkap.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
	Meta  Meta      `json:"meta"`
}

// Response adalah envelope sukses (satu objek atau daftar).
type Response struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

// DecodeJSON membaca body JSON secara ketat: menolak field asing dan data
// tambahan setelah objek pertama (API_SPEC §1). Kegagalan → MalformedRequest.
func DecodeJSON(r *http.Request, dst any) *APIError {
	if r.Body == nil {
		return MalformedRequest()
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		// body dipotong MaxBytesReader -> 413, bukan 400
		var maxBytes *http.MaxBytesError
		if errors.As(err, &maxBytes) {
			return PayloadTooLarge().WithInternal(err)
		}
		return MalformedRequest().WithInternal(err)
	}
	// tolak trailing data setelah objek JSON pertama
	if dec.More() {
		return MalformedRequest().WithInternal(errors.New("trailing data after JSON object"))
	}
	// pastikan sisa stream benar-benar kosong / hanya whitespace
	var trailing json.RawMessage
	if err := dec.Decode(&trailing); err == nil {
		return MalformedRequest().WithInternal(errors.New("trailing data after JSON object"))
	}
	return nil
}
