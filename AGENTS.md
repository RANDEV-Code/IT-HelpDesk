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
| PostgreSQL | *(diisi saat TASK-003 — installer native Windows EDB)* | Docker tidak tersedia di mesin dev |
| golang-migrate | *(diisi saat TASK-003 — keputusan D-02)* | |

Versi library Go (Gin, pgx) dan frontend (React, Vite, Tailwind, TanStack Query) dikunci di `backend/go.mod` dan `frontend/package-lock.json` saat task terkait dikerjakan (keputusan D-04) — catat versinya di tabel ini.

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
