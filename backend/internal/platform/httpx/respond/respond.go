// Package respond menulis envelope respons API_SPEC §2 (data/error + meta)
// secara konsisten dan menyuntikkan request_id ke setiap body.
//
// request_id dibuat oleh middleware RequestID dan disimpan pada gin.Context
// (kunci RequestIDKey) serta pada request context (untuk korelasi log service).
package respond

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"randesk/backend/internal/platform/httpx/apierr"
)

// RequestIDKey adalah kunci gin.Context untuk request_id.
const RequestIDKey = "randesk.request_id"

// ActorIDKey adalah kunci gin.Context untuk actor_id (diisi middleware auth
// pada milestone berikutnya; dipakai logging bila tersedia).
const ActorIDKey = "randesk.actor_id"

type ctxKey struct{}

// RequestID mengambil request_id dari gin.Context (string kosong bila belum ada).
func RequestID(c *gin.Context) string {
	if v, ok := c.Get(RequestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RequestIDFromContext mengambil request_id dari request context.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

// ContextWithRequestID menurunkan context yang membawa request_id.
func ContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// meta membangun blok meta; request_id selalu disertakan.
func meta(c *gin.Context) apierr.Meta {
	return apierr.Meta{RequestID: RequestID(c)}
}

// JSON menulis respons sukses satu objek: {data, meta}.
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, apierr.Response{Data: data, Meta: meta(c)})
}

// Created menulis 201 dengan envelope satu objek.
func Created(c *gin.Context, data any) {
	JSON(c, http.StatusCreated, data)
}

// List menulis respons daftar terpaginasi: {data:[], meta:{page,per_page,total,request_id}}.
func List(c *gin.Context, data any, page, perPage int, total int64) {
	p, pp, t := page, perPage, total
	c.JSON(http.StatusOK, apierr.Response{
		Data: data,
		Meta: apierr.Meta{Page: &p, PerPage: &pp, Total: &t, RequestID: RequestID(c)},
	})
}

// NoContent menulis 204 tanpa body (logout, delete lampiran).
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error memetakan err ke APIError dan menulis envelope error. Detail internal
// (SQL, stack, credential) TIDAK dikirim ke klien; hanya dicatat ke log server.
func Error(c *gin.Context, err error) {
	ae := apierr.From(err)

	if internal := ae.Internal(); internal != nil {
		slog.ErrorContext(c.Request.Context(), "request gagal",
			"request_id", RequestID(c),
			"code", ae.Code(),
			"error", internal.Error(),
		)
	}

	body := apierr.ErrorBody{
		Code:    ae.Code(),
		Message: ae.Error(),
		Fields:  ae.Fields(),
	}
	if cv, ok := ae.CurrentVersion(); ok {
		body.CurrentVersion = &cv
	}
	if ra := ae.RetryAfter(); ra > 0 {
		seconds := int(ra.Seconds())
		if seconds < 1 {
			seconds = 1
		}
		c.Header("Retry-After", strconv.Itoa(seconds))
	}

	c.JSON(ae.Status(), apierr.ErrorResponse{Error: body, Meta: meta(c)})
}
