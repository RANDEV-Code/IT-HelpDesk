# REST API Specification — RANDesk

Versi 1.0 · 26 September 2026 · Base path: `/api/v1`

## 1. Konvensi kontrak

- JSON UTF-8; nama field snake_case; UUID berupa string.
- Timestamp RFC 3339 UTC, contoh `2026-09-26T04:00:00Z`.
- Field opsional yang tidak dikirim berarti tidak diubah pada PATCH; null hanya diizinkan jika disebutkan.
- Tolak field asing dan request JSON dengan nilai tambahan setelah objek pertama.
- Semua endpoint selain login dan health membutuhkan sesi. Semua mutasi selain login memerlukan Origin valid serta `X-CSRF-Token`.
- Respons sensitif memakai `Cache-Control: no-store`; jangan cache API di service worker.
- Nomor internal BIGINT tidak dikirim sebagai JavaScript number; gunakan `ticket_code` string. Event seq juga string jika diekspos.
- `request_id` dibuat server dan disertakan dalam respons JSON serta header `X-Request-ID`.

## 2. Envelope

### Respons satu objek

```json
{
  "data": {
    "id": "11111111-1111-4111-8111-111111111111",
    "ticket_code": "HD-000042",
    "title": "Wi-Fi ruang rapat tidak terhubung",
    "description": "Laptop tidak dapat terhubung sejak pagi, sedangkan perangkat lain normal.",
    "status": "open",
    "priority": "normal",
    "version": 1,
    "requester": {"id": "22222222-2222-4222-8222-222222222222", "name": "Pegawai Demo"},
    "assignee": null,
    "category": {"id": "33333333-3333-4333-8333-333333333333", "name": "Jaringan", "is_active": true},
    "first_response_at": null,
    "resolved_at": null,
    "resolution_summary": null,
    "closed_at": null,
    "reopen_count": 0,
    "created_at": "2026-09-26T04:00:00Z",
    "updated_at": "2026-09-26T04:00:00Z",
    "allowed_actions": ["edit", "comment_public", "upload_attachment"]
  },
  "meta": {"request_id": "44444444-4444-4444-8444-444444444444"}
}
```

`allowed_actions` adalah bantuan UI yang dihitung untuk actor saat ini; backend tetap memeriksa ulang saat mutasi. Nilai yang mungkin: edit, change_priority, assign, claim, unassign, start, resolve, close, reopen, comment_public, comment_internal, upload_attachment. Izin hapus tiap attachment dikirim sebagai `can_delete` pada DTO attachment.

### Respons daftar

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 0,
    "request_id": "44444444-4444-4444-8444-444444444444"
  }
}
```

Semua daftar menggunakan page ≥ 1, default per_page 20, maksimum 100. Total dihitung setelah scope dan filter. Daftar comment/event/notification/audit juga dipaginasi. Tidak ada pagination untuk detail atau dashboard summary.

### Respons error

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Periksa kembali data yang dikirim.",
    "fields": {"title": "Judul harus berisi 5–150 karakter."}
  },
  "meta": {"request_id": "44444444-4444-4444-8444-444444444444"}
}
```

`fields` hanya muncul untuk error field. Error tidak berisi SQL, stack trace, credential, atau keberadaan objek yang tidak boleh diketahui actor.

## 3. Status dan error codes

| HTTP | Code contoh | Pemakaian |
| --- | --- | --- |
| 200 | — | Read/update/login berhasil |
| 201 | — | Tiket, komentar, lampiran, akun, atau master dibuat |
| 204 | — | Logout atau delete lampiran berhasil, tanpa body |
| 400 | MALFORMED_REQUEST | JSON rusak, field asing, ID/parameter tidak dapat diparse |
| 401 | UNAUTHENTICATED | Sesi tidak ada, kedaluwarsa, dicabut, atau akun nonaktif |
| 403 | FORBIDDEN / CSRF_FAILED / PASSWORD_CHANGE_REQUIRED | Role/aksi tidak diizinkan atau kontrol sesi belum dipenuhi |
| 404 | NOT_FOUND | Objek tidak ada atau tidak berada dalam scope baca actor |
| 409 | VERSION_CONFLICT / ALREADY_ASSIGNED / INVALID_STATE | Konflik perubahan, claim kalah, atau status tidak cocok |
| 409 | ACTIVE_TICKETS_EXIST / LAST_ADMIN / EMAIL_EXISTS / NAME_EXISTS | Konflik integritas bisnis |
| 413 | PAYLOAD_TOO_LARGE | Batas body/file terlampaui |
| 415 | UNSUPPORTED_MEDIA_TYPE | Content-Type atau jenis file tidak diterima |
| 422 | VALIDATION_ERROR / MASTER_INACTIVE / NO_CHANGE / ATTACHMENT_LIMIT | Input dapat diparse tetapi melanggar aturan |
| 429 | RATE_LIMITED | Batas request terlampaui; berikan Retry-After |
| 500 | INTERNAL_ERROR | Kegagalan tak terduga dengan pesan generik |
| 503 | SERVICE_UNAVAILABLE | DB/storage tidak siap; tanpa detail infrastruktur |

Untuk versi stale pada objek yang boleh diakses, error boleh menambahkan `current_version`. Jangan memberikan objek terbaru secara otomatis jika izin aktor telah berubah.

## 4. Autentikasi

| Method dan endpoint | Input | Hasil |
| --- | --- | --- |
| POST `/auth/login` | email, password | 200; Set-Cookie sesi, data user dan csrf_token |
| GET `/auth/me` | Cookie | 200; user dan csrf_token untuk bootstrap aplikasi |
| POST `/auth/logout` | `{}` dan CSRF | 204; revoke sesi saat ini dan hapus cookie |
| POST `/auth/change-password` | current_password, new_password | 204; ubah hash, matikan must_change_password, revoke seluruh sesi, hapus cookie |

User DTO sesi: id, name, email, role, is_active, must_change_password, serta department `{id,name}` atau null. Tidak ada password_hash atau token_hash.

Sesi hidup maksimal 8 jam dari login, tanpa sliding expiry pada MVP. CSRF token diberikan JSON, bukan token autentikasi. Login memerlukan JSON, pemeriksaan Origin, dan rate limit; tidak membutuhkan token CSRF sesi karena sesi belum ada.

Ketika `must_change_password=true`, hanya me, logout, dan change-password yang dapat diakses. Login mengganti sesi lama pada browser tersebut jika ada. Reset password administratif dijabarkan pada bagian akun.

## 5. Akun, master data, dan audit

| Method dan endpoint | Akses | Kontrak utama |
| --- | --- | --- |
| GET `/users` | Admin | q, role, is_active, department_id, page, per_page |
| POST `/users` | Admin | name, email, role, department_id opsional/null; 201 dengan user dan temporary_password sekali tampil |
| PATCH `/users/:id` | Admin | name, department_id nullable, is_active; role/email tidak dapat diubah pada MVP |
| POST `/users/:id/reset-password` | Admin | `{}`; 200 temporary_password sekali tampil; cabut semua sesi; tidak untuk diri sendiri |
| GET `/assignees` | Staff | q/page/per_page; daftar teknisi aktif dengan id dan name saja |
| GET `/categories` | Semua | include_inactive default false, q/page/per_page |
| POST `/categories` | Admin | name, description opsional; is_active awal true |
| PATCH `/categories/:id` | Admin | name, description nullable, is_active |
| GET `/departments` | Admin | include_inactive default false, q/page/per_page |
| POST `/departments` | Admin | name, description opsional; is_active awal true |
| PATCH `/departments/:id` | Admin | name, description nullable, is_active |
| GET `/audit-logs` | Admin | actor_id, action, from, to, page, per_page |

Password sementara dibuat server secara acak dengan minimal 16 karakter dan tidak masuk log. Setelah respons tersebut, password tidak bisa dibaca ulang; lakukan reset baru jika diperlukan. Admin memakai change-password untuk dirinya sendiri.

Perubahan akun/master dan auditnya commit bersama. Penonaktifan teknisi dengan tiket open/in_progress ditolak 409. Menonaktifkan diri sendiri ditolak 403 `FORBIDDEN`. Role immutable dan larangan menghilangkan admin aktif tetap berlaku.

DTO akun admin mencakup field user sesi serta created_at/updated_at, tanpa secret. DTO master: id, name, description, is_active, created_at, updated_at. Daftar aktif tetap dipaginasi; dropdown harus mendukung pencarian/pemuatan halaman berikutnya.

## 6. Tiket dan pencarian

| Method dan endpoint | Hasil | Izin |
| --- | --- | --- |
| GET `/tickets` | Daftar summary | Scope RULES |
| POST `/tickets` | 201 detail | Semua peran untuk diri sendiri |
| GET `/tickets/:id` | Detail | Scope RULES |
| PATCH `/tickets/:id` | 200 detail terbaru | Requester pemilik pada open; admin selain closed |
| PUT `/tickets/:id/assignee` | 200 detail terbaru | Admin atau claim technician |
| POST `/tickets/:id/transitions` | 200 detail terbaru | Sesuai aksi dan lifecycle |

Daftar menerima q, status, priority, category_id, from, to, page, per_page, sort. Staff juga boleh mengirim requester_id dan assignee_id; assignee_id dapat berupa UUID, `me`, atau `unassigned`. Requester yang mengirim filter staff mendapat 403, bukan scope yang diperluas.

`from` dan `to` adalah RFC 3339 dengan `from < to`, harus dikirim bersama, serta mengacu pada created_at, from inklusif/to eksklusif. Tanpa keduanya seluruh rentang data berlaku. Sort allowlist: `created_at:desc` (default), `created_at:asc`, `updated_at:desc`, `updated_at:asc`; ID digunakan sebagai tie-breaker dengan arah sama.

Search q maksimal 100 karakter pada nomor, judul, dan deskripsi; tidak mencari komentar. DTO summary berisi id, ticket_code, title, status, priority, version, requester, assignee, category, created_at, updated_at. Field description dan resolution_summary hanya pada detail.

### Buat tiket

```json
{
  "title": "Wi-Fi ruang rapat tidak terhubung",
  "description": "Laptop tidak dapat terhubung sejak pagi, sedangkan perangkat lain normal.",
  "category_id": "33333333-3333-4333-8333-333333333333",
  "priority": "normal"
}
```

Priority opsional dengan default normal. Server mengisi reporter, nomor, status open, version 1, dan timestamps. Tidak menerima requester_id, assignee_id, atau status dari klien. Lampiran menggunakan request terpisah setelah ID tiket tersedia.

### Edit metadata

```json
{
  "expected_version": 2,
  "priority": "high",
  "reason": "Gangguan berdampak pada rapat seluruh departemen."
}
```

Field mutable: title, description, category_id; admin juga priority. Minimal satu field mutable wajib ada. Reason wajib bila priority berubah. Changed fields, alasan, dan nilai kategori/prioritas yang relevan disimpan pada ticket.updated. Form dengan expected_version stale ditolak 409.

### Assignment

```json
{
  "expected_version": 2,
  "assignee_id": "55555555-5555-4555-8555-555555555555",
  "reason": "Dialihkan ke teknisi yang menangani jaringan kantor."
}
```

Technician hanya dapat menggunakan ID dirinya sendiri pada tiket open tanpa assignee. Admin dapat initial assignment dan reassignment pada open/in_progress; null berarti unassign dan hanya boleh pada open. Reason wajib saat mengganti atau melepas assignee yang sudah ada; initial assignment/claim tidak mewajibkannya.

Jika claim bersaing, request yang kalah mendapat 409 VERSION_CONFLICT atau ALREADY_ASSIGNED sesuai pemeriksaan yang pertama gagal. Klien menangani keduanya sebagai tiket perlu dimuat ulang.

### Transisi

```json
{
  "expected_version": 4,
  "action": "resolve",
  "resolution_summary": "Konfigurasi adapter jaringan diperbarui dan koneksi berhasil diuji bersama pelapor."
}
```

| Action | Input tambahan | Aturan |
| --- | --- | --- |
| start | Tidak ada | Open → in_progress, assignee/admin |
| resolve | resolution_summary | In_progress → resolved, assignee/admin |
| close | Tidak ada | Resolved → closed, reporter/admin |
| reopen | reason | Resolved → open, reporter/admin; reset assignment dan field resolusi |

Tolak field tambahan yang tidak relevan dengan action. Tidak ada PATCH status umum, direct close dari open, atau reopen dari closed.

## 7. Komentar, lampiran, dan event

| Method dan endpoint | Input/hasil |
| --- | --- |
| GET `/tickets/:id/comments` | visibility public/internal opsional, page/per_page; scope diterapkan sebelum count |
| POST `/tickets/:id/comments` | body, visibility default public; 201 DTO komentar |
| GET `/tickets/:id/attachments` | Daftar metadata lampiran aktif, page/per_page |
| POST `/tickets/:id/attachments` | multipart/form-data, tepat satu field file; 201 DTO attachment |
| GET `/tickets/:id/attachments/:attachment_id/download` | Binary dengan Content-Disposition attachment dan nosniff |
| DELETE `/tickets/:id/attachments/:attachment_id` | 204 setelah soft delete dan event commit |
| GET `/tickets/:id/events` | page/per_page; terbaru dahulu berdasarkan seq; visibility server |

Komentar tidak membutuhkan expected_version. Otorisasi dan status diperiksa di bawah ticket lock. Requester yang meminta visibility internal mendapat 403; tanpa filter, requester hanya mendapat public. Staff boleh menerima keduanya bila visibility tidak disertakan.

DTO komentar: id, ticket_id, author `{id,name}`, visibility, body, created_at. Tidak ada endpoint update/delete komentar.

DTO attachment: id, ticket_id, original_name, content_type, size_bytes, uploader `{id,name}`, created_at, can_delete, serta download_path relatif API. Storage key dan checksum tidak dikirim ke browser. Nested attachment_id harus cocok dengan ticket_id dan akses actor; objek dihapus/tidak sesuai mendapat 404.

Ukuran maksimum file 5 MiB, request multipart maksimum 6 MiB termasuk overhead. Batasi satu file per request; maksimal 5 metadata aktif per tiket di bawah row lock. Attachment limit menghasilkan 422, ukuran body/file berlebih 413, MIME yang tidak cocok 415.

DTO event: id, seq string, event_type, visibility, actor `{id,name}`, payload allowlist, created_at. Public projection tidak berisi internal events. Payload mengikuti SCHEMA.

## 8. Notifikasi

| Method dan endpoint | Kontrak |
| --- | --- |
| GET `/notifications` | unread_only boolean opsional, page/per_page; meta juga memuat unread_count seluruh notifikasi milik actor |
| PATCH `/notifications/:id/read` | `{}`; 200 DTO dengan read_at; pengulangan idempoten dan mempertahankan read_at pertama |

DTO: id, event_id, ticket_id dari join event, ticket_code, message, read_at, created_at. Semua query dibatasi recipient_id sesi. Mark-read objek orang lain mengembalikan 404. Tidak ada endpoint menerima recipient_id dari requester untuk membuat notifikasi.

## 9. Dashboard

GET `/dashboard/summary` menerima from/to dengan aturan cohort yang sama dengan daftar tiket. Backend menggunakan scope requester atau staff, bukan scope yang dikirim melalui body.

```json
{
  "data": {
    "status_counts": {"open": 8, "in_progress": 5, "resolved": 3, "closed": 12},
    "active_backlog": 13,
    "average_first_response_seconds": 1800,
    "first_response_sample_count": 20,
    "average_resolution_seconds": 21600,
    "resolution_sample_count": 15,
    "category_counts": [
      {"category_id": "33333333-3333-4333-8333-333333333333", "name": "Jaringan", "count": 28}
    ]
  },
  "meta": {"request_id": "44444444-4444-4444-8444-444444444444"}
}
```

Angka di atas merupakan contoh payload sintetis. Rata-rata tanpa sampel menggunakan null dan sample_count 0. Rumus otoritatif berada pada PRD.

## 10. Retry dan konsistensi

Read boleh diulang dengan backoff terbatas. Jangan otomatis mengulang create-ticket/comment/upload jika hasilnya tidak diketahui: request mungkin sudah commit. MVP belum memiliki idempotency-key pada operasi create.

Core mutation dengan expected_version tidak akan mengeksekusi perubahan yang sama dua kali pada versi yang sama. Setelah jaringan pulih, ambil detail/event untuk memeriksa hasil. Mark-read bersifat idempoten; logout dengan sesi sudah hilang dapat diperlakukan frontend sebagai berhasil setelah cookie/cache lokal dibersihkan.

Kontrak machine-readable OpenAPI merupakan deliverable implementasi pada ROADMAP; file Markdown ini adalah acuan awal dan harus diselaraskan saat OpenAPI dibuat.
