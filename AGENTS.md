# AGENTS.md — Konvensi repository RANDesk

Aplikasi IT Helpdesk satu organisasi (proyek portofolio). Backend Go + Gin + PostgreSQL; frontend React + TypeScript + Vite; arsitektur modular monolith.

## Sumber kebenaran

- Spesifikasi kontrak: `RANDesk_Planning/` (read-only — PRD, RULES, DESIGN, ARCHITECTURE, SCHEMA, API_SPEC, SECURITY, TEST_PLAN, ROADMAP, OPERATIONS, DECISIONS, REFERENCES). **Jangan mengubah dokumen sumber tanpa keputusan tercatat.**
- Backlog kerja: `docs/planning/` (IMPLEMENTATION_PLAN.md, TASKS.md, TRACEABILITY.md, GAPS_AND_DECISIONS.md). Kerjakan task sesuai urutan dependensi; perbarui status task dengan bukti.
- Status task: `todo` / `blocked` / `needs_verification` / `done`. **`done` hanya dengan bukti** (hasil tes, log, artefak — TEST_PLAN §9). Tes yang belum dijalankan tidak boleh dilaporkan lulus.

## Toolchain (diverifikasi 26 September 2026 di mesin dev Windows)

| Tool | Versi terkunci | Catatan |
| --- | --- | --- |
| Go | 1.26.4 windows/amd64 | GOROOT `C:\Program Files\Go` |
| Node.js | v22.16.0 (LTS) | via Laragon |
| npm | 10.9.2 | |
| Git | 2.37.3.windows.1 | |
| PostgreSQL | 17.11 (service `postgresql-x64-17`, port 5432) | Installer native Windows EDB via winget, 28 Sep 2026; superuser `postgres` |
| golang-migrate | v4.20.1 (`C:\Users\ricoa\go\bin\migrate.exe`) | Keputusan D-02. **Wajib build tag**: `go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1` — tanpa tag, driver postgres tidak ikut |

Database dev: `randesk_dev` dan `randesk_test` (owner role `randesk_migrate`); role runtime API: `randesk_runtime` (grants penuh ditata di TASK-027). Bootstrap: `backend/scripts/dev-bootstrap.sql` + `backend/scripts/setup-task003.ps1`. Kredensial di `.env`/skrip adalah placeholder development — bukan secret produksi.

Versi library Go terkunci di `backend/go.mod` (keputusan D-04, diperbarui 28 September 2026):

| Library | Versi | Dipakai di |
| --- | --- | --- |
| github.com/gin-gonic/gin | v1.12.0 | HTTP router (TASK-002) |
| github.com/jackc/pgx/v5 | v5.11.0 | Driver PostgreSQL via database/sql (TASK-002) |
| golang.org/x/sys | v0.41.0 | Sinyal CTRL_BREAK di tes Windows (TASK-002) |

Versi library frontend terkunci di `frontend/package-lock.json` (keputusan D-04, diselesaikan saat TASK-005 29 September 2026):

| Library | Versi | Catatan |
| --- | --- | --- |
| react / react-dom | ^19.3 | SPA |
| react-router-dom | ^7.18 | Router (data/createBrowserRouter) |
| @tanstack/react-query | ^5.104 | Server state |
| vite | ^8.3 | Build tool + dev proxy `/api`, `/health` → `127.0.0.1:8080` |
| @vitejs/plugin-react | ^6.1 | |
| tailwindcss + @tailwindcss/vite | ^4.3 | Config CSS-first `@theme` (DESIGN §4), tanpa tailwind.config.js |
| typescript | ^6.0 | `strict` penuh; typecheck = `tsc --noEmit` |
| eslint + typescript-eslint + react-hooks/refresh | ^10 / ^8.70 | Flat config (`eslint.config.js`) |
| vitest (+ jsdom) | ^5 / ^29 | Unit test; `npm run test` |

Utilitas tanggal memakai `Intl` bawaan (zone `Asia/Makassar`) — tanpa dependensi date library. UUID frontend belum diperlukan (request ID dibuat backend).

## Perintah umum

Backend (dari `backend/`):

```powershell
go build ./...          # kompilasi
go vet ./...            # static analysis
gofmt -l .              # format (harus kosong)
go test ./...           # unit test
go test -race ./...     # test dengan race detector (wajib untuk tes konkurensi)
go run ./cmd/api        # jalankan API (butuh .env di root repo atau env yang di-export)
```

Frontend (dari `frontend/`, berlaku setelah TASK-005):

```powershell
npm ci
npm run dev
npm run build
npm run typecheck
npm run test
```

Migrasi database (berlaku setelah TASK-003):

```powershell
migrate -path backend/migrations -database "$env:DATABASE_URL" up
```

## Aturan kerja

1. **Secret**: jangan pernah commit `.env`, password, token, atau connection string berisi kredensial. `.env.example` hanya placeholder. Log tidak boleh mencatat body, cookie, token, atau secret (SECURITY §7).
2. **Migrasi immutable** (RULES §8): migrasi yang sudah dipakai tidak diedit; perubahan schema = migrasi bernomor baru. DDL mengikuti SCHEMA §4.
3. **Kontrak API**: ikuti API_SPEC (envelope, kode error, endpoint). Perubahan kontrak = usulan tercatat di `docs/planning/GAPS_AND_DECISIONS.md` dulu.
4. **Transaksi**: mutasi tiket menulis tiket + event + notifikasi dalam satu transaksi DB (BR-12); urutan lock sesuai ARCHITECTURE §6 (advisory → master → user by UUID → ticket).
5. **Tes integrasi** memakai PostgreSQL nyata — SQLite dilarang untuk uji lock/konkurensi (TEST_PLAN §1).
6. **Shell**: PowerShell di Windows — gunakan `;` bukan `&&` sebagai pemisah statement.
7. **Commit**: pesan singkat imperative (contoh: `TASK-002: add API skeleton with health endpoints`). Sertakan ID task.
8. Struktur folder mengikuti ARCHITECTURE §3; service tidak menerima `gin.Context` — gunakan `context.Context` + actor + DTO (ARCHITECTURE §4).
