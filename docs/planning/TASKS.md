# TASKS — Backlog Implementasi MVP RANDesk

Versi 1.0 · 26 September 2026 · Seluruh task berstatus awal `todo` kecuali ditandai lain.
ID task stabil (TASK-001…TASK-054); jangan menggunakan ulang ID yang dibatalkan. Semua path file di bagian "File rencana" adalah **usulan** mengikuti ARCHITECTURE.md §3 — belum ada yang eksis di repository.

**Legenda status:** `todo` · `blocked` · `needs_verification` · `done` (wajib bukti: hasil tes/log/artefak, sesuai TEST_PLAN §9).

## 0. Ringkasan urutan pengerjaan

| ID | Judul singkat | Milestone | Dependensi | Estimasi | Status |
| --- | --- | --- | --- | --- | --- |
| TASK-001 | Inisialisasi repo, struktur, konfigurasi lingkungan | M0 | — | 2–4 jam | done |
| TASK-002 | Skeleton backend Go (cmd/api, health, shutdown) | M0 | 001 | 3–5 jam | done |
| TASK-003 | PostgreSQL dev + migrasi baseline DDL | M0 | 001 | 3–5 jam | done |
| TASK-004 | Platform: request ID, log, error envelope, timeout | M0 | 002 | 4–6 jam | done |
| TASK-005 | Skeleton frontend React (Vite, TS, Tailwind, shell) | M0 | 001 | 4–6 jam | done |
| TASK-006 | API client frontend + error mapping | M0 | 004, 005 | 3–4 jam | done |
| TASK-007 | CI + harness tes integrasi PostgreSQL | M0 | 002, 003, 005 | 4–6 jam | needs_verification |
| TASK-008 | Password Argon2id + session store | M1 | 003, 004 | 4–6 jam | todo |
| TASK-009 | Middleware auth + CSRF + Origin | M1 | 008 | 4–6 jam | todo |
| TASK-010 | Endpoint auth (login/me/logout/change-password) | M1 | 009 | 4–6 jam | todo |
| TASK-011 | cmd/admin bootstrap admin pertama | M1 | 008 | 2–4 jam | todo |
| TASK-012 | Frontend auth (login, change-password, guard) | M1 | 006, 010 | 4–6 jam | todo |
| TASK-013 | Audit infra (writer transaksional, GET /audit-logs) | M1 | 009 | 3–5 jam | todo |
| TASK-014 | User management backend (CRUD akun, reset password) | M1 | 013 | 4–6 jam | todo |
| TASK-015 | Guard aktivasi akun + advisory lock (BR-15/17) | M1 | 014 | 4–6 jam | todo |
| TASK-016 | Master data backend (kategori, departemen, assignees) | M1 | 013 | 4–6 jam | todo |
| TASK-017 | Seed sintetis + fixture tes minimum | M1 | 003, 016 | 3–5 jam | todo |
| TASK-018 | Admin UI akun | M1 | 012, 015 | 4–6 jam | todo |
| TASK-019 | Admin UI master data + audit | M1 | 012, 016, 013 | 3–5 jam | todo |
| TASK-020 | Policy tiket terpusat + allowed_actions | M2 | 009 | 4–6 jam | todo |
| TASK-021 | Event & notification writer transaksional | M2 | 003, 017 | 4–6 jam | todo |
| TASK-022 | Ticket list/search backend (scope, filter, pagination) | M2 | 016, 020 | 4–6 jam | todo |
| TASK-023 | Ticket create backend (POST /tickets) | M2 | 021, 022 | 3–5 jam | todo |
| TASK-024 | Edit metadata backend (expected_version, NO_CHANGE) | M2 | 023 | 3–5 jam | todo |
| TASK-025 | Assignment backend (claim/assign/reassign/unassign) | M2 | 015, 023 | 4–6 jam | todo |
| TASK-026 | Transitions backend (start/resolve/close/reopen) | M2 | 025 | 4–6 jam | todo |
| TASK-027 | Runtime DB grants (append-only enforcement) | M2 | 003, 021 | 2–3 jam | todo |
| TASK-028 | UI daftar tiket (filter, search, pagination) | M2 | 012, 022 | 4–6 jam | todo |
| TASK-029 | UI form buat tiket | M2 | 023, 028 | 3–5 jam | todo |
| TASK-030 | UI detail tiket + aksi lifecycle + conflict UX | M2 | 024, 025, 026, 028 | 4–6 jam | todo |
| TASK-031 | Comments backend (public/internal, first_response_at) | M3 | 021, 026 | 4–6 jam | todo |
| TASK-032 | Events timeline API | M3 | 021, 023 | 2–4 jam | todo |
| TASK-033 | Storage adapter + validasi file (staging) | M3 | 004 | 4–6 jam | todo |
| TASK-034 | Attachment upload backend (limit + kompensasi) | M3 | 026, 033 | 4–6 jam | todo |
| TASK-035 | Attachment download/delete + cleanup maintenance | M3 | 034 | 3–5 jam | todo |
| TASK-036 | UI komentar + timeline | M3 | 030, 031, 032 | 4–6 jam | todo |
| TASK-037 | UI lampiran (picker, upload per file, list) | M3 | 029, 030, 035 | 4–6 jam | todo |
| TASK-038 | Notifications API (list, unread_count, mark-read) | M4 | 021 | 3–4 jam | todo |
| TASK-039 | UI notifikasi + polling lifecycle | M4 | 012, 038 | 4–6 jam | todo |
| TASK-040 | Dashboard backend (agregat, cohort, metrik) | M4 | 022, 026 | 4–6 jam | todo |
| TASK-041 | Dashboard UI | M4 | 012, 040 | 4–6 jam | todo |
| TASK-042 | Rate limiting + security headers + cookie produksi | M5 | 010, 034 | 3–5 jam | todo |
| TASK-043 | Verifikasi log redaction + request_id tracing | M5 | 004, 010 | 2–4 jam | todo |
| TASK-044 | Sweep negative tests + gate keamanan SECURITY §10 | M5 | 031, 035, 038 | 4–6 jam | todo |
| TASK-045 | Suite konkurensi + failure injection berulang | M5 | 015, 025, 034 | 4–6 jam | todo |
| TASK-046 | Pass aksesibilitas + responsif | M5 | 018, 019, 030, 036, 037, 039, 041 | 4–6 jam | todo |
| TASK-047 | E2E suite (TC-35 + alur kritis) | M5 | 036, 037, 039, 041 | 4–6 jam | todo |
| TASK-048 | Benchmark dataset + load test NFR-03 | M5 | 017, 042, 047 | 4–6 jam | needs_verification |
| TASK-049 | UAT / walkthrough terstruktur | M5 | 047 | 3–5 jam | todo |
| TASK-050 | Deployment staging (proxy TLS, config, smoke test) | M6 | 042, 043, 046 | 4–6 jam | todo |
| TASK-051 | Backup + restore drill (TC-33) | M6 | 035, 050 | 3–5 jam | todo |
| TASK-052 | Perintah maintenance & retensi | M6 | 008, 013, 035 | 3–5 jam | todo |
| TASK-053 | OpenAPI deliverable + rekonsiliasi API_SPEC | M6 | 034, 038, 040 | 3–5 jam | todo |
| TASK-054 | README aplikasi, panduan setup, paket demo portofolio | M6 | 050, 051, 053 | 4–6 jam | todo |

Catatan: TASK-048 `needs_verification` karena lingkungan acuan 2 vCPU/4 GiB (NFR-03) belum terbukti tersedia. TASK-049 tidak lagi `blocked`: D-06 diputuskan (26 September 2026) dengan mengadopsi resmi alternatif walkthrough terstruktur sesuai TEST_PLAN §7 / asumsi P-07 — varian UAT 5 peserta dilepas dari cakupan MVP dan hasil walkthrough wajib dicatat "belum mewakili pengguna nyata".

---

## M0 — Fondasi

### TASK-001 — Inisialisasi repository, struktur folder, dan konfigurasi lingkungan

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (27 September 2026) |
| DEV / ADR | DEV-01; D-02, D-04 (dikunci di task ini) |
| Referensi | ARCHITECTURE §1, §3; RULES §8; SECURITY §8; OPERATIONS §2 |
| Dependensi | Tidak ada |
| Estimasi | 2–4 jam |

**Tujuan:** repository siap dikembangkan: version control, struktur folder target, konfigurasi lingkungan placeholder, dan keputusan versi toolchain terkunci.

**File rencana (usulan):** `.gitignore`, `.env.example`, `backend/` (go.mod kosong-isi), `frontend/` (placeholder), `docs/` (memuat paket perencanaan — lokasi final lihat GAP-01), `Makefile`/`Taskfile` opsional, `AGENTS.md` (konvensi repo untuk agent/developer).

**Langkah implementasi:**
1. Toolchain telah diverifikasi pada tahap planning (26 September 2026): Go 1.26.4 windows/amd64, Node v22.16.0, npm 10.9.2, git 2.37.3 tersedia; **Docker, PostgreSQL, dan make TIDAK terpasang** (Laragon hanya menyediakan MySQL). Sisa pekerjaan: instal PostgreSQL dev (rekomendasi: installer native Windows EDB — Docker tidak tersedia; CI memakai service container sehingga tidak butuh Docker lokal) dan catat versinya.
2. `git init`, commit awal berisi dokumen perencanaan + empat dokumen `docs/planning/`.
3. Buat struktur folder sesuai ARCHITECTURE §3 (backend/cmd, backend/internal, backend/migrations, frontend/src, docs).
4. Tulis `.env.example` berisi nama variabel OPERATIONS §2 dengan placeholder saja (APP_ENV, APP_ORIGIN, HTTP_ADDR, DATABASE_URL, STORAGE_ROOT, SESSION_TTL_HOURS, LOG_LEVEL, TRUSTED_PROXIES, DB_MAX_OPEN_CONNS, DB_MAX_IDLE_CONNS) — tanpa nilai secret.
5. Kunci keputusan D-02 (golang-migrate) dan D-04 (versi Go, Gin, pgx, Node, React, Vite, Tailwind, TanStack Query) — default rutin telah **disetujui pemilik proyek 26 September 2026** (GAPS_AND_DECISIONS.md §4.2); catat versi final yang dikunci di AGENTS.md.
6. Tulis `AGENTS.md` singkat: perintah build/test/lint, konvensi commit, larangan menyimpan secret.

**Acceptance criteria:**
- [x] `git log` menunjukkan commit awal; working tree bersih.
- [x] Seluruh folder target ARCHITECTURE §3 ada (boleh berisi placeholder/README).
- [x] `.env.example` berisi semua variabel OPERATIONS §2, tanpa nilai rahasia.
- [x] Versi toolchain dan library tercatat di satu tempat (AGENTS.md atau DECISIONS bagian D-04 yang diperbarui sebagai usulan).

**Verifikasi & bukti selesai:** tangkapan `git log --stat`, isi `.env.example`, daftar versi.

**Bukti selesai (27 September 2026):** commit `028d09e` (branch `main`, 38 file: dokumen perencanaan + `.gitignore`, `.gitattributes`, `.env.example`, `AGENTS.md`); 18 folder target ARCHITECTURE §3 dibuat dengan `.gitkeep`; `.env.example` memuat 10 variabel OPERATIONS §2 dengan placeholder saja; versi dikunci di AGENTS.md (Go 1.26.4, Node v22.16.0, npm 10.9.2, git 2.37.3; library menyusul di TASK-002). Catatan: PostgreSQL dev baru terpasang saat TASK-003 (PostgreSQL 17.11 via installer EDB).

---

### TASK-002 — Skeleton backend Go: cmd/api, konfigurasi, health, graceful shutdown

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (27 September 2026) |
| DEV | DEV-01, DEV-03 |
| Referensi | ARCHITECTURE §1, §3, §5, §10; OPERATIONS §2; NFR-08 |
| Dependensi | TASK-001 |
| Estimasi | 3–5 jam |

**Tujuan:** API Go dapat start, membaca konfigurasi dari environment, merespons `/health/live` dan `/health/ready`, dan shutdown maksimal 15 detik.

**File rencana (usulan):** `backend/go.mod`, `backend/cmd/api/main.go`, `backend/internal/platform/config/`, `backend/internal/platform/httpx/` (router setup), `backend/internal/platform/db/` (pool pgx stdlib).

**Langkah implementasi:**
1. `go mod init`, tambahkan Gin dan pgx/v5 stdlib pada versi yang dikunci di TASK-001.
2. Buat loader config: baca env, validasi wajib/opsional, gagal start bila APP_ENV=production dan konfigurasi tidak aman (OPERATIONS §2).
3. Setup pool `database/sql` + pgx: DB_MAX_OPEN_CONNS/IDLE, timeout koneksi.
4. Router Gin dengan `/health/live` (proses hidup) dan `/health/ready` (ping DB + cek STORAGE_ROOT writable) tanpa detail rahasia (ARCHITECTURE §10).
5. Graceful shutdown: stop accept, drain maks 15 detik, tutup pool.
6. Tes unit config loader (env hilang → error jelas) dan tes smoke health via `httptest` (ready 503 saat DB tak tersedia).

**Acceptance criteria:**
- [x] `go run ./cmd/api` start dengan `.env` dev; `GET /health/live` → 200.
- [x] `GET /health/ready` → 200 saat DB hidup, 503 saat DB mati, body tanpa connection string.
- [x] Ctrl+C menghentikan proses ≤15 detik tanpa panic.
- [x] Config production tanpa DATABASE_URL menolak start dengan pesan jelas.
- [x] `gofmt`, `go vet ./...`, `go test ./...` lulus.

**Verifikasi & bukti selesai:** output perintah di atas + hasil tes.

**Bukti selesai (27 September 2026):** commit `870f30a`; `go test ./...` lulus semua paket (config, httpx, cmd/api termasuk `TestGracefulShutdown` — sinyal CTRL_BREAK→SIGTERM, shutdown bersih ~4 ms, stdout memuat "api berhenti bersih", tanpa panic); smoke test `smoke-task002.ps1`: live → 200 `{"status":"ok"}`, ready → 503 saat DB mati (body tanpa detail), APP_ENV=production + DATABASE_URL kosong → exit 1 dengan pesan jelas; ready → 200 dibuktikan setelah DB hidup via `ready-check-task003.ps1`. Versi terkunci: Gin v1.12.0, pgx v5.11.0, x/sys v0.41.0 (AGENTS.md).

---

### TASK-003 — PostgreSQL development + tooling migrasi + migrasi baseline

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (27 September 2026) |
| DEV / Keputusan | DEV-02; D-02 dikunci di sini |
| Referensi | SCHEMA §1–§5, §8; RULES §8 (migrasi immutable); TEST_PLAN §1 (integration pakai PostgreSQL nyata) |
| Dependensi | TASK-001 |
| Estimasi | 3–5 jam |

**Tujuan:** database dev PostgreSQL tersedia dan DDL baseline SCHEMA §4 terpasang lewat migrasi bernomor yang dapat diulang dari DB kosong.

**File rencana (usulan):** `backend/migrations/000001_baseline.up.sql` / `.down.sql`, skrip setup dev (`backend/scripts/` atau docker-compose), dokumentasi perintah migrate di AGENTS.md.

**Langkah implementasi:**
1. Siapkan PostgreSQL dev: **keputusan 26 September 2026 — installer native Windows (EDB)** karena Docker tidak terpasang di mesin dev; CI tetap memakai PostgreSQL service container GitHub Actions (disetujui). Buat service/instance lokal + `pg_ctl` terdokumentasi di AGENTS.md.
2. Pasang tool migrasi pilihan (default: golang-migrate CLI); buat dua database: `randesk_dev` dan `randesk_test`.
3. Salin DDL SCHEMA §4 ke migrasi `000001_baseline` apa adanya; jangan mengubah kontrak — deviasi (mis. error versi PostgreSQL tertentu) dicatat di GAPS_AND_DECISIONS.md.
4. Jalankan migrate up pada DB kosong; verifikasi seluruh tabel, indeks, CHECK, dan identity ada (query `information_schema`/`\d`).
5. Uji migrate down → up kembali bersih.
6. Verifikasi privilege sequence/identity yang dibutuhkan driver (SCHEMA §5 catatan) dan catat hasilnya.
7. Buat role DB terpisah: migration owner vs runtime role (grants penuh ditunda ke TASK-027).

**Acceptance criteria:**
- [x] Migrasi dari DB kosong berhasil tanpa error; seluruh 10 tabel + indeks SCHEMA §4 ada.
- [x] CHECK constraint teruji manual minimal: status invalid ditolak, version 0 ditolak, email tidak lowercase ditolak.
- [x] Down migration menghapus bersih; re-up berhasil.
- [x] Perintah migrate terdokumentasi dan dapat dijalankan ulang oleh orang lain.

**Verifikasi & bukti selesai:** log eksekusi migrate up/down, hasil query verifikasi tabel/constraint.

**Bukti selesai (27 September 2026):** PostgreSQL 17.11 terpasang (service `postgresql-x64-17`); golang-migrate v4.20.1 di-build dengan `-tags postgres`; `setup-task003.ps1` berakhir "SEMUA LANGKAH TASK-003 LULUS": 3 role (`randesk_migrate`, `randesk_runtime`, `randesk_operator` NOLOGIN) + DB `randesk_dev`/`randesk_test`; verifikasi 11 tabel (10 + `schema_migrations`), 40 indeks, 32 CHECK constraint; uji constraint via `check-constraints-task003.sql` — status invalid ditolak (`tickets_lifecycle_ck`), version 0 ditolak (`tickets_version_check`), email tidak lowercase ditolak (`users_email_normalized_ck`), deskripsi pendek ditolak (`tickets_description_check`), kontrol positif INSERT 0 1; identity `tickets.ticket_no` & `ticket_events.seq` = GENERATED ALWAYS; `down 1` menyisakan hanya `schema_migrations` lalu re-up bersih; `randesk_test` berada di version 1. Deviasi dicatat di GAPS_AND_DECISIONS.md §10.

---

### TASK-004 — Platform HTTP: request ID, logging terstruktur, error envelope, batas body, timeout

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (28 September 2026) |
| DEV | DEV-03 |
| Referensi | ARCHITECTURE §5, §10; API_SPEC §1–§3; NFR-05, NFR-08; SECURITY §8 (log tanpa body sensitif) |
| Dependensi | TASK-002 |
| Estimasi | 4–6 jam |

**Tujuan:** seluruh respons API memakai envelope `data`/`error` + `meta.request_id` yang konsisten; middleware platform (recovery, body limit, timeout, log durasi) terpasang dan teruji.

**File rencana (usulan):** `backend/internal/platform/httpx/middleware/` (requestid, logging, recovery, bodylimit, timeout), `backend/internal/platform/httpx/errors/` (domain error → HTTP mapping), `backend/internal/platform/logger/`.

**Langkah implementasi:**
1. Middleware request ID: UUID per request, header `X-Request-ID` + `meta.request_id` di body (API_SPEC §1).
2. Logger JSON terstruktur: timestamp, level, request_id, route template (bukan path mentah), method, status, duration_ms, actor_id bila ada (ARCHITECTURE §10). Tanpa body/cookie/header sensitif.
3. Recovery middleware → 500 `INTERNAL_ERROR` generik + log stack internal saja.
4. Batas body: JSON 64 KiB → 413; siapkan hook batas multipart 6 MiB untuk M3 (SECURITY §7).
5. Timeout: query 3 detik, handler JSON 10 detik (context turunan) — nilai dari ARCHITECTURE §5.
6. Paket error domain: tipe error (Validation, Unauthorized, Forbidden, NotFound, VersionConflict, NoChange, MasterInactive, AttachmentLimit, RateLimited, dst. sesuai API_SPEC §3) + mapper ke status/code/fields; decode JSON dengan `DisallowUnknownFields` dan penolakan trailing data (API_SPEC §1 → 400 MALFORMED_REQUEST).
7. Tes: envelope tiap kelas error; field asing → 400; body >64 KiB → 413; timeout menghasilkan 500/504 terkontrol; log tidak memuat body.

**Acceptance criteria:**
- [x] Setiap respons (sukses/error) memuat `meta.request_id` yang sama dengan header `X-Request-ID`.
- [x] Tabel error API_SPEC §3 terwakili mapper dengan code yang tepat (diuji table-driven).
- [x] JSON field asing/extra ditolak 400 tanpa mengeksekusi handler.
- [x] Log output JSON mengandung field wajib dan tidak mengandung nilai body.

**Verifikasi & bukti selesai:** hasil `go test ./internal/platform/...`, contoh output log.

**Bukti selesai (28 September 2026):** `go test ./...` lulus semua paket; `gofmt -l` bersih; `go vet ./...` bersih. Paket baru: `httpx/apierr` (envelope + 22 konstruktor error + mapper `From`), `httpx/middleware` (RequestID/Recovery/BodyLimit/Logging/Timeout), `httpx/respond`, `platform/logger`. Tes table-driven `TestErrorMappingTableDriven` memetakan 22 kelas error ke status+code tepat. Bukti perilaku: request_id header == `meta.request_id` (`TestRequestID_HeaderMatchesEnvelopeMeta`), UUID unik antar-request; field asing → 400 MALFORMED_REQUEST (`TestDecodeJSON_UnknownFieldRejected400`, `TestDecodeJSON_Strict/field_asing_ditolak`), trailing data → 400; body berlebih → 413 PAYLOAD_TOO_LARGE (`TestBodyLimit_OversizedRejected413`); panic → 500 INTERNAL_ERROR tanpa stack di body klien sementara stack tercatat di log server (`TestRecovery_PanicBecomes500WithoutStack`); timeout → 500 terkontrol (kooperatif + backstop). Log JSON terverifikasi memuat `request_id/method/route(template)/status/duration_ms/time/level` dan TIDAK memuat nilai body (`TestLogging_RequiredFieldsAndNoBody`). Contoh output log nyata (server `APP_ENV=staging`, `scripts/sample-log-task004.ps1`):
```
X-Request-ID: 258a911e-9012-4cd0-8e58-431ebe020e13
{"time":"2026-09-28T23:48:02.633+08:00","level":"INFO","msg":"http request","request_id":"258a911e-9012-4cd0-8e58-431ebe020e13","method":"GET","route":"/health/live","status":200,"duration_ms":0}
```
(request_id header == request_id log). Tidak ada dependensi baru: UUID v4 dihasilkan lewat `crypto/rand` stdlib. Deviasi: timeout handler dipetakan ke **500 INTERNAL_ERROR** (bukan 504) karena API_SPEC §3 tidak mendefinisikan code untuk 504; pendekatan timeout kooperatif via context turunan (bukan goroutine pembatalan) sesuai ARCHITECTURE §5 — dicatat di GAPS_AND_DECISIONS.md §10.

---

### TASK-005 — Skeleton frontend: Vite + TypeScript strict + Tailwind + Router + TanStack Query + AppShell

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (29 September 2026) |
| DEV | DEV-01 |
| Referensi | ARCHITECTURE §1, §7; DESIGN §4 (tokens), §5 (layout), §7 (komponen); RULES §8 |
| Dependensi | TASK-001 |
| Estimasi | 4–6 jam |

**Tujuan:** SPA shell berjalan dengan TypeScript strict, design tokens DESIGN §4, layout responsif §5, dan proxy dev `/api` ke backend.

**File rencana (usulan):** `frontend/package.json`, `frontend/vite.config.ts`, `frontend/tsconfig.json`, `frontend/tailwind` config/theme tokens, `frontend/src/app/` (router, providers, AppShell/Sidebar/Topbar), `frontend/src/components/` (Button, Input, LoadingState, ErrorState, EmptyState versi awal), `frontend/src/lib/` (date utils WITA).

**Langkah implementasi:**
1. Scaffold Vite React-TS; aktifkan strict; pasang Tailwind, React Router, TanStack Query pada versi terkunci.
2. Konfigurasi proxy dev `/api` → `HTTP_ADDR` backend (ARCHITECTURE §7 — satu origin).
3. Implementasi design tokens sebagai Tailwind theme: warna, tipografi, spacing, radius, focus ring, target 44 px (DESIGN §4).
4. AppShell: sidebar 240 px ≥1024 px, drawer <1024 px, topbar 64 px, max-width konten 1440 px (DESIGN §5).
5. Route placeholder sesuai DESIGN §3 (`/login`, `/dashboard`, `/tickets`, dst.) tanpa logika akses dulu.
6. Utilitas tanggal: format UTC → Asia/Makassar dengan label WITA (BR-13).
7. Komponen state dasar: LoadingState, EmptyState, ErrorState; verifikasi halaman demo memanggil `/health/live` lewat proxy.
8. Gate: `npm run lint`, `typecheck`, `build` lulus.

**Acceptance criteria:**
- [x] Dev server menampilkan shell; resize 360/768/1024/1440 px sesuai perilaku DESIGN §5.
- [x] Panggilan `/api` dari browser dev mencapai backend (bukti: health check tampil).
- [x] `tsc --noEmit` tanpa error dengan strict; lint dan production build lulus.
- [x] Token warna/tipografi DESIGN §4 terpakai (bukan default framework mentah).
- [x] Tanggal contoh dirender dalam WITA dengan label.

**Verifikasi & bukti selesai:** screenshot shell 3 viewport, output lint/typecheck/build.

**Bukti selesai (29 September 2026):** `npm run typecheck` (`tsc --noEmit` strict) bersih; `npm run lint` 0 problem; `npm run test` (vitest) 5/5 lulus; `npm run build` sukses (dist index 0,42 kB + CSS 12,10 kB berisi token + JS 352,96 kB). Proxy Vite diverifikasi end-to-end: `Invoke-WebRequest http://localhost:5173/health/live` → `200 {"status":"ok"}` dengan header `X-Request-ID` (artinya request menembus ke Go + middleware TASK-004 aktif); di browser badge dashboard menampilkan **"Backend hidup: ok"**. WITA: `2026-09-26T04:00:00Z` dirender **"26 Sep 2026 12.00 WITA"** (UTC+8, benar). Responsif DESIGN §5 diverifikasi lewat pengukuran DOM/computed-style tiap lebar (via iframe same-origin): sidebar `w-60`=240px `position:static` terlihat permanen di 1440 & 1024; tersembunyi (`x=-240`) dengan tombol hamburger 44×44 `lg:hidden` di 768 & 360; klik hamburger → `aria-expanded` false→true, aside `translate-x-0` + backdrop muncul. Screenshot asli tersimpan 1 (docs/testing/probe3.png, lebar 976px mode drawer — menegakkan hamburger + badge + WITA); screenshot 4 lebar lain **terhambat keterbatasan otomasi browser** (tab `visibilityState=hidden` → rAF tidak fire → capture timeout), bukan bug aplikasi; perilaku telah dibuktikan terukur di atas. Versi frontend terkunci di `package-lock.json` (D-04).

---

### TASK-006 — API client frontend: fetch wrapper, envelope/error mapping, CSRF memory

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / done (29 September 2026) |
| DEV | DEV-01 |
| Referensi | API_SPEC §1–§3, §10; SECURITY §4 (CSRF token di memori); ARCHITECTURE §7 |
| Dependensi | TASK-004, TASK-005 |
| Estimasi | 3–4 jam |

**Tujuan:** satu jalur akses HTTP di frontend yang mem-parse envelope, memetakan error ke bentuk UI (field errors, conflict, unauthorized), mengirim `X-CSRF-Token`, dan mendukung retry read terbatas.

**File rencana (usulan):** `frontend/src/lib/api/client.ts`, `frontend/src/lib/api/errors.ts`, `frontend/src/lib/api/csrf.ts` (store in-memory), `frontend/src/app/providers.tsx` (QueryClient default).

**Langkah implementasi:**
1. Wrapper fetch: credentials same-origin, JSON strict parse, timeout client-side.
2. Mapping error: `error.code` + `error.fields` → tipe error TS (VALIDATION_ERROR dengan field map, VERSION_CONFLICT dengan current_version opsional, UNAUTHENTICATED → trigger alur logout, 403, 404, 429 dengan Retry-After).
3. CSRF: simpan token hanya di memori (module state), attach header untuk POST/PUT/PATCH/DELETE; hook re-fetch `/auth/me` setelah reload (endpoint tersedia di TASK-010 — untuk M0 cukup interface + tes mock).
4. Kebijakan retry TanStack Query: retry read dengan backoff terbatas; jangan retry mutasi create/comment/upload (API_SPEC §10).
5. Query key convention menyertakan user ID (ARCHITECTURE §7) + helper invalidation.
6. Tes komponen/unit (vitest): mapping tiap code, header CSRF terkirim, mutasi tidak di-retry.

**Acceptance criteria:**
- [x] Error 422 VALIDATION_ERROR terpetakan ke struktur field yang bisa dirender form.
- [x] 409 VERSION_CONFLICT terpetakan ke error konflik dengan `current_version` bila ada.
- [x] Header `X-CSRF-Token` terkirim pada semua mutasi; token tidak pernah masuk localStorage.
- [x] Mutasi tidak di-retry otomatis; read di-retry terbatas.

**Verifikasi & bukti selesai:** hasil unit test frontend.

**Bukti selesai (29 September 2026):** `npm run typecheck` bersih; `npm run lint` 0 problem; `npm run test` **21/21 lulus** (vitest, +5 tes date TASK-005); `npm run build` sukses. File: `frontend/src/lib/api/{types,errors,csrf,client}.ts` + `client.test.ts` + `errors.test.ts`; `providers.tsx` QueryClient memakai `queryRetryPolicy` + `mutationRetryPolicy=false`. Bukti per kriteria: `errors.test.ts` assert 422→`fields` map (isValidation true) dan 409→`current_version`=7 (isConflict) + 401 isUnauthorized + 429 retryAfter=30 + body non-JSON→generik; `client.test.ts` assert GET tanpa header CSRF meski token ter-set, POST/PUT/PATCH/DELETE mengirim `X-CSRF-Token`, credentials same-origin, fetch-reject→NETWORK_ERROR status 0, 204→data undefined, dan **`storageWrites` kosong** (token tidak pernah ditulis ke localStorage); retry policy: read retry utk 5xx/429/408/network/timeout s.d. maks 2, TIDAK utk 404/422/403, mutasi=false. Catatan: hook auto-refetch `/auth/me` setelah reload menunggu endpoint tersedia (TASK-010); interface + store CSRF sudah siap dan teruji lewat mock.

---

### TASK-007 — CI pipeline + harness tes integrasi PostgreSQL

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M0 / P0 / needs_verification (29 September 2026) — harness, gerbang lokal, dan `go test -race` (via MinGW GCC) lulus; hanya run CI hijau di GitHub menunggu push ke remote |
| DEV | DEV-13 (bertahap sejak awal) |
| Referensi | TEST_PLAN §1, §5, §6; RULES §8; NFR-07 |
| Dependensi | TASK-002, TASK-003, TASK-005 |
| Estimasi | 4–6 jam |

**Tujuan:** setiap push menjalankan gate TEST_PLAN §6: gofmt/vet/test/build backend, race detector, lint/typecheck/test/build frontend, dan migrasi ke PostgreSQL test yang dapat di-reset.

**File rencana (usulan):** `.github/workflows/ci.yml` (atau CI pilihan lain — catat), `backend/internal/platform/testdb/` (helper spin-up DB test: migrasi otomatis, truncate antar suite), `frontend` vitest config.

**Langkah implementasi:**
1. Pilih dan catat platform CI (default usulan: GitHub Actions).
2. Job backend: setup Go versi terkunci → gofmt check → `go vet ./...` → `go test ./...` → `go test -race ./...` → build binary.
3. Service PostgreSQL pada CI; helper testdb menjalankan migrasi ke schema test unik per run (TEST_PLAN §6 "database terisolasi per run").
4. Job frontend: `npm ci` → lint → typecheck → unit test → production build.
5. Job dependency scan: `govulncheck` + `npm audit` (temuan dikaji, bukan auto-fail — SECURITY §10).
6. Badge/status check wajib pada branch utama; commit terfokus sesuai RULES §8.

**Acceptance criteria:**
- [ ] Push contoh memicu CI hijau untuk seluruh gate di atas. — **belum terverifikasi**: repo belum punya remote GitHub (`git remote -v` kosong), sehingga run CI nyata belum bisa dijalankan/dicatat.
- [x] Suite integrasi contoh (1 tes DB sederhana lewat helper testdb) lulus di lokal terhadap `randesk_test`. (Di CI: otomatis lewat service container.)
- [x] Race detector aktif dan lulus. — dibuktikan lokal: toolchain MinGW-w64 GCC 16.1.0 (WinLibs POSIX/UCRT via winget) dipasang; `CGO_ENABLED=1 go test -race ./...` **exit 0** (semua paket `ok`, termasuk `testdb` terhadap `randesk_test`). Skrip: `backend/scripts/run-race-task007.ps1`.
- [x] Tidak ada secret di konfigurasi CI (hanya kredensial service container sementara + `env:` injection; `actions/*` tanpa token).

**Verifikasi & bukti selesai:** URL/log run CI hijau, tangkapan daftar job.

**Bukti sejauh ini (29 September 2026):** Workflow dibuat di `.github/workflows/ci.yml` (GitHub Actions — keputusan D-02/§4.2; tanpa karakter tab; 3 job: backend + service `postgres:17`, frontend, security-scan non-blocking `govulncheck`+`npm audit`). Helper `backend/internal/platform/testdb` menerapkan migrasi dari FS tertanam (`backend/migrations` `//go:embed`) via golang-migrate library + `postgres.WithInstance` di atas pool pgx (tanpa memakai DATABASE_URL dev; skip bila `TEST_DATABASE_URL` kosong). Bukti lokal: `gofmt -l` kosong; `go vet ./...` bersih; `go build ./...` sukses; `go test ./...` semua `ok` (testdb SKIP tanpa env); **`TestIntegrationMigrateAndCRUD` PASS terhadap `randesk_test`**; dan **`go test -race ./...` exit 0** dengan MinGW GCC 16.1.0 (CGO aktif) termasuk paket `testdb`. Frontend gate (lint/typecheck/vitest 21 lulus/build) hijau. **Yang tersisa untuk `done`:** push ke remote GitHub + tautkan URL run CI hijau.

---

## M1 — Akun dan Sesi

### TASK-008 — Password hashing Argon2id + session store

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-04 |
| Referensi | SECURITY §2, §3; SCHEMA §3 (sessions, users); API_SPEC §4; BR-17; FR-01 |
| Dependensi | TASK-003, TASK-004 |
| Estimasi | 4–6 jam (area baru bagi pengembang — ketidakpastian lebih tinggi) |

**Tujuan:** modul `internal/auth` inti: hash password Argon2id (PHC), verifikasi, pembuatan sesi (token 32 byte base64url, DB menyimpan SHA-256), CSRF token 43 karakter per sesi, expiry absolut 8 jam, revocation.

**File rencana (usulan):** `backend/internal/auth/password.go`, `backend/internal/auth/session_repo.go`, `backend/internal/auth/session_service.go`, tes unit + integrasi.

**Langkah implementasi:**
1. Pilih library Argon2id terawat (catat pilihan); parameter SECURITY §2: memory 64 MiB, iterations 3, parallelism 1, salt ≥16 byte acak, output 32 byte, format PHC.
2. Password 12–128 karakter, tidak ditrim/dipotong; batas panjang dicek sebelum hashing (RULES §7); verifikasi memakai perbandingan library.
3. Session service: generate token acak crypto 32 byte → base64url; simpan `token_hash` SHA-256; csrf_token acak terpisah 32 byte → 43 char base64url (SCHEMA CHECK); TTL absolut dari SESSION_TTL_HOURS (default 8, tanpa sliding).
4. Fungsi revoke: sesi kini, seluruh sesi user (dipakai change-password/reset/deactivate — BR-17).
5. Login menggantikan sesi lama pada browser sama (API_SPEC §4) — revokasi token lama dari cookie yang dikirim bila ada.
6. Tes: hash/verify roundtrip; password dengan spasi & Unicode; token berbeda tiap login; revoke-all menghapus keaktifan; expiry dihormati (clock dapat dikontrol — TEST_PLAN §2).

**Acceptance criteria:**
- [ ] Hash tersimpan format PHC; plaintext tidak pernah persisten/log.
- [ ] Verifikasi salah → false tanpa panic; password 11 atau 129 karakter → error validasi.
- [ ] `sessions.token_hash` 64 hex, `csrf_token` 43 char — lolos CHECK DB.
- [ ] Revoke-all menandai `revoked_at` semua sesi user dalam satu transaksi.
- [ ] Unit test lulus termasuk kasus batas panjang dan Unicode.

**Verifikasi & bukti selesai:** hasil `go test ./internal/auth/...`.

---

### TASK-009 — Middleware autentikasi sesi + CSRF + validasi Origin

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-04 |
| Referensi | SECURITY §3, §4; ARCHITECTURE §5, §7; API_SPEC §1, §4; BR-17; NFR-01; TC-02, TC-03 |
| Dependensi | TASK-008 |
| Estimasi | 4–6 jam |

**Tujuan:** setiap request terproteksi: baca cookie, cek hash/expiry/revoked/is_active per request, actor masuk context; mutasi lolos validasi Origin + `X-CSRF-Token` (constant-time); gate `must_change_password` membatasi ke me/logout/change-password.

**File rencana (usulan):** `backend/internal/auth/middleware.go`, `backend/internal/auth/csrf.go`, `backend/internal/platform/httpx/cookies.go`, tes integrasi.

**Langkah implementasi:**
1. Cookie session: nama `__Host-randesk_session` (produksi HTTPS) / `randesk_session_dev` (development) dengan atribut SECURITY §3; helper set/clear.
2. Middleware auth: lookup `token_hash`, periksa `expires_at`, `revoked_at`, `users.is_active` setiap request (BR-17); gagal → 401 UNAUTHENTICATED generik.
3. Middleware CSRF untuk mutasi: validasi Origin terhadap APP_ORIGIN; fallback Referer; keduanya tidak ada/cocok → tolak (SECURITY §4). Bandingkan `X-CSRF-Token` dengan `sessions.csrf_token` constant-time. Login dikecualikan dari token sesi tetapi tetap wajib JSON + Origin sah + tolak form content-type.
4. Gate `must_change_password=true`: hanya `/auth/me`, `/auth/logout`, `/auth/change-password` lolos; lainnya 403 `PASSWORD_CHANGE_REQUIRED` (API_SPEC §3).
5. Actor (id, role, name) disuntikkan ke context untuk service; handler tidak mengakses gin.Context di service (ARCHITECTURE §4).
6. Tes integrasi (TC-02/03 sebagian): cookie dicabut → 401; akun nonaktif → 401; CSRF hilang/salah → 403 CSRF_FAILED tanpa perubahan data; Origin asing → ditolak; GET tidak mengubah state.

**Acceptance criteria:**
- [ ] Request tanpa/dengan cookie invalid → 401; respons tidak membocorkan alasan spesifik (expired vs dicabut boleh generik).
- [ ] Mutasi tanpa header CSRF atau Origin tidak sah → 403; data tidak berubah (dicek DB).
- [ ] Perbandingan token constant-time (bukti kode memakai `subtle.ConstantTimeCompare`).
- [ ] Akun `must_change_password` menerima 403 PASSWORD_CHANGE_REQUIRED pada endpoint lain.
- [ ] `Cache-Control: no-store` pada respons sensitif auth.

**Verifikasi & bukti selesai:** hasil tes integrasi auth; contoh request/response.

---

### TASK-010 — Endpoint autentikasi: login, me, logout, change-password

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-04 |
| Referensi | API_SPEC §4; SECURITY §2, §3; DESIGN D-01 (kontrak pesan); FR-01; TC-01, TC-02, TC-03, TC-28 |
| Dependensi | TASK-009 |
| Estimasi | 4–6 jam |

**Tujuan:** empat endpoint auth sesuai kontrak API_SPEC §4 beserta DTO user sesi dan audit `auth.password_changed`.

**File rencana (usulan):** `backend/internal/auth/handler.go`, `backend/internal/auth/dto.go`, `backend/internal/users/` (lookup user), tes API contract.

**Langkah implementasi:**
1. `POST /auth/login`: validasi email (normalisasi trim+lowercase) & password; gagal → pesan generik tunggal (tidak membedakan akun tidak ada/salah/nonaktif — DESIGN D-01); sukses → Set-Cookie + `{data:{user, csrf_token}}`.
2. `GET /auth/me`: user DTO (id, name, email, role, is_active, must_change_password, department{id,name}|null) + csrf_token; tanpa hash apapun.
3. `POST /auth/logout`: `{}` + CSRF → revoke sesi kini, hapus cookie, 204. Frontend memperlakukan logout dengan sesi hilang sebagai berhasil (API_SPEC §10).
4. `POST /auth/change-password`: verifikasi current_password, validasi new (12–128, tidak ditrim), update hash, `must_change_password=false`, **revoke seluruh sesi dalam transaksi yang sama**, hapus cookie, 204; tulis audit `auth.password_changed` (transaksional — bergantung TASK-013; bila 013 belum siap, stub writer dengan interface).
5. Tes API contract: TC-01 (benar/salah/nonaktif), TC-02 (logout/change-password mencabut; cookie lama tidak bisa dipakai), TC-28 sebagian (must_change_password membatasi endpoint).
6. Rate limit login disusul TASK-042; catat di task tersebut.

**Acceptance criteria:**
- [ ] Login sukses menghasilkan cookie HttpOnly + atribut lingkungan benar; respons memuat csrf_token.
- [ ] Kredensial salah / akun nonaktif → 401 dengan pesan generik identik.
- [ ] Setelah change-password: seluruh sesi user revoked (verifikasi DB) dan cookie dibersihkan; audit tercatat.
- [ ] DTO tidak pernah memuat password_hash/token_hash.
- [ ] TC-01/02 lulus sebagai tes otomatis.

**Verifikasi & bukti selesai:** hasil tes API contract auth; contoh Set-Cookie (nilai disamarkan).

---

### TASK-011 — cmd/admin: bootstrap admin pertama

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-04 (prasyarat DEV-05) |
| Referensi | OPERATIONS §3; SCHEMA §6 (`system.admin_bootstrapped`); SECURITY §2; ARCHITECTURE §3 (`backend/cmd/admin/`) |
| Dependensi | TASK-008 |
| Estimasi | 2–4 jam |

**Tujuan:** perintah operator interaktif membuat admin pertama dengan password dari prompt (bukan argumen CLI), hash via kode aplikasi, `must_change_password=true`, audit `actor_id NULL`.

**File rencana (usulan):** `backend/cmd/admin/main.go` (subcommand `bootstrap-admin`).

**Langkah implementasi:**
1. Prompt email, nama, department opsional, password (input tersembunyi); tolak password di argumen/history (OPERATIONS §3).
2. Validasi input RULES §7; hash Argon2id; insert user role=admin; audit `system.admin_bootstrapped` dengan actor_id NULL.
3. Idempoten-aman: bootstrap kedua kali ditolak bila sudah ada admin aktif (LAST_ADMIN logic arah sebaliknya), pesan jelas.
4. Tes integrasi: bootstrap pada DB kosong berhasil; password tidak muncul di log/output; run kedua ditolak.

**Acceptance criteria:**
- [ ] Admin pertama dibuat dengan `must_change_password=true`; login berikutnya dipaksa ganti password.
- [ ] Audit `system.admin_bootstrapped` tertulis dengan actor_id NULL.
- [ ] Tidak ada password pada process list (bukti: prompt interaktif), log, maupun audit metadata.

**Verifikasi & bukti selesai:** transkrip eksekusi bootstrap (password disamarkan), baris audit hasil query.

---

### TASK-012 — Frontend auth: halaman login, change-password, bootstrap sesi, route guard

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-04 (sisi UI) |
| Referensi | DESIGN D-01, §3 (guard), §8; ARCHITECTURE §7; API_SPEC §4; TC-02 (cache dibersihkan), TC-35 (refresh browser) |
| Dependensi | TASK-006, TASK-010 |
| Estimasi | 4–6 jam |

**Tujuan:** alur login → dashboard, paksa change-password, pemulihan sesi setelah reload via `/auth/me`, route guard peran, dan pembersihan cache TanStack Query + draft sensitif saat sesi berakhir.

**File rencana (usulan):** `frontend/src/features/auth/` (LoginPage, ChangePasswordPage, AuthProvider, useAuth), `frontend/src/app/guards.tsx`.

**Langkah implementasi:**
1. LoginPage sesuai D-01: field email/password, toggle tampil password, tombol loading, cegah double-submit, satu pesan error generik, teks bantuan "Hubungi admin untuk bantuan akun"; tanpa link registrasi/reset.
2. Sukses login → simpan csrf_token (memori, via lib TASK-006) → redirect dashboard; `must_change_password=true` → redirect `/change-password`.
3. ChangePasswordPage: password lama, baru, konfirmasi lokal; aturan panjang ditampilkan; sukses → logout lokal → login (sesi dicabut server).
4. AuthProvider: bootstrap via `GET /auth/me` saat mount/refresh (TC-35 prasyarat); state unauthenticated → semua route guard redirect `/login`; pengguna bersesi yang membuka `/login` → dashboard.
5. Guard peran untuk route admin (UI saja; backend tetap otoritatif — DESIGN §3).
6. Pada 401/logout: `queryClient.clear()` + reset CSRF memory (ARCHITECTURE §7, TC-02).
7. Tes komponen: render error generik, redirect must_change_password, guard role; aksesibilitas label form (DESIGN §8).

**Acceptance criteria:**
- [ ] Login benar → dashboard; salah → pesan generik D-01; tombol tidak bisa double-submit.
- [ ] Reload browser mempertahankan sesi (me) tanpa kehilangan state auth.
- [ ] Akun must_change_password hanya bisa mengakses change-password (UI + bukti 403 API).
- [ ] Logout mengosongkan cache query; halaman terproteksi kembali ke login.
- [ ] Form memiliki label eksplisit dan error terlihat sampai diperbaiki.

**Verifikasi & bukti selesai:** hasil tes komponen; rekaman/screenshot alur login 3 peran.

---

### TASK-013 — Audit infrastructure: writer transaksional + GET /audit-logs

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 |
| Referensi | FR-13; SCHEMA §6 (registry action), §5 (append-only); API_SPEC §5; RULES §2 (hanya admin membaca); TC-29; BR-18 |
| Dependensi | TASK-009 |
| Estimasi | 3–5 jam |

**Tujuan:** modul `internal/audit` yang menulis `audit_logs` **dalam transaksi pemanggil** (syarat TC-29: audit gagal → perubahan ikut rollback), registry action awal, dan endpoint baca khusus admin.

**File rencana (usulan):** `backend/internal/audit/writer.go`, `backend/internal/audit/handler.go`, `backend/internal/audit/repository.go`, tes integrasi.

**Langkah implementasi:**
1. Writer menerima `*sql.Tx` + actor + action (registry SCHEMA §6: user.created, user.updated, user.password_reset, auth.password_changed, category.*, department.*, system.admin_bootstrapped) + target + metadata allowlist + request_id dari context.
2. Metadata hanya nama field berubah & nilai role/status/departemen relevan — tanpa secret/payload akun penuh.
3. `GET /audit-logs` (admin): filter actor_id, action, from, to; pagination; scope admin-only (403 untuk non-admin).
4. Tes: rollback transaksi menggagalkan audit+perubahan bersama (pola TC-29 diuji di task pemakai, writer diuji di sini dengan tx gagal); non-admin → 403; filter/pagination benar.

**Acceptance criteria:**
- [ ] Audit tertulis atomik bersama mutasi pemanggil (bukti tes: insert audit error → tx rollback).
- [ ] Registry action sesuai SCHEMA §6; metadata tanpa password/token.
- [ ] Hanya admin dapat membaca; respons paginated sesuai envelope.
- [ ] Tidak ada endpoint update/delete audit (BR-18).

**Verifikasi & bukti selesai:** hasil tes integrasi audit.

---

### TASK-014 — User management backend: list, create, patch, reset-password

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 |
| Referensi | FR-02; API_SPEC §5; RULES §2, §7 (validasi), BR-15 (role immutable — bagian create/patch); SCHEMA §3 (users); TC-28; DESIGN D-07 (kontrak data) |
| Dependensi | TASK-013 |
| Estimasi | 4–6 jam |

**Tujuan:** endpoint admin `/users` (GET list + filter, POST create dengan temporary_password sekali tampil, PATCH name/department_id/is_active, POST reset-password) dengan audit transaksional.

**File rencana (usulan):** `backend/internal/users/handler.go`, `service.go`, `repository.go`, `dto.go`, tes API contract.

**Langkah implementasi:**
1. DTO allowlist: role/email tidak dapat diubah via PATCH (BR-15 — enforcement penuh aktivasi di TASK-015); tolak field asing.
2. Create: validasi RULES §7 (nama 2–100, email ≤254 normalisasi lowercase unik, dst.); generate temporary_password acak ≥16 karakter server-side; `must_change_password=true`; hash Argon2id; audit `user.created` transaksional; password hanya di respons ini, tidak dilog.
3. Reset-password: bukan untuk diri sendiri (403); revoke semua sesi user target (BR-17) dalam transaksi; audit `user.password_reset`; temporary_password sekali tampil.
4. List: filter q/role/is_active/department_id + pagination; EMAIL_exists → 409 EMAIL_EXISTS; NAME tidak unik case-insensitive → 409 NAME_EXISTS (SCHEMA unique index lower).
5. `GET /assignees` (staff): teknisi aktif saja, field id+name (API_SPEC §5) — bisa di task ini atau TASK-016; pilih di sini karena query users.
6. Tes: create/list/patch/reset happy path + 403 non-admin + 409 email duplikat + temporary password tidak muncul di log/audit + reset mencabut sesi (TC-28 sebagian).

**Acceptance criteria:**
- [ ] Admin membuat akun → 201 dengan temporary_password; respons berikutnya tidak dapat membaca ulang password.
- [ ] PATCH role/email ditolak (field ditolak sebagai unknown/immutable).
- [ ] Reset password diri sendiri → 403; reset orang lain mencabut semua sesi target (bukti DB).
- [ ] `GET /assignees` hanya mengembalikan teknisi aktif; requester → 403.
- [ ] Audit user.created/user.updated/user.password_reset commit bersama perubahan.

**Verifikasi & bukti selesai:** hasil tes API contract users.

---

### TASK-015 — Guard aktivasi akun: self-deactivation, LAST_ADMIN, ACTIVE_TICKETS_EXIST, advisory lock

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 |
| Referensi | BR-15, BR-17; ARCHITECTURE §6.2 (urutan lock, advisory lock); RULES §5; API_SPEC §3 (409 ACTIVE_TICKETS_EXIST/LAST_ADMIN, 403); TC-24, TC-25; SCHEMA §5 (invariant) |
| Dependensi | TASK-014 |
| Estimasi | 4–6 jam |

**Tujuan:** aturan penonaktifan/aktivasi akun yang aman terhadap konkurensi: tidak kehilangan admin aktif terakhir, tidak menonaktifkan teknisi bertiket aktif, tidak menonaktifkan diri sendiri, serialisasi via advisory transaction lock.

**File rencana (usulan):** `backend/internal/users/activation.go` (+ repository queries), tes konkurensi integrasi.

**Langkah implementasi:**
1. Urutan lock ARCHITECTURE §6.2: advisory lock administrasi (key tetap per instalasi) → lock akun target (FOR UPDATE) → cek referensi; **jangan** ambil ticket lock pada jalur deaktivasi.
2. Aturan: deactivate diri sendiri → 403 FORBIDDEN; deactivate admin yang merupakan admin aktif terakhir → 409 LAST_ADMIN (hitung ulang dalam transaksi di bawah advisory lock); deactivate teknisi dengan tiket open/in_progress → 409 ACTIVE_TICKETS_EXIST.
3. Setiap deaktivasi/reset mencabut semua sesi (sudah dari TASK-014 — pastikan dalam transaksi yang sama).
4. Tes konkurensi (prasyarat TC-25, diselesaikan di TASK-045): dua koneksi menonaktifkan dua admin berbeda secara bersamaan → tepat satu berhasil, admin aktif tersisa ≥1; deactivate + assignment bersamaan (pola TC-24 sisi akun; sisi tiket dilengkapi TASK-025/045).

**Acceptance criteria:**
- [ ] Self-deactivate → 403; last-admin → 409 LAST_ADMIN; teknisi bertiket aktif → 409 ACTIVE_TICKETS_EXIST dengan detail yang tidak membocorkan objek di luar scope.
- [ ] Dua deaktivasi admin konkuren tidak pernah menghasilkan 0 admin aktif (tes 2 koneksi, diulang ≥20x).
- [ ] Sesi user nonaktif langsung tidak valid pada request berikutnya (BR-17).

**Verifikasi & bukti selesai:** hasil tes konkurensi (log run berulang), query verifikasi jumlah admin aktif.

---

### TASK-016 — Master data backend: categories, departments

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 |
| Referensi | FR-03; BR-14; API_SPEC §5; RULES §2, §7; ARCHITECTURE §6.2 (shared row lock master); TC-26; DESIGN D-07 |
| Dependensi | TASK-013 |
| Estimasi | 4–6 jam |

**Tujuan:** CRUD kategori/departemen (create admin, patch admin termasuk is_active, list dengan include_inactive/q/pagination), nama unik case-insensitive, deactivation bukan delete, lock master untuk konsistensi referensi.

**File rencana (usulan):** `backend/internal/masters/` (atau `internal/users/` sesuai struktur — putuskan dan catat), tes API contract.

**Langkah implementasi:**
1. GET `/categories` semua peran (untuk dropdown tiket nantinya); GET `/departments` admin; keduanya include_inactive default false, paginated — daftar aktif tetap dipaginasi (API_SPEC §5).
2. POST/PATCH admin-only; validasi nama 2–80 unik case-insensitive (409 NAME_EXISTS); description ≤500 opsional/nullable.
3. PATCH is_active=false → deactivation (BR-14): referensi historis tetap valid; tidak ada DELETE endpoint (BR-18).
4. Mutasi status aktif mengambil row lock master yang berkonflik dengan pemilihan referensi baru (ARCHITECTURE §6.2) — implementasikan shared lock helper yang akan dipakai TASK-023 (pemilihan kategori saat create ticket).
5. Audit category.*/department.* transaksional via TASK-013.
6. Tes: happy path; duplikat nama beda kapitalisasi → 409; non-admin create → 403; item nonaktif tetap terbaca dengan include_inactive=true (TC-26 sebagian — sisi pemilihan referensi diuji di TASK-023).

**Acceptance criteria:**
- [ ] Kategori/departemen nonaktif tidak hilang dan referensi lama tetap resolve.
- [ ] Nama duplikat case-insensitive ditolak 409.
- [ ] Hanya admin dapat mutate; audit tertulis transaksional.
- [ ] List paginated dengan filter q dan include_inactive sesuai kontrak.

**Verifikasi & bukti selesai:** hasil tes API contract master data.

---

### TASK-017 — Seed sintetis development + fixture tes minimum

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-02 |
| Referensi | SCHEMA §8 (seed); TEST_PLAN §2 (fixture minimum); PRD §4.2 (demo dari seed); SECURITY §9 (data sintetis, domain reserved) |
| Dependensi | TASK-003, TASK-016 |
| Estimasi | 3–5 jam |

**Tujuan:** perintah seed dev (2 requester, 2 technician, 2 admin, akun nonaktif, kategori nonaktif, departemen, tiket lintas status termasuk reopened, komentar internal dengan penanda unik, lampiran sintetis) dan builder fixture yang dipakai ulang suite tes.

**File rencana (usulan):** `backend/cmd/admin/` subcommand `seed-dev`, `backend/internal/platform/testdb/fixtures.go`, dataset terpisah dari migrasi (SCHEMA §8).

**Langkah implementasi:**
1. Seed memakai email domain `example.test`, nama sintetis (SECURITY §9); password dari env dev atau acak + hash helper aplikasi; ticket_no dari sequence, bukan insert manual.
2. Isi fixture TEST_PLAN §2: R1/R2 masing-masing ≥2 tiket; satu tiket dengan komentar internal berpembaca-unique-marker (untuk deteksi kebocoran TC-15); tiket open kosong, open assigned, in_progress, resolved, closed, reopened.
3. File lampiran sintetis valid (PNG/PDF kecil) + file uji tidak valid (MIME palsu, oversize) dibuat generator, bukan binary di repo.
4. Clock terkontrol untuk waktu bisnis (first_response_at/resolved_at fixture) sesuai TEST_PLAN §2.
5. Pisahkan seed dev, seed test, dan (nanti) benchmark dataset TASK-048.

**Acceptance criteria:**
- [ ] `seed-dev` pada DB kosong menghasilkan fixture lengkap dan dapat di-reset.
- [ ] Login dengan akun seed bekerja; tidak ada password seragam yang hardcoded di repo (nilai dari env/generated).
- [ ] Fixture tes menyediakan seluruh entitas TEST_PLAN §2 dan dipakai minimal oleh satu suite yang berjalan.

**Verifikasi & bukti selesai:** output perintah seed, query jumlah row per tabel.

---

### TASK-018 — Admin UI: manajemen akun

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 (sisi UI) |
| Referensi | DESIGN D-07; FR-02; API_SPEC §5; TC-28 (secret sekali tampil, tidak disimpan browser) |
| Dependensi | TASK-012, TASK-015 |
| Estimasi | 4–6 jam |

**Tujuan:** layar `/admin/users`: tabel akun, form buat akun, edit nama/departemen, aktivasi/nonaktifkan, reset password dengan tampilan temporary_password sekali.

**File rencana (usulan):** `frontend/src/features/admin/users/`, komponen bersama `frontend/src/components/` (ConfirmDialog, FormField, ErrorSummary, Pagination).

**Langkah implementasi:**
1. Tabel nama/email/role/departemen/status + filter (q, role, is_active, department) + pagination server; role read-only setelah dibuat.
2. Dialog buat akun → tampilkan temporary_password sekali dengan aksi salin; tidak disimpan di state setelah ditutup (TC-28).
3. Reset password dengan ConfirmDialog yang menjelaskan akibat (sesi dicabut, user wajib ganti password); hasil password sekali tampil.
4. Aktivasi/nonaktifkan dengan ConfirmDialog; error 409 (LAST_ADMIN, ACTIVE_TICKETS_EXIST) dan 403 self-deactivate ditampilkan dengan pesan bermakna dari error mapping TASK-006.
5. State loading/empty/error/success (DESIGN §8); label aksesibel; tabel boleh horizontal scroll di mobile.
6. Tes komponen: render tabel, alur create → dialog password, mapping error 409.

**Acceptance criteria:**
- [ ] CRUD akun admin berfungsi end-to-end terhadap API M1.
- [ ] Temporary password hanya tampil sekali dan hilang dari memori UI setelah dialog ditutup.
- [ ] Error 409/403 deaktivasi tampil sebagai pesan yang dapat dipahami, bukan toast generik saja.
- [ ] Non-admin tidak melihat route/menu admin (guard TASK-012).

**Verifikasi & bukti selesai:** hasil tes komponen; screenshot layar pada 1024 px dan 360 px.

---

### TASK-019 — Admin UI: master data + audit log

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M1 / P0 / todo |
| DEV | DEV-05 (sisi UI) |
| Referensi | DESIGN D-07; FR-03, FR-13; API_SPEC §5; BR-14 (tanpa tombol hapus) |
| Dependensi | TASK-012, TASK-013, TASK-016 |
| Estimasi | 3–5 jam |

**Tujuan:** layar `/admin/categories`, `/admin/departments`, `/admin/audit`.

**File rencana (usulan):** `frontend/src/features/admin/masters/`, `frontend/src/features/admin/audit/`.

**Langkah implementasi:**
1. Tabel kategori/departemen: nama, deskripsi opsional, status aktif; form create/edit; aksi "Nonaktifkan" dengan penjelasan data historis tetap ada; **tanpa tombol hapus** (BR-14/18).
2. Layar audit: tabel waktu/aktor/aksi/target/request_id, filter actor/action/from/to + pagination server; format waktu WITA.
3. Dropdown master mendukung pencarian/pemuatan halaman berikut (API_SPEC §5) — komponen reusable SelectAsync.
4. State loading/empty/error; error NAME_EXISTS 409 tampil inline.
5. Tes komponen ringan untuk form + filter.

**Acceptance criteria:**
- [ ] Kategori dapat dibuat, diubah, dinonaktifkan; item nonaktif tampil dengan label dan tetap pada data historis.
- [ ] Audit hanya dapat diakses admin (guard + bukti 403 API bila dipaksa).
- [ ] Tidak ada UI delete untuk master data maupun audit.

**Verifikasi & bukti selesai:** hasil tes komponen; screenshot.

---

## M2 — Tiket Inti

### TASK-020 — Policy tiket terpusat + allowed_actions

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-06 |
| Referensi | RULES §2 (matriks izin), §3 (lifecycle); ARCHITECTURE §4 (policy pusat); API_SPEC §2 (allowed_actions); NFR-01; BR-16; TC-04, TC-10 |
| Dependensi | TASK-009 |
| Estimasi | 4–6 jam |

**Tujuan:** satu modul policy yang menjawab: can_read, can_edit_metadata, can_change_priority, can_comment_public/internal, can_assign/claim/unassign, can_transition(start/resolve/close/reopen), can_upload/delete_attachment — dipakai detail, list scope, dan (nanti) dashboard & download. Plus kalkulasi `allowed_actions` DTO.

**File rencana (usulan):** `backend/internal/tickets/policy.go`, `backend/internal/tickets/policy_test.go`.

**Langkah implementasi:**
1. Model input policy: actor (id, role, is_active), ticket snapshot (requester_id, assignee_id, status); untuk attachment: uploader.
2. Implementasikan seluruh baris matriks RULES §2 sebagai fungsi murni (tanpa DB) — table-driven test mencakup **setiap pasangan role × status × kepemilikan**, termasuk: teknisi non-assignee read-only; requester hanya miliknya; admin semua kecuali closed untuk mutasi; close/reopen oleh reporter atau admin (teknisi hanya bila reporter sendiri).
3. `allowed_actions` dihitung untuk actor kini sesuai nilai enum API_SPEC §2; `can_delete` attachment terpisah.
4. Scope predicate untuk query daftar: requester → `requester_id = actor`; staff → semua (dipakai TASK-022) — policy yang sama, bukan salinan logika (ARCHITECTURE §4).
5. Tes unit table-driven lengkap; sertakan kasus TC-10 (izin mutasi mengikuti assignee terbaru setelah reassignment).

**Acceptance criteria:**
- [ ] Setiap sel matriks RULES §2 memiliki minimal satu kasus tes yang sesuai.
- [ ] Fungsi policy murni (deterministik, tanpa I/O) dan diekspor untuk dipakai modul lain.
- [ ] allowed_actions untuk requester pada tiket orang lain → kosong/tidak dapat dibaca (read → false).
- [ ] Tidak ada duplikasi logika izin di tempat lain (bukti: pemanggil hanya melalui policy).

**Verifikasi & bukti selesai:** hasil `go test ./internal/tickets/ -run Policy` dengan daftar kasus table-driven.

---

### TASK-021 — Event & notification writer transaksional (pemilihan penerima, dedup)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-07 |
| Referensi | BR-12, BR-10; RULES §6 (penerima, dedup, isi aman); SCHEMA §3, §6 (payload allowlist); ARCHITECTURE §6.1, §8; NFR-02; FR-10, FR-11 (sisi tulis); TC-07, TC-20 |
| Dependensi | TASK-003, TASK-017 (fixture untuk tes penerima) |
| Estimasi | 4–6 jam (ketidakpastian tinggi — pola transaksi Go baru bagi pengembang) |

**Tujuan:** modul penulis aktivitas: dalam `*sql.Tx` yang sama menulis `ticket_events` (payload allowlist, visibility) dan `notifications` untuk seluruh penerima sesuai RULES §6 — fondasi semua mutasi tiket.

**File rencana (usulan):** `backend/internal/notifications/writer.go`, `backend/internal/tickets/events.go`, tes integrasi.

**Langkah implementasi:**
1. Event writer: insert `ticket_events` (event_type enum SCHEMA, seq identity, visibility, payload JSONB **hanya field allowlist SCHEMA §6** — tanpa body komentar, token, storage_key, password).
2. Recipient resolver per event_type mengikuti tabel RULES §6: tiket baru → admin aktif; assignment → assignee baru + reporter; reassign/unassign → lama + baru + reporter; status berubah → reporter + assignee kini (+ assignee lama untuk reopen); komentar publik → reporter + assignee (admin aktif bila tanpa assignee); komentar internal → assignee + admin aktif; prioritas/kategori/judul/deskripsi → reporter + assignee; lampiran → reporter + assignee.
3. Aturan: hilangkan duplikat, akun nonaktif, dan aktor pemicu; unique `(recipient_id, event_id)`; message ringkas aman ≤200 karakter, tanpa isi catatan internal (RULES §6).
4. Semua insert menerima `*sql.Tx` pemanggil — tidak ada koneksi/goroutine terpisah (ARCHITECTURE §8).
5. Tes integrasi: TC-20 (penerima benar per event, tanpa duplikat, tanpa internal ke requester, aktor dikecualikan, akun nonaktif dilewati); pola rollback TC-07 (event insert gagal → mutasi batal — diuji penuh di TASK-023); unique constraint terbukti menahan duplikasi.

**Acceptance criteria:**
- [ ] Setiap event_type RULES §6 memiliki tes penerima yang lulus.
- [ ] Payload event hanya berisi field allowlist SCHEMA §6 (dibuktikan tes yang menolak field ekstra saat build payload).
- [ ] Notifikasi tidak memuat teks catatan internal/lampiran; message ≤200 char.
- [ ] Insert event/notifikasi gagal → seluruh tx rollback (tidak ada data setengah jadi).
- [ ] `(recipient_id, event_id)` unik ditegakkan (uji duplikasi → tidak menghasilkan dua baris).

**Verifikasi & bukti selesai:** hasil tes integrasi writer; query contoh isi notifications untuk event internal.

---

### TASK-022 — Ticket list & search backend (scope-first, filter, pagination, sort allowlist)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-06 |
| Referensi | FR-05; BR-16; RULES §5; API_SPEC §6 (query params, sort, DTO summary); SCHEMA §7 (ILIKE escape, LIMIT 100); TC-04, TC-15, TC-27; SECURITY §5 |
| Dependensi | TASK-016, TASK-020 |
| Estimasi | 4–6 jam |

**Tujuan:** `GET /tickets` dengan scope akses diterapkan **sebelum** filter/count/pagination; `GET /tickets/:id` detail; filter status, priority, category_id, from/to (created_at; from inklusif/to eksklusif; wajib berpasangan; from<to), q (≤100 char; nomor via parse `HD-`, judul, deskripsi; escape `%`/`_`), staff-only requester_id/assignee_id (UUID | `me` | `unassigned`); sort allowlist + tie-breaker id; DTO summary tanpa description/resolution_summary.

**File rencana (usulan):** `backend/internal/tickets/list_handler.go`, `list_service.go`, `list_repository.go`, `dto.go`, tes API contract + integrasi.

**Langkah implementasi:**
1. Validasi parameter: requester mengirim filter staff → 403 (bukan scope diperluas); sort di luar allowlist → ditolak; per_page max 100; page ≥1.
2. Query builder dengan predicate scope dari policy (TASK-020); parameterized SQL; identifier sort dari allowlist internal (RULES §8).
3. Search: parse `HD-000042` → ticket_no; selain itu ILIKE title/description dengan escape wildcard; tidak mencari komentar (SCHEMA §7 — mencegah bocornya catatan internal).
4. Total dihitung setelah scope+filter; ORDER BY created_at + id sesuai arah sort.
5. `GET /tickets/:id`: policy can_read → 404 bila di luar scope (TC-04: keberadaan objek tidak bocor); DTO detail lengkap + allowed_actions + version; `Cache-Control: no-store`.
6. Tes: R1 tidak pernah melihat tiket R2 di list/search/count (TC-04); penanda internal tidak muncul di q search (bagian TC-15); sort terlarang & injeksi SQL pada q/filter ditolak aman (TC-27); from/to edge (from==to → 422; hanya from → 422).

**Acceptance criteria:**
- [ ] Scope requester membatasi list, total, dan search secara identik.
- [ ] Objek di luar scope baca → 404 (bukan 403) untuk detail.
- [ ] Filter staff dari requester → 403; sort/q injeksi tidak mengubah query (tes dengan input `'`, `%`, `; DROP`).
- [ ] Pagination default 20, max 100; meta.total akurat pasca-scope.
- [ ] DTO summary/detail sesuai API_SPEC §6 (ticket_code string; BIGINT tidak dikirim sebagai JS number).

**Verifikasi & bukti selesai:** hasil tes API contract list/search/detail.

---

### TASK-023 — Ticket create backend (POST /tickets)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-07 |
| Referensi | FR-04; BR-01, BR-02, BR-12; API_SPEC §6 (buat tiket); RULES §7 (batas judul/deskripsi); SCHEMA §1 (ticket_no), §6 (payload ticket.created); TC-05, TC-06, TC-07; ARCHITECTURE §6.2 (shared lock master); TC-26 (sisi pemilihan kategori) |
| Dependensi | TASK-021, TASK-022 |
| Estimasi | 3–5 jam |

**Tujuan:** pembuatan tiket transaksional: validasi, ticket_no dari identity, reporter dari sesi, event `ticket.created` + notifikasi admin aktif dalam satu tx, kategori harus aktif dengan shared lock master.

**File rencana (usulan):** `backend/internal/tickets/create_service.go`, `create_handler.go`, tes API contract + integrasi.

**Langkah implementasi:**
1. Input allowlist: title (5–150), description (20–10.000), category_id, priority opsional default normal; **tolak** requester_id/assignee_id/status/ticket_no/timestamps dari klien (TC-05 — decode DisallowUnknownFields → field asing 400).
2. Panjang dihitung karakter Unicode setelah trim (RULES §7); kategori dicek aktif dengan shared row lock (helper TASK-016) → nonaktif = 422 MASTER_INACTIVE.
3. Transaksi: insert ticket (version 1, status open) → event ticket.created (payload: priority, category_id) → notifikasi semua admin aktif → commit. Reporter = actor sesi (BR-01).
4. ticket_no: identity DB; format respons `ticket_code` = `HD-` + minimal 6 digit; angka besar tidak dipotong (SCHEMA §1).
5. Failure injection (TC-07): buat insert event/notifikasi gagal (test hook) → tiket tidak ada, tidak ada notifikasi; sequence boleh melompat.
6. Tes: validasi batas (judul 4/5 char, deskripsi 19/20, Unicode), kategori nonaktif, field asing ditolak, rollback TC-07, 201 + DTO detail.

**Acceptance criteria:**
- [ ] Tiket dibuat dengan reporter dari sesi; klien tidak dapat memalsukan reporter/status (TC-05 lulus).
- [ ] Batas validasi RULES §7 ditegakkan dengan pesan field yang sesuai (TC-06).
- [ ] ticket_code berformat benar dan unik; version=1; status open.
- [ ] Kegagalan tulis event/notifikasi → tidak ada tiket setengah jadi (TC-07 lulus).
- [ ] Kategori nonaktif → 422 MASTER_INACTIVE.

**Verifikasi & bukti selesai:** hasil tes API contract + failure injection.

---

### TASK-024 — Edit metadata backend (PATCH /tickets/:id)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-07 |
| Referensi | BR-05, BR-06, BR-08, BR-14; RULES §2 (siapa boleh edit), §5 (urutan pemeriksaan); API_SPEC §6 (edit metadata); ARCHITECTURE §6.1; TC-09, TC-14 |
| Dependensi | TASK-023 |
| Estimasi | 3–5 jam |

**Tujuan:** PATCH title/description/category_id (+priority admin) dengan expected_version, row lock, NO_CHANGE, reason wajib untuk priority, event `ticket.updated` + notifikasi transaksional.

**File rencana (usulan):** `backend/internal/tickets/update_service.go`, `update_handler.go`, tes.

**Langkah implementasi:**
1. Urutan pemeriksaan RULES §5: auth → scope objek (di luar scope → 404) → `SELECT … FOR UPDATE` → validasi role/akun → expected_version (mismatch → 409 VERSION_CONFLICT, boleh sertakan current_version bila actor boleh akses) → status/aturan → perubahan → event/notifikasi → commit.
2. UPDATE inti tetap `WHERE id=$id AND version=$expected` + cek rows affected sebagai pertahanan ganda (ARCHITECTURE §6.1); version +1; updated_at kini.
3. Aturan field: requester pemilik hanya saat open; admin semua kecuali closed; technician tidak dapat edit (RULES §2). Minimal satu field mutable wajib; seluruh nilai identik → 422 NO_CHANGE (BR-06). Priority hanya admin, reason 10–1.000 wajib bila priority berubah (BR-08); kategori baru harus aktif (shared lock master).
4. Closed → mutasi ditolak 409 INVALID_STATE (TC-14).
5. Event ticket.updated: changed_fields + old/new priority/category + reason; notifikasi ke reporter + assignee kini.
6. Tes: stale version 409 (TC-09 — perubahan pemenang tetap utuh); NO_CHANGE 422; priority tanpa reason 422; requester edit in_progress → 403; edit pada closed → 409; payload event benar.

**Acceptance criteria:**
- [ ] Hanya edit yang diizinkan matriks RULES §2 yang lolos; lainnya 403/409 sesuai kontrak.
- [ ] Version stale → 409 tanpa perubahan data; tidak ada last-write-wins.
- [ ] NO_CHANGE terdeteksi setelah normalisasi trim.
- [ ] Event + notifikasi commit atomik dengan perubahan; payload hanya allowlist.

**Verifikasi & bukti selesai:** hasil tes API contract edit + konkurensi dua klien (stale version).

---

### TASK-025 — Assignment backend (claim, assign, reassign, unassign)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-07 |
| Referensi | FR-06; BR-03, BR-04, BR-05; RULES §5 (urutan lock akun→tiket); ARCHITECTURE §6.2; API_SPEC §6 (assignment); PRD §7 (acceptance scenario FR-06); TC-08, TC-10, TC-24 |
| Dependensi | TASK-015, TASK-023 |
| Estimasi | 4–6 jam (ketidakpastian tinggi — konkurensi) |

**Tujuan:** `PUT /tickets/:id/assignee`: technician claim (open, assignee kosong, ID diri sendiri); admin initial assign/reassign (open/in_progress, reason 10–1.000 saat mengganti/melepas) / unassign (hanya open, null); satu pemenang saat claim bersaing; assignee selalu teknisi aktif.

**File rencana (usulan):** `backend/internal/tickets/assign_service.go`, `assign_handler.go`, tes konkurensi.

**Langkah implementasi:**
1. Lock order ARCHITECTURE §6.2: lock akun teknisi target (FOR UPDATE, urut UUID bila lebih dari satu) → cek role technician + is_active → lock tiket → cek expected_version + state → update.
2. Claim: status harus open, assignee NULL; expected_version wajib; kalah → 409 VERSION_CONFLICT atau ALREADY_ASSIGNED sesuai pemeriksaan yang gagal pertama (API_SPEC §6).
3. Assignment tidak mengubah status (BR-03); reassign/unassign oleh admin dengan reason (BR-04); unassign hanya open.
4. Event ticket.assigned/reassigned/unassigned (payload old/new assignee, reason) + notifikasi sesuai RULES §6, transaksional.
5. Tes: skenario PRD §7 FR-06 / TC-08 (dua teknisi expected_version sama → tepat satu 200, satu 409; satu assignee; satu event; notifikasi konsisten dengan pemenang); TC-10 (setelah reassign, teknisi lama kehilangan izin mutasi — bersama policy); TC-24 pola claim vs deaktivasi teknisi bersamaan (diperluas di TASK-045); assign ke akun nonaktif/non-teknisi → ditolak.

**Acceptance criteria:**
- [ ] Claim bersaing: tepat satu pemenang (tes 2 koneksi, diulang ≥20x, verifikasi state akhir DB — bukan hanya kode respons).
- [ ] Assignee non-teknisi atau nonaktif tidak pernah tersimpan.
- [ ] Reason wajib pada reassign/unassign; initial assign/claim tanpa reason diterima.
- [ ] Unassign pada in_progress → 409 INVALID_STATE.
- [ ] Event + notifikasi assignment atomik; tidak ada notifikasi untuk pihak yang kalah.

**Verifikasi & bukti selesai:** hasil tes konkurensi (log run berulang), query verifikasi satu assignee/event.

---

### TASK-026 — Transitions backend (start, resolve, close, reopen)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-07 |
| Referensi | FR-07; BR-07, BR-05; RULES §3 (tabel lifecycle); API_SPEC §6 (transitions); SCHEMA §4 (tickets_lifecycle_ck); PRD §9 (first response); TC-11, TC-12, TC-13, TC-14 |
| Dependensi | TASK-025 |
| Estimasi | 4–6 jam |

**Tujuan:** `POST /tickets/:id/transitions` dengan action start/resolve/close/reopen — hanya empat transisi RULES §3; set field turunan (first_response_at, resolved_at/by, closed_at/by, reopen_count); event + notifikasi transaksional.

**File rencana (usulan):** `backend/internal/tickets/transition_service.go`, `transition_handler.go`, tes.

**Langkah implementasi:**
1. Validasi input per action (API_SPEC §6): resolve → resolution_summary 20–4.000; reopen → reason 10–1.000; field tambahan tidak relevan → 400. Tidak ada PATCH status umum, direct close dari open, atau reopen dari closed.
2. start: open→in_progress oleh assignee/admin; assignee harus terisi & teknisi aktif; set first_response_at bila masih NULL dan aksi memenuhi definisi PRD §9 (aksi publik staff selain reporter).
3. resolve: in_progress→resolved oleh assignee/admin; set resolved_at/by + resolution_summary; solusi dicatat sebagai event publik (RULES §3).
4. close: resolved→closed oleh reporter/admin; set closed_at/by; closed terminal & read-only (TC-14 — mutasi isi/komentar/lampiran ditolak, diuji lintas modul di TASK-031/034).
5. reopen: resolved→open oleh reporter/admin; reason wajib; **kosongkan assignee + seluruh field resolusi/penutupan**; reopen_count+1; created_at & first_response_at tidak berubah; solusi lama tetap di event (TC-13).
6. Seluruh transisi: lock tiket, expected_version, CHECK DB tickets_lifecycle_ck sebagai jaring pengaman, event + notifikasi atomik (status berubah → reporter + assignee kini; reopen juga assignee lama).
7. Tes: matriks penuh status × action × actor (TC-11 — hanya 4 transisi sah lolos, lainnya 409 INVALID_STATE); TC-12 (resolve tanpa/short summary → 422, status tetap); TC-13 (state pasca-reopen diverifikasi kolom per kolom).

**Acceptance criteria:**
- [ ] Matriks TC-11 lulus sebagai table-driven integration test.
- [ ] Resolve tanpa ringkasan valid → 422 dan status tetap in_progress.
- [ ] Reopen menghasilkan: status open, assignee NULL, resolved_*/closed_* NULL, reopen_count+1, first_response_at & created_at tidak berubah, event lama tetap ada.
- [ ] Closed menolak transisi lanjutan (409).
- [ ] CHECK tickets_lifecycle_ck tidak pernah tertembus (bukti: tidak ada error constraint saat tes sah).

**Verifikasi & bukti selesai:** hasil tes transisi + query state akhir tiket reopened.

---

### TASK-027 — Runtime DB grants: enforcement append-only komentar/event/audit

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-02 (lanjutan) |
| Referensi | SCHEMA §5 (grants per tabel, runtime vs migration role); SECURITY §5; BR-18 |
| Dependensi | TASK-003, TASK-021 |
| Estimasi | 2–3 jam |

**Tujuan:** runtime DB role tidak memiliki UPDATE/DELETE pada ticket_comments, ticket_events, audit_logs; tabel mutable mendapat grant sesuai kebutuhan; terdokumentasi sebagai migrasi/script grants.

**File rencana (usulan):** `backend/migrations/00000X_runtime_grants.sql` (atau script operator — putuskan dan catat; grants bukan schema), dokumentasi setup di AGENTS.md.

**Langkah implementasi:**
1. Buat role runtime terpisah dari migration owner (OPERATIONS §3).
2. GRANT SELECT/INSERT pada tabel append-only; SELECT/INSERT/UPDATE terbatas pada tabel mutable; verifikasi privilege sequence/identity untuk insert (catatan SCHEMA §5).
3. Arahkan koneksi API dev/test ke role runtime; jalankan seluruh suite yang ada — pastikan tidak ada fitur sah yang rusak.
4. Tes bukti: koneksi runtime mencoba UPDATE ticket_comments → permission denied.

**Acceptance criteria:**
- [ ] Runtime role ditolak saat UPDATE/DELETE komentar/event/audit (bukti eksekusi SQL).
- [ ] Seluruh tes API yang ada tetap lulus dengan runtime role.
- [ ] Script grants terdokumentasi dan idempoten.

**Verifikasi & bukti selesai:** transkrip psql percobaan UPDATE ditolak + CI hijau.

---

### TASK-028 — UI daftar tiket (filter, search, pagination, URL state)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-08 |
| Referensi | DESIGN D-03, §5, §8; API_SPEC §6; FR-05; TC-35 (bagian) |
| Dependensi | TASK-012, TASK-022 |
| Estimasi | 4–6 jam |

**Tujuan:** halaman `/tickets` sesuai D-03: toolbar search+filter, tabel desktop/kartu mobile, pagination, filter tersimpan di URL, pembedaan empty vs no-results.

**File rencana (usulan):** `frontend/src/features/tickets/TicketsPage.tsx`, komponen `TicketFilters`, `TicketTable`, `TicketCard`, `TicketStatusBadge`, `PriorityBadge`, `Pagination` (DESIGN §7).

**Langkah implementasi:**
1. Query params ↔ URL state (status, priority, category, q, assignee, page, per_page, sort) agar navigasi kembali konsisten.
2. Search debounce ±300 ms + pembatalan request lama (TanStack Query + AbortController).
3. Kolom desktop & label status/prioritas Bahasa Indonesia sesuai D-03 dan DESIGN §2 prinsip 2; badge dengan label (tidak warna saja); waktu WITA + waktu absolut pada title/tooltip.
4. Filter staff (requester_id/assignee_id, shortcut "Ditugaskan kepada saya"/"Belum ditugaskan") hanya dirender untuk staff.
5. Mobile <768 px: daftar menjadi kartu; tabel boleh horizontal scroll.
6. State: loading skeleton, "Belum ada tiket" vs "Tidak ada hasil yang sesuai filter", error dengan retry.
7. Tes komponen: render filter, debounce, mapping empty state, sort link.

**Acceptance criteria:**
- [ ] Filter/pagination/sort bekerja terhadap API dan tercermin di URL (refresh mempertahankan tampilan).
- [ ] Requester tidak melihat filter staff maupun kolom pelapor.
- [ ] Debounce search tidak memicu request per keystroke (bukti network log).
- [ ] Empty vs no-results berbeda teks.

**Verifikasi & bukti selesai:** hasil tes komponen; screenshot desktop+mobile.

---

### TASK-029 — UI form buat tiket

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-08 |
| Referensi | DESIGN D-04; API_SPEC §6; FR-04; RULES §7; TC-30 (bagian create) |
| Dependensi | TASK-023, TASK-028 |
| Estimasi | 3–5 jam |

**Tujuan:** `/tickets/new`: form judul→kategori→prioritas→deskripsi (+lampiran menyusul TASK-037), contoh deskripsi, penjelasan prioritas, validasi klien, draft dipertahankan saat gagal.

**File rencana (usulan):** `frontend/src/features/tickets/NewTicketPage.tsx`, komponen `FormField`, `ErrorSummary`, async-select kategori (dari TASK-019).

**Langkah implementasi:**
1. Urutan field & copy sesuai D-04; prioritas dengan penjelasan dampak; catatan urgent ≠ jaminan SLA; contoh deskripsi (masalah, waktu kejadian, langkah yang dicoba).
2. Validasi klien (batas RULES §7) sebagai bantuan; backend tetap otoritatif — tampilkan field errors dari 422.
3. Submit → POST /tickets → sukses: invalidasi cache list → redirect detail. Gagal: draft dipertahankan (DESIGN §2 prinsip 3).
4. Hasil create tidak diketahui (network error) → pesan minta periksa daftar tiket terbaru sebelum kirim ulang; **jangan auto-retry create** (DESIGN D-04, API_SPEC §10).
5. Struktur dua langkah untuk lampiran disiapkan (state file lokal) — eksekusi upload di TASK-037.
6. Tes komponen: validasi, error mapping 422, draft dipertahankan.

**Acceptance criteria:**
- [ ] Tiket valid terkirim → redirect ke detail dengan ticket_code tampil.
- [ ] Error 422 menampilkan pesan per field; input tidak hilang.
- [ ] Tidak ada retry otomatis create saat network gagal.
- [ ] Kategori nonaktif tidak muncul di pilihan (include_inactive=false).

**Verifikasi & bukti selesai:** hasil tes komponen; rekaman alur create.

---

### TASK-030 — UI detail tiket: panel info, aksi lifecycle, conflict UX

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M2 / P0 / todo |
| DEV | DEV-08 |
| Referensi | DESIGN D-05, §2, §3 (halaman tidak ditemukan), §8; API_SPEC §2 (allowed_actions), §6; FR-07; TC-09 (draft UI dipertahankan), TC-35 |
| Dependensi | TASK-024, TASK-025, TASK-026, TASK-028 |
| Estimasi | 4–6 jam |

**Tujuan:** `/tickets/:id`: header + kolom utama + panel kanan (±320 px desktop), tombol aksi sesuai allowed_actions, dialog transisi (ResolutionDialog, ConfirmDialog, reason untuk reopen/reassign/priority), penanganan 409 konflik version sesuai D-05.

**File rencana (usulan):** `frontend/src/features/tickets/TicketDetailPage.tsx`, `TicketInfoPanel`, `TicketActions`, `ResolutionDialog`, `ConfirmDialog`, `AssignDialog` (admin), placeholder tab percakapan/timeline (diisi TASK-036) dan lampiran (TASK-037).

**Langkah implementasi:**
1. Render detail: nomor, judul, status/prioritas badge, waktu WITA, deskripsi, panel pelapor/kategori/assignee, resolusi bila ada, reopen_count.
2. Aksi utama per kondisi D-05 (Ambil tiket / Mulai penanganan / Tandai selesai ditangani / Konfirmasi selesai + Buka kembali / banner closed read-only) — dirender dari `allowed_actions`, bukan role saja.
3. Dialog: resolve → form ringkasan 20–4.000; reopen/reassign/priority → reason 10–1.000; loading & error state; focus management dialog + fokus kembali ke pemicu (DESIGN §8).
4. Konflik 409: pesan "Tiket berubah sejak halaman ini dibuka. Muat data terbaru sebelum menyimpan."; **draft dipertahankan** (TC-09); aksi muat ulang data; invalidate query detail.
5. Setelah mutasi sukses: invalidation list/detail/dashboard/notifikasi terkait (ARCHITECTURE §7).
6. 404 → halaman "Tiket tidak ditemukan" (DESIGN §3 — tidak mengonfirmasi keberadaan objek).
7. Tes komponen: render per allowed_actions (matriks peran × status), dialog resolve validasi, konflik 409 mempertahankan draft.

**Acceptance criteria:**
- [ ] Tombol aksi hanya muncul sesuai allowed_actions untuk 3 peran × 4 status (matriks tes komponen).
- [ ] Alur claim → start → resolve → close dapat dijalankan dua akun berbeda (manual, dicatat).
- [ ] 409 menampilkan pesan konflik + draft tetap ada; reload data memperbaiki version.
- [ ] Closed → banner + seluruh kontrol mutasi hilang (kecuali mark-read notifikasi — RULES §3).
- [ ] URL tiket di luar scope → halaman "Tiket tidak ditemukan".

**Verifikasi & bukti selesai:** hasil tes komponen; catatan/screenshot alur manual lintas peran.

---

## M3 — Kolaborasi dan Berkas

### TASK-031 — Comments backend (public/internal, isolasi requester, first_response_at)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-09 |
| Referensi | FR-08; BR-09, BR-10, BR-12; RULES §2, §5 (komentar tanpa bump version, tetap lock); API_SPEC §7; PRD §9 (first response); TC-14, TC-15, TC-19 |
| Dependensi | TASK-021, TASK-026 |
| Estimasi | 4–6 jam |

**Tujuan:** `GET/POST /tickets/:id/comments` dengan visibility server-side; requester tidak pernah menerima internal (filter sebelum count/pagination); first_response_at diisi komentar publik staff selain reporter bila NULL tanpa bump version; closed menolak komentar.

**File rencana (usulan):** `backend/internal/comments/handler.go`, `service.go`, `repository.go`, tes.

**Langkah implementasi:**
1. POST: validasi body 1–5.000 (trim, Unicode); visibility default public; requester hanya boleh public (mengirim internal → 403); technician hanya bila assignee; admin semua kecuali closed; lock tiket, cek status terkini — closed → 409 (TC-19: tiket ditutup bersamaan → operasi yang melihat closed ditolak).
2. first_response_at: komentar publik oleh staff selain reporter dan first_response_at NULL → isi; version inti tidak naik (RULES §5).
3. GET: scope policy; requester → hanya public (filter visibility pada WHERE **sebelum** count/limit); requester mengirim `visibility=internal` → 403; staff tanpa filter → keduanya.
4. Event `ticket.comment_added` (payload comment_id; visibility event mengikuti komentar) + notifikasi RULES §6 — transaksional; tanpa body komentar di payload/notifikasi.
5. Append-only: tidak ada endpoint update/delete (BR-09/18); runtime grants TASK-027 melindungi.
6. Tes: TC-15 untuk jalur komentar (penanda internal tidak muncul di GET comments requester, count, meta.total); TC-19 (komentar vs close bersamaan, 2 koneksi); validasi batas; first_response_at terisi tepat sekali dan tidak oleh komentar reporter sendiri.

**Acceptance criteria:**
- [ ] Requester tidak dapat membaca maupun membuat komentar internal (403/tersaring).
- [ ] Komentar pada closed → 409; data lain tidak berubah (TC-14).
- [ ] first_response_at hanya dari komentar publik staff selain reporter pertama; version tiket tidak naik karena komentar.
- [ ] Event+notifikasi atomik; payload hanya comment_id.
- [ ] Pagination/count konsisten dengan filter visibility.

**Verifikasi & bukti selesai:** hasil tes API contract comments + konkurensi TC-19.

---

### TASK-032 — Events timeline API (GET /tickets/:id/events)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-09 |
| Referensi | FR-10; API_SPEC §7 (DTO event, seq string); SCHEMA §6, §7 (urut seq); BR-10; TC-15 (jalur timeline) |
| Dependensi | TASK-021, TASK-023 |
| Estimasi | 2–4 jam |

**Tujuan:** endpoint timeline: paginated, terbaru dahulu by seq, visibility difilter server (requester tidak menerima event internal), DTO: id, seq string, event_type, visibility, actor {id,name}, payload allowlist, created_at.

**File rencana (usulan):** `backend/internal/tickets/events_handler.go`, `events_repository.go`, tes.

**Langkah implementasi:**
1. Query events dengan scope policy + filter visibility sebelum count/limit; ORDER BY seq DESC; seq dikirim sebagai string (API_SPEC §1 — BIGINT bukan JS number).
2. Proyeksi payload: hanya field allowlist SCHEMA §6 per event_type (server-side projection).
3. Tes: penanda internal tidak bocor ke requester (bagian TC-15); urutan seq stabil; event reopen memuat reason & previous_assignee_id; event resolved memuat resolution_summary.

**Acceptance criteria:**
- [ ] Requester hanya menerima event public; staff menerima keduanya.
- [ ] DTO sesuai kontrak (seq string, actor ringkas, payload allowlist).
- [ ] Pagination akurat pasca-filter.

**Verifikasi & bukti selesai:** hasil tes API contract events.

---

### TASK-033 — Storage adapter + validasi file (staging privat)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-10 |
| Referensi | SECURITY §6; ARCHITECTURE §4, §6.3; FR-09; BR-11; TC-17; OPERATIONS §2 (STORAGE_ROOT); REFERENCES (OWASP File Upload) |
| Dependensi | TASK-004 |
| Estimasi | 4–6 jam |

**Tujuan:** interface storage (Put/Move/Delete/Open) dengan implementasi filesystem privat di STORAGE_ROOT; validasi upload: extension dinormalisasi + MIME deteksi server + magic-byte signature (jpeg/png/pdf), ukuran 1–5.242.880 byte, streaming dengan batas byte, staging privat, key final acak, sanitasi nama tampilan 1–180 char, tolak path traversal.

**File rencana (usulan):** `backend/internal/attachments/storage.go` (interface), `storage_fs.go`, `validate.go`, `staging.go`, tes unit + integrasi.

**Langkah implementasi:**
1. Interface storage agar service tidak bergantung filesystem (ARCHITECTURE §4); implementasi FS: direktori staging + final di luar webroot, key acak server, tanpa nama asli pada path.
2. Validasi berlapis SECURITY §6: cocokkan extension ↔ MIME terdeteksi ↔ signature; file kosong → tolak; >5 MiB → 413 (streaming limit, jangan baca penuh ke memori); nama file berbahaya (`../`, control char, >180) → sanitasi untuk tampilan, tolak sebagai path.
3. Alur staging: tulis ke staging privat dengan batas 6 MiB multipart (termasuk overhead) → validasi → siap dipindah (move terjadi di TASK-034 dalam transaksi).
4. Hook failure injection untuk tes kompensasi (TEST_PLAN §5: boundary file move).
5. Tes TC-17: MIME palsu (exe berekstensi .png), extension terlarang (.php/.html), path traversal pada filename, kosong, oversize → masing-masing error contract tepat (415/413/400); signature check dibuktikan unit test.

**Acceptance criteria:**
- [ ] Hanya jpeg/png/pdf valid (extension+MIME+signature cocok) yang lolos staging.
- [ ] Oversize → 413; MIME mismatch → 415; file kosong/traversal ditolak; tidak ada file tertulis di luar staging/final STORAGE_ROOT.
- [ ] Nama asli tersimpan hanya sebagai metadata tersanitasi.
- [ ] Streaming: file 5 MiB tidak menaikkan memori proses secara signifikan (bukti pengamatan).

**Verifikasi & bukti selesai:** hasil tes validasi file dengan fixture generator TASK-017.

---

### TASK-034 — Attachment upload backend (limit 5 aktif, transaksi, kompensasi)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-10 |
| Referensi | FR-09; BR-11, BR-12; ARCHITECTURE §6.3; API_SPEC §7 (multipart 1 file, 201 DTO, 413/415/422); RULES §2 (izin upload); TC-16, TC-18 (sisi commit gagal), TC-30 (sisi API) |
| Dependensi | TASK-026, TASK-033 |
| Estimasi | 4–6 jam |

**Tujuan:** `POST /tickets/:id/attachments`: tepat satu file per request; hitung lampiran aktif ≤5 di bawah ticket lock; move staging→final + insert metadata+event+notifikasi dalam transaksi; commit gagal → hapus file final (kompensasi) + log.

**File rencana (usulan):** `backend/internal/attachments/handler.go`, `upload_service.go`, `repository.go`, tes integrasi + failure injection.

**Langkah implementasi:**
1. Izin upload via policy (requester pemilik / assignee / admin; status open/in_progress saja; closed/resolved → 409).
2. Parse multipart: tepat satu field file (lebih → 400); batas 6 MiB request; validasi TASK-033.
3. Transaksi ARCHITECTURE §6.3: lock tiket → hitung lampiran aktif → ≥5 → 422 ATTACHMENT_LIMIT (tanpa memindah file) → insert metadata (storage_key final, sha256, size, original_name tersanitasi, content_type) + event `ticket.attachment_added` (payload attachment_id + original_name sanitized) + notifikasi → commit.
4. Pindah file ke key final **sebelum commit**; bila commit gagal → hapus file final (kompensasi); kegagalan cleanup dilog untuk operator (ARCHITECTURE §6.3).
5. DTO respons: tanpa storage_key/checksum; memuat can_delete & download_path (API_SPEC §7).
6. Tes: TC-16 (dua upload bersamaan saat 4 aktif → tepat satu menjadi kelima, lainnya 422, file sisa dibersihkan); TC-18 sebagian (commit gagal → tidak ada file final orphan, tidak ada metadata); closed → 409; multi-file → 400.

**Acceptance criteria:**
- [ ] Tidak pernah >5 lampiran aktif per tiket (tes konkurensi diulang ≥20x, verifikasi DB + filesystem).
- [ ] Commit gagal → file final terkompensasi terhapus; log mencatat; tiket tidak berubah.
- [ ] DTO tidak membocorkan storage_key/sha256.
- [ ] Upload pada closed/resolved → 409 INVALID_STATE.

**Verifikasi & bukti selesai:** hasil tes upload + failure injection; listing STORAGE_ROOT pasca-tes (tanpa orphan).

---

### TASK-035 — Attachment download & delete + cleanup staging/orphan

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-10 |
| Referensi | FR-09; SECURITY §6 (Content-Disposition, nosniff, otorisasi download); ARCHITECTURE §6.3 (soft delete, kompensasi byte, maintenance 24 jam); API_SPEC §7; RULES §2 (hapus lampiran); TC-04, TC-18; OPERATIONS §7 |
| Dependensi | TASK-034 |
| Estimasi | 3–5 jam |

**Tujuan:** `GET …/download` terotorisasi (attachment_id harus cocok ticket & akses actor; soft-deleted → ditolak; header attachment+nosniff), `DELETE` soft-delete transaksional + event + hapus byte pasca-commit dengan retry, dan perintah maintenance pembersih staging/orphan >24 jam.

**File rencana (usulan):** `backend/internal/attachments/download_handler.go`, `delete_service.go`, `backend/cmd/admin/` subcommand `cleanup-attachments`, tes.

**Langkah implementasi:**
1. Download: validasi nested route (attachment bukan milik tiket → 404); scope policy; `Content-Disposition: attachment`, `X-Content-Type-Options: nosniff`; streaming; deleted_at terisi → ditolak.
2. Delete: izin RULES §2 (pengunggah dengan akses & open/in_progress; teknisi masih assignee; admin); lock tiket → set deleted_at/deleted_by (CHECK attachments_delete_pair_ck) + event `ticket.attachment_deleted` + notifikasi → commit → hapus byte; gagal hapus → log + dibersihkan maintenance (ARCHITECTURE §6.3).
3. Maintenance: hapus file staging dan file final tanpa metadata yang lebih tua dari 24 jam **setelah** verifikasi tidak direferensikan DB; dry-run + log jumlah + batas batch (OPERATIONS §7: jangan menghapus file hanya karena query gagal).
4. Tes: TC-04 (R1 download lampiran R2 → 404); download soft-deleted → ditolak (TC-18 bagian akhir); delete dengan gagal-hapus-byte → metadata tetap soft-deleted, orphan dibersihkan maintenance; cleanup tidak menyentuh file aktif/referenced.

**Acceptance criteria:**
- [ ] Download selalu melewati otorisasi tiket; header keamanan terpasang; tidak ada public static mapping ke STORAGE_ROOT.
- [ ] Soft-deleted tidak dapat diunduh namun metadata/timeline tetap ada.
- [ ] Maintenance membersihkan orphan staging/final >24 jam tanpa menghapus file aktif (tes dengan fixture umur file).
- [ ] Delete pada closed → 409.

**Verifikasi & bukti selesai:** hasil tes download/delete/cleanup; listing filesystem sebelum/sesudah maintenance.

---

### TASK-036 — UI komentar + timeline

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-09 (sisi UI) |
| Referensi | DESIGN D-05 (tab Publik/Internal, label permanen, draft tidak terkirim saat pindah tab), §2 prinsip 6, §7 (CommentComposer/List, ActivityTimeline); FR-08; TC-15; SECURITY §5 (text node) |
| Dependensi | TASK-030, TASK-031, TASK-032 |
| Estimasi | 4–6 jam |

**Tujuan:** percakapan & riwayat pada halaman detail: tab Publik/Internal untuk staff (requester hanya publik), composer dengan draft per-tab, timeline event dengan label aktor/waktu absolut WITA dan penanda staff untuk internal.

**File rencana (usulan):** `frontend/src/features/tickets/comments/`, `frontend/src/features/tickets/timeline/`, komponen `CommentComposer`, `CommentList`, `ActivityTimeline`.

**Langkah implementasi:**
1. Query terpisah untuk komentar publik vs internal (cache key berbeda — ARCHITECTURE §7); requester tidak pernah meminta internal.
2. Composer: validasi 1–5.000; label permanen "Hanya terlihat oleh tim IT" pada tab internal; pindah tab tidak mengirim draft; submit → append + invalidasi; body dirender sebagai text node (tanpa raw HTML).
3. Timeline: urut seq, event_type → label manusiawi Bahasa Indonesia, aktor, waktu WITA; event internal berlabel staff; payload reason/resolution_summary dirender aman.
4. Closed → composer disabled + banner read-only.
5. State loading/empty ("Belum ada komentar")/error; aria-live untuk pesan sukses singkat.
6. Tes komponen: tab internal tidak dirender untuk requester; draft per-tab independen; closed menonaktifkan composer.

**Acceptance criteria:**
- [ ] Requester tidak melihat tab/komposer/event internal pada UI manapun.
- [ ] Komentar publik & internal terkirim dan muncul sesuai visibility; append-only (tanpa tombol edit/hapus).
- [ ] Timeline menampilkan seluruh event dengan label dan waktu WITA benar.
- [ ] Pindah tab tidak mengirim draft; closed read-only.

**Verifikasi & bukti selesai:** hasil tes komponen; screenshot detail tiket sebagai staff dan requester.

---

### TASK-037 — UI lampiran (picker, upload per file, list, download, delete)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M3 / P0 / todo |
| DEV | DEV-10 (sisi UI) |
| Referensi | DESIGN D-04 (alur create+upload terpisah, pesan "Dua lampiran belum terunggah", tombol coba lagi), D-05 (ikon+nama+ukuran, tanpa preview inline); FR-09; TC-30; API_SPEC §7 |
| Dependensi | TASK-029, TASK-030, TASK-035 |
| Estimasi | 4–6 jam |

**Tujuan:** komponen AttachmentPicker (validasi awal klien: jumlah ≤5, ukuran ≤5 MiB, tipe), AttachmentList (unduh via API, can_delete), alur create-ticket dua langkah dengan retry per file gagal, dan upload/hapus pada halaman detail.

**File rencana (usulan):** `frontend/src/features/tickets/attachments/`, komponen `AttachmentPicker`, `AttachmentList`.

**Langkah implementasi:**
1. Picker: cek awal ukuran/jenis sebelum kirim (bantuan UX; backend otoritatif); hasil per file ditampilkan (sukses/gagal/alasan).
2. Create flow D-04: tiket dibuat dahulu (JSON) → unggah tiap file; bila sebagian gagal → tetap redirect ke detail dengan pesan "Tiket berhasil dibuat. N lampiran belum terunggah" + tombol coba lagi per file; **jangan** ulang create-ticket; hasil create tidak diketahui → minta periksa daftar (DESIGN D-04).
3. Detail: list lampiran aktif (ikon tipe, nama tersanitasi, ukuran, pengunggah, waktu WITA), tautan download lewat endpoint terotorisasi (tanpa preview inline), tombol hapus bila can_delete + ConfirmDialog.
4. Error mapping: 413/415/422 ATTACHMENT_LIMIT tampil sebagai pesan per file yang jelas (TC-30).
5. Closed/resolved → picker disabled (mengikuti allowed_actions).
6. Tes komponen: validasi awal, render hasil per file, retry hanya file gagal.

**Acceptance criteria:**
- [ ] Alur create + 2 lampiran (1 gagal) → tiket ada, pesan tepat, retry hanya file gagal (TC-30).
- [ ] Download melewati API dengan cookie sesi; URL langsung storage tidak tersedia.
- [ ] Hapus hanya muncul bila can_delete; soft-deleted hilang dari list aktif.
- [ ] Batas klien (5 file/5 MiB/tipe) mencegah request sia-sia, tetapi error server tetap ditangani.

**Verifikasi & bukti selesai:** hasil tes komponen; rekaman alur create dengan lampiran gagal.

---

## M4 — Notifikasi dan Dashboard

### TASK-038 — Notifications API (list, unread_count, mark-read)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M4 / P0 / todo |
| DEV | DEV-11 |
| Referensi | FR-11; API_SPEC §8 (kontrak DTO, unread_count, mark-read idempoten, 404 milik orang lain); BR-10 (isi aman); SCHEMA §3 (notifications); TC-21 |
| Dependensi | TASK-021 |
| Estimasi | 3–4 jam |

**Tujuan:** `GET /notifications` (unread_only opsional, paginated, meta.unread_count milik actor) dan `PATCH /notifications/:id/read` (idempoten, mempertahankan read_at pertama, 404 untuk milik orang lain).

**File rencana (usulan):** `backend/internal/notifications/handler.go`, `read_service.go`, `repository.go`, tes.

**Langkah implementasi:**
1. Query selalu `recipient_id = actor` (BR-16); DTO: id, event_id, ticket_id (join event), ticket_code, message, read_at, created_at; unread_count dihitung untuk seluruh notifikasi actor (bukan hanya halaman kini).
2. Mark-read: `{}` + CSRF; UPDATE read_at hanya bila masih NULL (idempoten, read_at pertama dipertahankan); notifikasi orang lain → 404 (TC-21); boleh pada tiket closed (RULES §3).
3. Tidak ada endpoint yang menerima recipient_id dari klien (API_SPEC §8).
4. Tes: TC-21 (R1 mark-read notifikasi R2 → 404, read_at R2 tidak berubah); idempotency mark-read; unread_count akurat pasca-filter; isi notifikasi tidak memuat penanda internal (bagian TC-15/TC-20).

**Acceptance criteria:**
- [ ] List hanya memuat notifikasi milik actor; unread_count benar.
- [ ] Mark-read berulang → read_at tidak berubah setelah yang pertama.
- [ ] Mark-read milik orang lain → 404 tanpa efek.
- [ ] DTO sesuai API_SPEC §8 (ticket_code string, tanpa data internal).

**Verifikasi & bukti selesai:** hasil tes API contract notifications.

---

### TASK-039 — UI notifikasi + polling lifecycle 30 detik

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M4 / P0 / todo |
| DEV | DEV-11 (sisi UI) |
| Referensi | DESIGN D-06; FR-11; ARCHITECTURE §8; README paket §2 (polling 30 detik tab aktif); TC-22 |
| Dependensi | TASK-012, TASK-038 |
| Estimasi | 4–6 jam |

**Tujuan:** bell dengan unread_count di topbar, halaman `/notifications`, polling 30 detik hanya saat tab visible + pengguna bersesi + online, backoff setelah gagal, refetch saat tab kembali aktif dan setelah aksi sendiri, mark-read + navigasi ke detail.

**File rencana (usulan):** `frontend/src/features/notifications/` (NotificationBell, NotificationsPage, useNotificationPolling).

**Langkah implementasi:**
1. TanStack Query `refetchInterval` 30_000 dengan `refetchIntervalInBackground: false`; pause saat `navigator.onLine` false; resume + refetch saat event `visibilitychange`/`online` (TC-22 — tanpa request loop).
2. Refetch setelah mutasi sendiri (invalidation dari TASK-006 convention).
3. Bell: unread_count badge; klik item → PATCH read (kegagalan mark-read tidak menghalangi navigasi — D-06) → navigate detail tiket.
4. Halaman `/notifications`: list paginated, filter unread, status dibaca visual + label.
5. Pembaruan tidak memindahkan focus keyboard (D-06); aria-live untuk perubahan jumlah.
6. Backoff eksponensial terbatas setelah kegagalan beruntun; kembali ke interval normal setelah sukses.
7. Tes komponen: polling berhenti saat hidden/offline (fake timers + mock visibility), backoff, mark-read gagal tetap navigasi.

**Acceptance criteria:**
- [ ] Tidak ada request polling saat tab hidden/offline/tanpa sesi (bukti network log).
- [ ] Refetch terjadi saat tab aktif kembali dan setelah aksi sendiri.
- [ ] Kegagalan mark-read tidak menghalangi navigasi ke tiket.
- [ ] Focus keyboard tidak berpindah karena pembaruan daftar.

**Verifikasi & bukti selesai:** hasil tes komponen (fake timers); rekaman network log siklus hidden→visible.

---

### TASK-040 — Dashboard backend (agregat sesuai PRD §9)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M4 / P0 / todo |
| DEV | DEV-12 |
| Referensi | FR-12; PRD §9 (definisi metrik otoritatif); API_SPEC §9 (payload, cohort from/to); RULES §2 (scope dashboard); BR-16; TC-23 |
| Dependensi | TASK-022, TASK-026 |
| Estimasi | 4–6 jam |

**Tujuan:** `GET /dashboard/summary`: status_counts, active_backlog (open+in_progress), rata-rata first response & resolution (+sample_count, null tanpa sampel), category_counts (termasuk kategori nonaktif); scope requester/staff dari sesi; cohort from/to sama dengan daftar tiket (created_at, from inklusif/to eksklusif UTC).

**File rencana (usulan):** `backend/internal/dashboard/service.go`, `repository.go`, `handler.go`, tes.

**Langkah implementasi:**
1. Satu predicate scope policy yang sama dengan list (ARCHITECTURE §4 — hindari salinan logika scope).
2. Query agregat: per status; first response = avg(first_response_at - created_at) untuk yang memiliki first_response_at; resolution = avg(resolved_at - created_at) untuk resolved/closed dengan resolved_at; durasi kalender (BR-13); tanpa sampel → null + sample_count 0 (API_SPEC §9).
3. Category_counts per kategori dalam cohort, termasuk kategori nonaktif (PRD §9).
4. Validasi from/to sama dengan TASK-022 (berpasangan, from<to, RFC 3339).
5. Tes TC-23 dengan fixture terkontrol (TASK-017): cohort filter benar (boundary from/to), reopen tidak mereset first_response, resolved_at mengikuti penyelesaian terakhir, nilai kosong → null, scope requester hanya menghitung tiket sendiri, kategori nonaktif tetap terhitung.

**Acceptance criteria:**
- [ ] Seluruh metrik cocok dengan perhitungan manual atas fixture (bukan angka yang diatur agar lulus).
- [ ] Tanpa sampel → null dan sample_count 0, bukan 0 detik.
- [ ] Scope requester vs staff konsisten dengan list/search (policy sama).
- [ ] from/to UTC boundary inklusif/eksklusif terbukti di tes.

**Verifikasi & bukti selesai:** hasil tes dashboard + tabel perbandingan nilai fixture vs hasil API.

---

### TASK-041 — Dashboard UI

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M4 / P0 / todo |
| DEV | DEV-12 (sisi UI) |
| Referensi | DESIGN D-02; FR-12; PRD §9 (konversi hari WITA → rentang UTC); TC-23 (bagian UI) |
| Dependensi | TASK-012, TASK-040 |
| Estimasi | 4–6 jam |

**Tujuan:** `/dashboard` sesuai D-02: judul scope-aware ("Ringkasan tiket saya" vs "Ringkasan layanan"), 4 kartu status yang menaut ke daftar terfilter, backlog + metrik durasi dengan jumlah sampel, distribusi kategori (bar horizontal + padanan tabel), filter tanggal, "Belum ada data" untuk metrik kosong.

**File rencana (usulan):** `frontend/src/features/dashboard/DashboardPage.tsx`, komponen kartu/bar.

**Langkah implementasi:**
1. Filter tanggal: user memilih hari WITA; klien mengonversi ke rentang UTC from/to (PRD §9) sebelum kirim; tooltip menjelaskan basis created_at.
2. Kartu status klik → `/tickets?status=…`; metrik durasi ditampilkan dalam format manusiawi (jam/hari kalender) + sample_count di bawahnya; null → "Belum ada data".
3. Distribusi kategori: bar horizontal berlabel jumlah + tabel padanan untuk screen reader (D-02).
4. Invalidasi setelah mutasi tiket (convention TASK-006); state loading/empty/error.
5. Tes komponen: render null metric, konversi tanggal WITA→UTC, link kartu status.

**Acceptance criteria:**
- [ ] Metrik kosong menampilkan "Belum ada data", bukan nol.
- [ ] Filter tanggal mengirim rentang UTC yang benar untuk hari WITA yang dipilih (tes unit konversi, termasuk kasus lintas tengah malam).
- [ ] Kartu status menaut ke daftar dengan filter sesuai.
- [ ] Distribusi kategori memiliki padanan tabel yang terbaca screen reader.

**Verifikasi & bukti selesai:** hasil tes komponen; screenshot dashboard staff & requester.

---

## M5 — Hardening dan UX

### TASK-042 — Rate limiting + security headers + cookie produksi

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | SECURITY §3 (nama cookie produksi), §7 (batas trafik), §8 (header, HSTS, CSP); API_SPEC §3 (429 + Retry-After); TC-34 (bagian header); NFR-05 |
| Dependensi | TASK-010, TASK-034 |
| Estimasi | 3–5 jam |

**Tujuan:** limiter memory satu instance sesuai tabel SECURITY §7 (login per IP 20/15mnt; login per email+IP 5 gagal/15mnt; API umum 120/mnt burst 30 per user; upload 10/mnt), respons 429 + Retry-After, header keamanan produksi (nosniff, Referrer-Policy, CSP default-src 'self' dst., HSTS saat HTTPS stabil), cookie `__Host-` di production, trusted proxy headers hanya dari TRUSTED_PROXIES.

**File rencana (usulan):** `backend/internal/platform/httpx/middleware/ratelimit.go`, `security_headers.go`, pembaruan `cookies.go`, tes.

**Langkah implementasi:**
1. Limiter in-memory (catat keterbatasan reset saat restart — SECURITY §7); kunci per IP untuk login dan per user untuk API/upload; 429 generik login tidak mengungkap keberadaan akun.
2. Middleware header keamanan sesuai CSP SECURITY §8; HSTS hanya ketika APP_ENV production + HTTPS.
3. Cookie produksi `__Host-randesk_session` (Secure, tanpa Domain) vs dev — sudah dari TASK-009, verifikasi konfigurasi per lingkungan.
4. Trusted proxies Gin dikonfigurasi dari TRUSTED_PROXIES (bukan default semua).
5. Tes: 429 + Retry-After setelah batas terlampaui (login & API umum); header terpasang di respons; limiter tidak memblokir permanen; IP dari X-Forwarded-For hanya dipercaya dari proxy terkonfigurasi.

**Acceptance criteria:**
- [ ] Seluruh batas SECURITY §7 terpasang dan teruji minimal untuk login dan API umum.
- [ ] Respons 429 memuat Retry-After; pesan tidak mengungkap akun.
- [ ] Header CSP/nosniff/Referrer-Policy ada pada seluruh respons; HSTS hanya produksi HTTPS.
- [ ] Cookie produksi memenuhi atribut `__Host-`.

**Verifikasi & bukti selesai:** hasil tes rate limit/header; contoh header respons produksi-like.

---

### TASK-043 — Verifikasi log redaction + request_id tracing

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | NFR-08; SECURITY §3 (jangan log cookie/token/password), §8; ARCHITECTURE §5 poin 6, §10; TC-34; OPERATIONS §8 (route template) |
| Dependensi | TASK-004, TASK-010 |
| Estimasi | 2–4 jam |

**Tujuan:** membuktikan seluruh jalur log (middleware, error, service, upload) bebas secret dan body sensitif; request_id konsisten header↔body↔log; route template dipakai sebagai dimensi log.

**File rencana (usulan):** pembaruan `backend/internal/platform/logger/`, tes redaction.

**Langkah implementasi:**
1. Inventarisasi titik log; pastikan tidak ada yang mencatat body, cookie, Set-Cookie, Authorization, password, csrf_token, temporary_password, storage_key.
2. Tes otomatis: jalankan alur login/create/upload terhadap buffer log; assert string sensitif (password fixture, token, dsb.) tidak muncul; assert request_id muncul konsisten.
3. Verifikasi log memakai route template (`/tickets/:id`), bukan path berisi UUID (OPERATIONS §8).
4. Periksa pesan error 500 tidak membocorkan SQL/stack ke klien (API_SPEC §2).

**Acceptance criteria:**
- [ ] Tes redaction lulus untuk alur login, change-password, create user (temporary_password), create ticket, upload.
- [ ] request_id di header = body = entri log untuk request yang sama.
- [ ] Dimensi log memakai route template.

**Verifikasi & bukti selesai:** hasil tes redaction + contoh entri log (nilai sensitif disamarkan).

---

### TASK-044 — Sweep negative tests + gate keamanan SECURITY §10

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | SECURITY §10 (daftar gate); NFR-01; TC-04, TC-05, TC-10, TC-15, TC-21, TC-27; BR-16; PRD §4.2 (0 kebocoran pada fixture) |
| Dependensi | TASK-031, TASK-035, TASK-038 |
| Estimasi | 4–6 jam |

**Tujuan:** satu suite negatif lintas endpoint yang membuktikan seluruh gate keamanan MVP: lintas requester (404 di semua jalur objek), isolasi internal di semua jalur (list/detail/comment/event/search/count/notifikasi/dashboard), forged role/reporter/status, CSRF, revocation, stale version, path/size/type file, input SQL pada q/filter/sort.

**File rencana (usulan):** `backend/tests/security/` (suite khusus), memakai fixture TASK-017 dengan penanda internal unik.

**Langkah implementasi:**
1.Enumerasi seluruh endpoint; untuk tiap endpoint jalankan matriks aktor (R1/R2/T1/T2/A1/akun nonaktif/tanpa sesi) × objek fixture; assert 404/403 sesuai kontrak (SECURITY §5: luar scope baca → 404; aksi terlarang objek terlihat → 403).
2. Pencarian penanda internal (string unik dari fixture TEST_PLAN §2) pada **seluruh** respons requester — harus nol kemunculan (TC-15 menyeluruh; PRD §4.2 target 0 kebocoran).
3. Uji forged input: requester mengirim role/reporter/status/visibility internal (TC-05); sort/q/filter injeksi (TC-27).
4.Dependency scan (govulncheck/npm audit) dan kajian temuan kritis (SECURITY §10) — hasil dicatat, bukan auto-pass.
5. Defect yang ditemukan → task perbaikan baru + regression test (OPERATIONS §9 pola bug otorisasi).

**Acceptance criteria:**
- [ ] Seluruh gate SECURITY §10 memiliki tes otomatis yang lulus: lintas requester, internal isolation, CSRF, revocation, forged role/reporter, file path/size/type, SQL input, stale version, log redaction (TASK-043).
- [ ] 0 kemunculan penanda internal pada semua respons requester.
- [ ] Hasil dependency scan dikaji; temuan kritis ditutup atau dicatat keputusannya.

**Verifikasi & bukti selesai:** laporan run suite negatif (tanggal, commit, hasil per gate) sesuai template TEST_PLAN §9.

---

### TASK-045 — Suite konkurensi + failure injection berulang

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | TEST_PLAN §5 (barrier, 2 klien independen, failure injection, tanpa mock untuk row lock); TC-07, TC-08, TC-16, TC-18, TC-19, TC-24, TC-25; NFR-02, NFR-05 |
| Dependensi | TASK-015, TASK-025, TASK-034 |
| Estimasi | 4–6 jam |

**Tujuan:** mengonsolidasikan tes konkurensi menjadi suite yang dijalankan berulang (bukan sekali): claim bersaing, attachment cap, deaktivasi vs assignment, dual-admin deactivation, comment vs close; plus failure injection di boundary event insert, notification insert, file move, commit — dengan verifikasi state akhir DB dan filesystem.

**File rencana (usulan):** `backend/tests/concurrency/` (build tag atau suite terpisah yang juga jalan di CI job race).

**Langkah implementasi:**
1. Harness 2 klien HTTP/DB independen + barrier (TEST_PLAN §5); tiap kasus diulang ≥20x per run.
2. Kasus: TC-08 (claim), TC-16 (cap lampiran), TC-24 (deactivate vs claim — tidak pernah ada tiket aktif ditugaskan ke teknisi nonaktif), TC-25 (dual admin — selalu tersisa admin aktif), TC-19 (comment/upload vs close).
3. Failure injection: event insert gagal, notification insert gagal, file move gagal, commit gagal → assert rollback + kompensasi + tidak ada orphan (TC-07/18).
4. Jalankan dengan `go test -race`; verifikasi state akhir lewat query DB dan listing filesystem, bukan hanya kode respons.
5. Integrasikan ke CI (job race-sensitive TEST_PLAN §6).

**Acceptance criteria:**
- [ ] Seluruh kasus kritis konkurensi lulus pada run berulang (≥20x) tanpa flake.
- [ ] Failure injection membuktikan rollback dan kompensasi file (tidak ada orphan pasca-run).
- [ ] Suite berjalan di CI dengan race detector.

**Verifikasi & bukti selesai:** log run berulang (jumlah iterasi, hasil), listing filesystem pasca failure-injection.

---

### TASK-046 — Pass aksesibilitas + responsif

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | NFR-04; DESIGN §4.1 (kontras), §5, §8; TC-31; TEST_PLAN §1 (manual+otomatis) |
| Dependensi | TASK-018, TASK-019, TASK-030, TASK-036, TASK-037, TASK-039, TASK-041 |
| Estimasi | 4–6 jam |

**Tujuan:** layar MVP memenuhi target WCAG 2.2 AA: audit otomatis (axe-core atau setara) + manual (keyboard-only, focus dialog, zoom 200%, lebar 360/768/1024/1440, reduced motion, screen reader untuk login/create-ticket/detail).

**File rencana (usulan):** perbaikan tersebar di `frontend/src/`, laporan audit di `docs/testing/a11y-report.md` (usulan).

**Langkah implementasi:**
1. Pasang audit otomatis di CI atau lokal (axe + vitest/playwright); perbaiki violation serius.
2. Checklist manual DESIGN §8: urutan tab logis, label eksplisit, accessible name icon-only button, focus dialog & restore, aria-live sukses, kontras (verifikasi token DESIGN §4.1 pada komponen final termasuk hover/disabled/focus), target 44 px, zoom 200%, reduced motion.
3. Verifikasi 4 lebar viewport; tidak ada status/judul/aksi utama tersembunyi di mobile (DESIGN §5).
4. Screen reader (NVDA/VoiceOver) untuk 3 alur: login, create-ticket, detail; catat temuan.
5. Temuan sedang boleh dicatat sebagai defect UX beralasan + rencana perbaikan (TEST_PLAN §3); kritis/tinggi diperbaiki.

**Acceptance criteria:**
- [ ] Audit otomatis tanpa violation serius (serious/critical) pada layar MVP.
- [ ] Alur utama dapat dijalankan keyboard-only; focus tidak terjebak.
- [ ] Zoom 200% dan lebar 360 px tidak memotong kontrol/teks utama (TC-31).
- [ ] Laporan audit manual (tanggal, layar, temuan, perbaikan) tersimpan — bukan klaim "WCAG compliant" tanpa pemeriksaan (DESIGN §8).

**Verifikasi & bukti selesai:** laporan a11y + screenshot 4 viewport + hasil audit otomatis.

---

### TASK-047 — E2E suite (TC-35 + alur kritis tiga peran)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / todo |
| DEV | DEV-13 |
| Referensi | TC-35; TEST_PLAN §1, §6 (E2E job + artefak screenshot saat gagal); ROADMAP §2 (M5); PRD §4.2 (integritas alur kritis) |
| Dependensi | TASK-036, TASK-037, TASK-039, TASK-041 |
| Estimasi | 4–6 jam |

**Tujuan:** browser otomatis (default usulan: Playwright — catat sebagai keputusan tool) menjalankan alur penuh: login requester → create ticket (+lampiran) → login technician → claim → start → komentar internal+publik → resolve → login requester → close; varian reopen; varian konflik claim 2 teknisi; refresh browser di tengah alur (state pulih via /auth/me); negative: R1 membuka URL tiket R2 → "Tiket tidak ditemukan".

**File rencana (usulan):** `frontend/e2e/` (atau `tests/e2e/`), konfigurasi Playwright, seed khusus E2E dari TASK-017.

**Langkah implementasi:**
1. Setup Playwright + webServer (API + frontend build) dengan DB test ter-seed; screenshot artefak saat gagal (TEST_PLAN §6).
2. Skrip alur utama TC-35 lintas 3 peran dengan assert status/badge/timeline/notifikasi di tiap tahap.
3. Varian: reopen oleh reporter; claim conflict (2 context browser); refresh di halaman detail; closed read-only.
4. Negative: akses silang requester → halaman tidak ditemukan; internal tab tidak ada untuk requester.
5. Masukkan ke CI sebagai job E2E.

**Acceptance criteria:**
- [ ] TC-35 lulus end-to-end tanpa kehilangan state setelah refresh.
- [ ] Varian reopen, claim conflict, dan negative cross-access lulus.
- [ ] E2E berjalan di CI dengan artefak screenshot saat gagal.

**Verifikasi & bukti selesai:** run E2E hijau (log + tanggal + commit), screenshot kunci.

---

### TASK-048 — Benchmark dataset + load test NFR-03

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P0 / **needs_verification** (lingkungan acuan 2 vCPU/4 GiB belum terbukti tersedia — lihat GAP-04) |
| DEV | DEV-13 |
| Referensi | NFR-03 (target & kondisi ukur); TC-32; PRD §8 (metodologi 80/20, pacing, error <1%); SCHEMA §8 (benchmark dataset terpisah); TEST_PLAN §6 |
| Dependensi | TASK-017, TASK-042, TASK-047 |
| Estimasi | 4–6 jam |

**Tujuan:** dataset benchmark (10.000 tiket, 50.000 komentar) terpisah dari seed dev; load test 25 pengguna virtual campuran 80% read/20% write selama 5 menit; catat p95, error rate, CPU/RAM/DB — tanpa mengarang hasil; login/upload dikecualikan dari target.

**File rencana (usulan):** `backend/cmd/admin/` subcommand `seed-benchmark`, skrip load test (k6 atau setara — putuskan & catat), `docs/testing/loadtest-report.md` (usulan).

**Langkah implementasi:**
1. Verifikasi lingkungan: apakah tersedia host 2 vCPU/4 GiB/SSD (lokal VM/staging)? Bila tidak, jalankan pada lingkungan yang ada dan **catat deviasi** — hasil tidak diklaim memenuhi NFR-03 (PRD §8: profilkan dahulu sebelum mengubah target).
2. Generator dataset sesuai skala NFR-03; pisahkan dari seed dev (SCHEMA §8).
3. Skenario: pacing 1 request/user/detik, 80% read (list/detail/dashboard) 20% write (create/comment/transisi), input valid; rate-limit test & konflik dijalankan terpisah (PRD §8).
4. Ukur dan catat: versi software, dataset, error rate, p95 latency, CPU, memori; bandingkan target (read ≤500 ms, write JSON ≤800 ms, error <1%).
5. Bila gagal: profilkan (query plan EXPLAIN, pool wait) sebelum mengusulkan perubahan — usulan perubahan target masuk GAPS_AND_DECISIONS.md, bukan diubah diam-diam.

**Acceptance criteria:**
- [ ] Load test dijalankan pada dataset skala target; laporan memuat seluruh parameter PRD §8.
- [ ] Hasil dicatat apa adanya (lulus/gagal/deviasi lingkungan) — tanpa angka karangan.
- [ ] Tidak ada regresi fungsional pasca-load (suite utama tetap hijau).

**Verifikasi & bukti selesai:** laporan load test (tanggal, commit, environment, hasil mentah).

---

### TASK-049 — UAT / walkthrough terstruktur

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M5 / P1-dalam-MVP / todo (D-06 diputuskan 26 Sep 2026: adopsi resmi walkthrough terstruktur P-07; varian UAT 5 peserta di luar cakupan MVP) |
| DEV | DEV-13 |
| Referensi | PRD §4.2 (target evaluasi); TEST_PLAN §7; DECISIONS D-06 |
| Dependensi | TASK-047 |
| Estimasi | 3–5 jam |

**Tujuan:** menjalankan walkthrough terstruktur oleh pengembang sebagai bukti penerimaan MVP, dengan catatan eksplisit bahwa hasilnya belum mewakili pengguna nyata (TEST_PLAN §7; keputusan D-06 26 September 2026). Bila di kemudian hari 5 peserta tersedia, varian UAT penuh dapat dijalankan ulang tanpa mengubah task ini.

**File rencana (usulan):** `docs/testing/uat-report.md` (usulan) memakai template TEST_PLAN §9.

**Langkah implementasi:**
1. Siapkan skenario tugas peserta TEST_PLAN §7 (membuat tiket; menemukan progres; menambah informasi; mengambil & menangani tiket untuk staff; konfirmasi/reopen).
2. Jalankan tanpa mengarahkan ke tombol sebelum mengamati kesulitan; catat berhasil/gagal, waktu, bantuan, kesalahan, komentar.
3. Triase temuan: kritis → perbaikan sebelum rilis; sedang → defect UX beralasan + rencana (TEST_PLAN §3).
4. Bila hanya walkthrough: tandai laporan "belum mewakili pengguna nyata".

**Acceptance criteria:**
- [ ] Laporan UAT/walkthrough terisi sesuai template (tanggal, peserta/peran, hasil per tugas).
- [ ] Temuan kritis ditutup atau memiliki keputusan risiko tertulis (TEST_PLAN §8).

**Verifikasi & bukti selesai:** laporan UAT tersimpan di `docs/testing/`.

---

## M6 — Deploy dan Portofolio

### TASK-050 — Deployment staging: proxy TLS, konfigurasi produksi, smoke test

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M6 / P0 / todo |
| DEV | DEV-14; D-03 dikunci di sini (default: Caddy, satu VM Linux/container) |
| Referensi | ARCHITECTURE §2 (topologi), §10; OPERATIONS §1–§4; SECURITY §8 (HTTPS, HSTS); ROADMAP §6 (risiko hosting berbeda — staging sejak awal bila memungkinkan) |
| Dependensi | TASK-042, TASK-043, TASK-046 |
| Estimasi | 4–6 jam |

**Tujuan:** aplikasi berjalan di lingkungan staging/produksi-demo: reverse proxy TLS (hanya proxy yang publik), static build React + routing `/api/v1/*` ke Go, SPA fallback, PostgreSQL & volume privat, konfigurasi produksi (APP_ENV, APP_ORIGIN, cookie `__Host-`, TRUSTED_PROXIES), smoke test tiga peran + readiness.

**File rencana (usulan):** `deploy/` (Caddyfile/nginx conf, systemd unit atau Dockerfile+compose, skrip rilis), dokumentasi runbook rilis OPERATIONS §4.

**Langkah implementasi:**
1. Kunci D-03; siapkan VM/container target; pasang PostgreSQL + direktori STORAGE_ROOT persistent.
2. Konfigurasi proxy: TLS, `/api/v1/*` → HTTP_ADDR privat, route SPA → index.html, **tanpa** public static mapping ke volume lampiran (ARCHITECTURE §2).
3. Build artifact reproducible (binary Go + `npm run build`); prosedur rilis OPERATIONS §4: backup → migrasi → deploy → readiness → smoke.
4. Smoke test: login 3 peran, satu tiket dibuat, satu lampiran diunduh, health/readiness OK, header keamanan & cookie produksi benar.
5. Production gagal start bila konfigurasi tidak aman (sudah dari TASK-002 — verifikasi nyata).

**Acceptance criteria:**
- [ ] Hanya proxy yang menerima trafik publik; DB/volume tidak terekspos (bukti: port scan/uji akses langsung).
- [ ] HTTPS + cookie `__Host-` + HSTS aktif; SPA deep-link (refresh `/tickets/:id`) berfungsi.
- [ ] Smoke test tiga peran lulus di staging; readiness probe benar.
- [ ] Urutan rilis & rollback terdokumentasi dan sekali dijalankan.

**Verifikasi & bukti selesai:** log smoke test staging, tangkapan konfigurasi proxy (tanpa secret), catatan rilis.

---

### TASK-051 — Backup + restore drill (TC-33)

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M6 / P0 / todo |
| DEV | DEV-14 |
| Referensi | NFR-06 (RPO 24 jam, RTO 4 jam — target diuji); TC-33; OPERATIONS §5 (backup konsisten maintenance window), §6 (restore drill 7 langkah); TEST_PLAN §8 |
| Dependensi | TASK-035, TASK-050 |
| Estimasi | 3–5 jam |

**Tujuan:** prosedur backup satu set (DB dump + volume lampiran + manifest + checksum, terenkripsi, lokasi terpisah) benar-benar dijalankan, lalu restore drill ke lingkungan kosong dengan verifikasi lengkap.

**File rencana (usulan):** `deploy/backup.sh` (usulan), `docs/testing/restore-drill-report.md` (usulan).

**Langkah implementasi:**
1. Implementasikan prosedur OPERATIONS §5: hentikan penerimaan request → tunggu drain → stop API → logical dump PostgreSQL + salinan volume → manifest (waktu, versi app, versi migrasi) + checksum → enkripsi → simpan terpisah.
2. Restore drill OPERATIONS §6 di lingkungan kosong: pulihkan DB+file dengan owner/permission benar; **cabut sesi hasil restore** sebelum dibuka; verifikasi jumlah objek, FK, tiket lintas status, timeline, checksum sampel file; unduh lampiran terotorisasi; smoke test login 3 peran.
3. Catat waktu pemulihan vs RTO 4 jam dan selisih data vs RPO 24 jam; kegagalan & perbaikan dicatat.
4. Retensi baseline: 7 harian + 4 mingguan (OPERATIONS §5) — dijadwalkan (cron) dan diverifikasi sekali.

**Acceptance criteria:**
- [ ] Satu backup set lengkap (DB + file + manifest + checksum) dibuat dan tersimpan terenkripsi di lokasi terpisah.
- [ ] Restore ke lingkungan kosong lulus seluruh verifikasi TC-33 termasuk download lampiran terotorisasi dan checksum cocok.
- [ ] Sesi hasil restore dicabut sebelum lingkungan dibuka.
- [ ] Laporan mencatat durasi restore vs RTO dan langkah yang gagal/diperbaiki.

**Verifikasi & bukti selesai:** laporan restore drill (tanggal, lingkungan, hasil verifikasi, durasi).

---

### TASK-052 — Perintah maintenance & retensi

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M6 / P0 / todo |
| DEV | DEV-14 |
| Referensi | OPERATIONS §7 (tabel retensi, dry-run, log jumlah, batas batch, role terpisah); ARCHITECTURE §6.3 (cleanup 24 jam — sudah sebagian di TASK-035); SCHEMA §5 (maintenance role) |
| Dependensi | TASK-008, TASK-013, TASK-035 |
| Estimasi | 3–5 jam |

**Tujuan:** satu perintah maintenance (role DB terbatas, terpisah dari runtime) untuk: hapus sesi expired/revoked >7 hari, notifikasi dibaca >90 hari (batch), purge audit >180 hari (dengan review operator), cleanup staging/orphan (konsisten TASK-035) — semua dengan dry-run, log jumlah item, dan batas batch.

**File rencana (usulan):** `backend/cmd/admin/` subcommand `maintenance` (flags: `--dry-run`, `--batch-size`, target job), dokumentasi penjadwalan.

**Langkah implementasi:**
1. Implementasikan job per tabel retensi OPERATIONS §7; jangan menghapus event/komentar/tiket (BR-18 — hanya data operasional yang listed).
2. Purge audit memerlukan flag eksplisit + dry-run default; log jumlah item dihapus per job.
3. Credential maintenance role terpisah (SCHEMA §5; OPERATIONS §7).
4. Jadwalkan (cron/systemd timer) di staging; verifikasi satu siklus berjalan.
5. Tes: data uji berumur sesuai fixture → job menghapus hanya yang jatuh tempo; dry-run tidak mengubah data.

**Acceptance criteria:**
- [ ] Semua job retensi OPERATIONS §7 tersedia dengan dry-run dan batas batch.
- [ ] Tidak ada penghapusan tiket/komentar/event/audit tanpa flag eksplisit.
- [ ] Satu siklus terjadwal terbukti berjalan di staging dengan log jumlah item.

**Verifikasi & bukti selesai:** log eksekusi maintenance (dry-run + riil), query jumlah row sebelum/sesudah.

---

### TASK-053 — OpenAPI deliverable + rekonsiliasi API_SPEC

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M6 / P0 / todo |
| DEV | DEV-14 (ROADMAP §3: OpenAPI dihasilkan saat implementasi dan divalidasi terhadap API_SPEC) |
| Referensi | API_SPEC §10; ROADMAP §3; TEST_PLAN §8 (release gate: OpenAPI sesuai implementasi) |
| Dependensi | TASK-034, TASK-038, TASK-040 (seluruh permukaan API selesai) |
| Estimasi | 3–5 jam |

**Tujuan:** spesifikasi OpenAPI 3 machine-readable yang selaras dengan implementasi nyata; selisih terhadap API_SPEC.md didaftar dan direkonsiliasi (dokumen sumber diperbarui lewat keputusan tercatat, bukan diam-diam).

**File rencana (usulan):** `docs/api/openapi.yaml` (atau generated dari kode — putuskan pendekatan: hand-written vs swaggo/oapi-codegen; catat), tes validasi contoh request/response.

**Langkah implementasi:**
1. Pilih pendekatan generasi (catat sebagai keputusan); tulis/generate spec seluruh endpoint `/api/v1` termasuk envelope, error codes, pagination.
2. Validasi: contract test yang membandingkan respons nyata endpoint terhadap contoh OpenAPI (minimal endpoint kritis: auth, tickets, transitions, attachments, notifications, dashboard).
3. Daftar selisih API_SPEC.md ↔ implementasi → usulan pembaruan dokumen di GAPS_AND_DECISIONS.md; endpoint yang belum ada ditandai backlog, bukan dipresentasikan berjalan (ROADMAP §3).

**Acceptance criteria:**
- [ ] openapi.yaml tervalidasi (tool linter OpenAPI) dan mencakup seluruh endpoint yang diimplementasikan.
- [ ] Contract test contoh respons lulus untuk endpoint kritis.
- [ ] Selisih dengan API_SPEC.md terdokumentasi dengan rencana rekonsiliasi.

**Verifikasi & bukti selesai:** file OpenAPI + hasil linter + hasil contract test.

---

### TASK-054 — README aplikasi, panduan setup, paket demo portofolio

| Atribut | Nilai |
| --- | --- |
| Milestone / Prioritas / Status | M6 / P0 / todo |
| DEV | DEV-14 |
| Referensi | ROADMAP §7 (paket portofolio, cerita demo, klaim CV sesuai hasil); TEST_PLAN §8 (release gate); PRD §4.2 (demo dari seed); SECURITY §9 (data demo aman) |
| Dependensi | TASK-050, TASK-051, TASK-053 |
| Estimasi | 4–6 jam |

**Tujuan:** README aplikasi (berbeda dari README paket perencanaan): masalah, pengguna, keputusan stack, fitur yang **benar-benar selesai**, cara menjalankan (setup reproducible dari clone), akun demo aman, batasan, bukti pengujian; diagram arsitektur + ERD; screenshot responsif; video demo 3–5 menit mengikuti cerita ROADMAP §7.

**File rencana (usulan):** `README.md` (root aplikasi), `docs/setup.md`, `docs/demo/` (screenshot, skrip video), pembaruan `docs/planning/*` status task → done dengan bukti.

**Langkah implementasi:**
1. Tulis README aplikasi mengikuti struktur ROADMAP §7; setiap klaim fitur merujuk bukti (test/laporan) — tidak ada angka performa tanpa data.
2. Panduan setup: prasyarat versi, perintah migrasi, seed, bootstrap admin, run dev & build produksi; uji sendiri dari clone bersih (atau minta orang lain menguji).
3. Siapkan demo: seed sintetis aman (SECURITY §9), skrip cerita demo ROADMAP §7 termasuk satu kasus konflik claim/akses ditolak; rekam video.
4. Verifikasi checklist release gate TEST_PLAN §8 item per item dengan tautan bukti; item yang belum terpenuhi dicatat jujur.
5. Perbarui status seluruh task di TASKS.md dan ringkasan milestone di IMPLEMENTATION_PLAN.md.

**Acceptance criteria:**
- [ ] Orang lain (atau pengembang di mesin bersih) dapat menjalankan aplikasi dari panduan setup tanpa informasi tambahan.
- [ ] README tidak mengklaim fitur/hasil uji yang belum ada buktinya.
- [ ] Video demo menampilkan alur tiga peran dari seed, termasuk satu kasus konflik.
- [ ] Checklist release gate TEST_PLAN §8 terisi dengan bukti atau catatan risiko tertulis.

**Verifikasi & bukti selesai:** README + setup guide + link video + checklist release gate terisi.

---

## Daftar lanjutan (P1/P2 — bukan dependensi MVP)

Tidak dijadwalkan pada backlog ini (ROADMAP §5, PRD §5.2/5.3): SSE, SLA + worker/outbox, email, ekspor laporan, filter tersimpan, cursor pagination, SSO, multi-organisasi, inventaris aset, integrasi chat, mobile. Setiap item memerlukan desain tambahan (DDL/API/permission/test plan baru — SCHEMA §9) dan ADR sebelum dikerjakan. Usulan peningkatan pasca-MVP dievaluasi berdasarkan masalah pengguna atau nilai pembelajaran terukur (ROADMAP §5).
