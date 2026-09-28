# Implementation Plan — RANDesk MVP

Versi 1.0 · 26 September 2026 · Status: rencana kerja engineering (dokumen task, bukan kode)

Dokumen ini menerjemahkan paket perencanaan `RANDesk_Planning/` menjadi rencana pengerjaan MVP yang terurut dependensi. Detail task berada di [TASKS.md](TASKS.md); pemetaan requirement di [TRACEABILITY.md](TRACEABILITY.md); asumsi, konflik, dan keputusan terbuka di [GAPS_AND_DECISIONS.md](GAPS_AND_DECISIONS.md).

## 1. Kondisi repository saat ini (hasil audit)

Audit dilakukan 26 September 2026 dengan memeriksa isi workspace `d:\laragon\www\ITHelpDesk`.

| Item | Status | Bukti |
| --- | --- | --- |
| Dokumen perencanaan | **Tersedia lengkap** | 13 file Markdown di `RANDesk_Planning/`: README, PRD, RULES, DESIGN, ARCHITECTURE, SCHEMA, API_SPEC, SECURITY, TEST_PLAN, ROADMAP, OPERATIONS, DECISIONS, REFERENCES |
| Kode backend Go | **Belum ada** | Tidak ada `backend/`, `go.mod`, atau file `.go` di workspace |
| Kode frontend React | **Belum ada** | Tidak ada `frontend/`, `package.json`, atau file `.ts/.tsx` |
| Migrasi database | **Belum ada** | DDL hanya berupa rancangan di `SCHEMA.md` §4 |
| Test suite | **Belum ada** | TEST_PLAN.md §6 menyatakan perintah CI "belum dapat dijalankan dari paket dokumen ini" |
| Version control | **Belum ada** | Tidak ditemukan direktori `.git` di root workspace |
| AGENTS.md / konvensi repo | **Belum ada** | Pencarian glob `**/AGENTS.md` menghasilkan 0 file |
| README aplikasi | **Belum ada** | README.md yang ada adalah README paket perencanaan, bukan README aplikasi |
| CI/CD | **Belum ada** | Tidak ada `.github/`, `.gitlab-ci.yml`, atau konfigurasi CI lain |
| OpenAPI spec | **Belum ada** | API_SPEC.md §10 menyebut OpenAPI sebagai deliverable implementasi |
| `docs/planning/` | **Baru dibuat** | Keempat dokumen planning ini adalah file pertama di direktori tersebut |

**Fakta lingkungan (diverifikasi 26 September 2026):** workspace berada di dalam `d:\laragon\www` (Laragon — stack PHP/MySQL lokal). Go 1.26.4, Node v22.16.0, npm 10.9.2, dan git 2.37.3 **tersedia**; Docker, PostgreSQL, dan make **tidak terpasang**. Prasyarat tersisa TASK-001/003: instal PostgreSQL dev via installer native Windows (keputusan disetujui — lihat GAPS_AND_DECISIONS.md §4.2).

**Kesimpulan audit:** ini adalah proyek greenfield. Seluruh path kode yang dirujuk ARCHITECTURE.md §3 adalah *rencana*, bukan berkas yang ada. Tidak ada fitur yang boleh dianggap selesai; semua task berstatus `todo` kecuali TASK-048 `needs_verification` (lingkungan acuan NFR-03 belum terbukti). Tidak ada task `blocked` per 26 September 2026 — D-06 diputuskan dengan adopsi resmi walkthrough terstruktur (P-07).

## 2. Cakupan MVP (P0)

Sumber: PRD.md §5.1, README paket §4.

**Masuk MVP:** login/logout/change-password, manajemen akun oleh admin, kategori & departemen, pembuatan tiket, assignment/claim, lifecycle (open → in_progress → resolved → closed, reopen), komentar publik & internal, lampiran (5 file, 5 MiB, JPEG/PNG/PDF), timeline event, notifikasi persisten dengan polling 30 detik, dashboard metrik, pencarian sederhana, audit administratif, pengujian akses & transaksi, deployment satu instance dengan TLS dan backup/restore.

**Tidak masuk MVP (jangan jadikan dependensi):** SSE, SLA/worker/outbox, email, ekspor laporan, filter tersimpan (P1); SSO, multi-organisasi, inventaris aset, mobile (P2); registrasi publik, reset password mandiri, rich text, lampiran internal, hard delete (di luar baseline).

**Alur utuh yang menjadi tulang punggung rencana** (README paket §6, ROADMAP §1): login → membuat tiket → assignment → penanganan → resolve → close, diuji sejak awal, bukan di akhir. Infrastruktur event/notifikasi transaksional dibangun bersama mutasi tiket (M2) karena BR-12/NFR-02 mewajibkan satu transaksi; UI notifikasi menyusul di M4.

## 3. Strategi pengerjaan

Disesuaikan untuk pengembang individual berlatar PHP/Laravel yang sedang mempelajari Go dan React (RULES.md §1):

1. **Vertical slice dulu, baru melebar** (ROADMAP §1). M0 membangun fondasi dua sisi (backend + frontend) yang langsung tersambung lewat health endpoint, bukan menyelesaikan seluruh backend dahulu.
2. **Policy otorisasi terpusat** (ARCHITECTURE §4). Satu modul policy dipakai detail, list, search, dashboard, dan download — mencegah risiko "otorisasi tersebar" (ROADMAP §6).
3. **Tes mengikuti setiap milestone**, bukan ditunda ke M5 (ROADMAP §2). Task backend menyertakan test TC kritis yang relevan sebagai acceptance criteria; suite negatif menyeluruh dan load test tetap di M5.
4. **Analogi Laravel → Go** untuk mempercepat belajar: middleware Gin ≈ middleware Laravel; service layer ≈ Form Request + Service class; `sql.Tx` ≈ `DB::transaction()` tetapi eksplisit; DTO eksplisit menggantikan Eloquent model serialization (RULES §8 melarang serialisasi struct DB langsung).
5. **Frontend: TanStack Query sebagai sumber kebenaran server state** (ARCHITECTURE §7); query key menyertakan user ID; cache dibersihkan saat sesi berakhir.
6. **Keputusan rutin diambil dan dicatat** sebagai asumsi di GAPS_AND_DECISIONS.md; hanya konflik kontrak antar-dokumen yang menghentikan pekerjaan bagian terpengaruh (README paket §5).
7. **Urutan lock database dipatuhi ketat**: advisory lock administrasi → master data → user (UUID) → ticket (ARCHITECTURE §6.2, RULES §5). Setiap task transaksi merujuk urutan ini.

## 4. Milestone, exit criteria, dan estimasi

Estimasi per milestone adalah penjumlahan rentang task penyusunnya (bottom-up, effort fokus termasuk tes per task). Pemetaan DEV-xx mengikuti ROADMAP.md §3.

| Milestone | Task | DEV terkait | Estimasi bottom-up | Exit criteria (bukti selesai) |
| --- | --- | --- | --- | --- |
| **M0 — Fondasi** | TASK-001…007 | DEV-01, DEV-02, DEV-03 | 23–36 jam | DB kosong dapat dimigrasi dan diverifikasi; API health merespons; UI shell memanggil API lewat proxy; CI menjalankan lint/build/test; error envelope konsisten; versi dependency terkunci di go.mod/lockfile (keputusan D-04 tercatat) |
| **M1 — Akun & sesi** | TASK-008…019 | DEV-04, DEV-05 | 43–67 jam | Login 3 peran berfungsi; TC-01/02/03/28 lulus; bootstrap admin dapat dijalankan; kelola akun/master data/audit admin berfungsi; TC-24/25/26/29 lulus; seed sintetis tersedia; UI admin selesai |
| **M2 — Tiket inti** | TASK-020…030 | DEV-06, DEV-07, DEV-08 | 39–60 jam | Alur tiket create → assign/claim → start → resolve → close → reopen berfungsi via API dan UI; TC-04/05/07/08/09/10/11/12/13/14 (mutasi closed) lulus; event+notifikasi tertulis satu transaksi; grants runtime append-only diterapkan |
| **M3 — Kolaborasi & berkas** | TASK-031…037 | DEV-09, DEV-10 | 25–39 jam | Komentar publik/internal dengan isolasi requester terbukti (TC-15/19); upload/download/delete dengan batas, kompensasi, dan cleanup berfungsi (TC-16/17/18/30); UI percakapan, timeline, dan lampiran selesai |
| **M4 — Notifikasi & dashboard** | TASK-038…041 | DEV-11, DEV-12 | 15–22 jam | Penerima notifikasi sesuai RULES §6 (TC-20/21); polling 30 detik dengan lifecycle tab/offline benar (TC-22); metrik dashboard cocok dengan fixture dan definisi PRD §9 (TC-23) |
| **M5 — Hardening & UX** | TASK-042…049 | DEV-13 | 28–44 jam | Gate keamanan SECURITY §10 lulus; suite konkurensi/failure-injection lulus berulang; a11y & responsif diverifikasi (TC-31); E2E TC-35 lulus; load test TC-32 dijalankan dengan hasil nyata (atau catatan deviasi lingkungan); UAT/walkthrough terstruktur tercatat |
| **M6 — Deploy & portofolio** | TASK-050…054 | DEV-14 | 17–27 jam | Deployment staging/produksi dapat dipulihkan; restore drill TC-33 lulus; OpenAPI selaras API_SPEC; README aplikasi + panduan setup dapat direproduksi orang lain; release gate TEST_PLAN §8 tercentang dengan bukti |

**Total bottom-up: ±190–295 jam.** ROADMAP §1 memperkirakan 82–122 jam inti (100–160 jam praktis). Selisih ini adalah **konflik estimasi yang disadari**, bukan kesalahan hitung — lihat GAPS_AND_DECISIONS.md GAP-02. ROADMAP sendiri mengamanatkan evaluasi ulang setelah milestone kedua; rencana ini memakai angka bottom-up sebagai baseline kerja dan merevisi ROADMAP secara resmi setelah M2.

## 5. Dependensi utama dan jalur kritis

### 5.1 Jalur kritis (urutan terpanjang)

```
TASK-001 (repo)
  → TASK-002 (backend skeleton) → TASK-004 (platform/error)
    → TASK-008 (password+session) → TASK-009 (auth+CSRF middleware)
      → TASK-013 (audit infra) → TASK-014 (users) → TASK-015 (activation guards)
        → TASK-020 (policy) → TASK-022 (list/search) → TASK-023 (create)
          → TASK-025 (assignment) → TASK-026 (transitions)
            → TASK-031 (comments) / TASK-034 (attachments)
              → TASK-036/037 (UI kolaborasi) → TASK-047 (E2E)
                → TASK-050 (deploy) → TASK-051 (restore) → TASK-054 (portfolio)
```

TASK-021 (event/notification writer) berada tepat di samping jalur kritis: ia wajib selesai **sebelum** TASK-023 karena create-ticket harus menulis event+notifikasi dalam transaksi yang sama (BR-12, TC-07).

### 5.2 Dependensi struktural penting

| Dependensi | Alasan |
| --- | --- |
| TASK-003 (migrasi) sebelum semua task backend ber-DB | Seluruh service bergantung pada schema SCHEMA.md §4 |
| TASK-021 (event+notification writer) sebelum TASK-023…026, 031, 034 | BR-12/NFR-02: mutasi inti + event + notifikasi satu transaksi; ROADMAP §2: "infrastruktur notifikasi disiapkan saat M2" |
| TASK-020 (policy pusat) sebelum TASK-022 | FR-05/BR-16: scope akses diterapkan sebelum pencarian/count/pagination |
| TASK-015 (guard deaktivasi) sebelum TASK-025 | TC-24: claim bersamaan dengan deaktivasi teknisi harus diuji setelah kedua sisi ada |
| TASK-033 (storage adapter) sebelum TASK-034 | ARCHITECTURE §4: bisnis tidak bergantung filesystem; failure injection di boundary file move (TEST_PLAN §5) |
| TASK-017 (seed/fixture) sebelum TASK-040…048 | TEST_PLAN §2: fixture minimum untuk uji metrik, scope, dan konkurensi |
| TASK-012 (frontend auth) sebelum semua task UI | Semua halaman memerlukan sesi, CSRF token, dan route guard |
| TASK-042 (rate limit/header) sebelum TASK-050 | SECURITY §7/§8 bagian dari gate deployment produksi |

### 5.3 Pekerjaan yang dapat diparalelkan

- Frontend M0 (TASK-005/006) paralel dengan backend M0 (TASK-002/004).
- Admin UI (TASK-018/019) paralel dengan awal M2 backend (TASK-020/021).
- TASK-032 (events API), TASK-027 (grants), dan TASK-038 (notification API) adalah task kecil yang dapat disisipkan.

## 6. Aturan execution

1. **Status task:** `todo` (belum dimulai), `blocked` (terhalang dependensi/keputusan), `needs_verification` (perlu bukti lingkungan/hasil), `done` (hanya dengan bukti: hasil tes, log eksekusi, atau artefak yang dapat diperiksa — sesuai README paket dan TEST_PLAN §9). Tes yang direncanakan tidak pernah dilaporkan lulus sebelum dijalankan.
2. **Definisi siap/selesai** mengikuti ROADMAP §4: task siap bila requirement, peran, input/output, validasi, dan skenario gagal jelas; selesai bila acceptance criteria terbukti, UI menangani loading/empty/error, tes bermakna lulus, log bebas secret, dan dokumen terdampak diperbarui.
3. **Perubahan kontrak** (API_SPEC/SCHEMA/RULES) tidak dilakukan diam-diam: catat sebagai usulan di GAPS_AND_DECISIONS.md dengan alasan dan dampak, lalu perbarui semua dokumen pemilik dalam perubahan yang sama (README paket §5).
4. **Migrasi yang sudah dirilis immutable** (RULES §8); perubahan schema lewat migrasi baru.
5. **Ukuran task** 2–6 jam fokus; task dengan estimasi lebih besar sudah dipecah. Ketidakpastian tertinggi ada pada TASK-008/009/021/025/034 (area baru bagi pengembang: crypto/session, transaksi konkuren Go) — rentang estimasinya sengaja lebar.
6. **Fitur P1/P2** dicatat di ROADMAP §5 dan tidak boleh masuk sprint MVP; usulan fitur mengikuti RULES §9.

## 7. Tiga task pertama sesuai urutan dependensi

1. **TASK-001** — Inisialisasi repository, struktur folder, dan konfigurasi lingkungan (mengunci D-04 versi dependency; prasyarat segala pekerjaan).
2. **TASK-002** — Skeleton backend Go: `cmd/api`, konfigurasi, Gin, health endpoint, graceful shutdown (paralel: TASK-003 migrasi dan TASK-005 frontend skeleton).
3. **TASK-003** — PostgreSQL development + tooling migrasi + migrasi baseline DDL SCHEMA §4.

## 8. Rujukan dokumen

| Keputusan | Sumber |
| --- | --- |
| Scope & acceptance | `RANDesk_Planning/PRD.md` (FR-01…13, NFR-01…08, definisi metrik §9) |
| Izin & aturan bisnis | `RANDesk_Planning/RULES.md` (matriks §2, lifecycle §3, BR-01…18 §4, konkurensi §5, notifikasi §6, validasi §7) |
| Perilaku UI | `RANDesk_Planning/DESIGN.md` (D-01…D-07, tokens §4, layout §5, a11y §8) |
| Pembagian sistem & transaksi | `RANDesk_Planning/ARCHITECTURE.md` (§3 struktur, §5 siklus request, §6 transaksi, §8 notifikasi) |
| Penyimpanan | `RANDesk_Planning/SCHEMA.md` (DDL §4, invariant §5, payload event §6, query §7, seed §8) |
| Kontrak HTTP | `RANDesk_Planning/API_SPEC.md` (envelope §2, error §3, endpoint §4–§9) |
| Keamanan | `RANDesk_Planning/SECURITY.md` (§2 password, §3 sesi, §4 CSRF, §6 lampiran, §7 rate limit, §10 gate) |
| Urutan & verifikasi | `RANDesk_Planning/ROADMAP.md` (milestone §2, DEV-01…14 §3), `RANDesk_Planning/TEST_PLAN.md` (TC-01…35 §3, traceability §4, CI §6, release gate §8) |
| Operasi | `RANDesk_Planning/OPERATIONS.md` (bootstrap §3, backup §5, restore §6, retensi §7) |
| Keputusan & asumsi | `RANDesk_Planning/DECISIONS.md` (ADR-001…012, A-01…07, D-01…06) |
