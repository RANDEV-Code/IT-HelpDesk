# Gaps, Conflicts, Assumptions & Decisions — RANDesk MVP Planning

Versi 1.0 · 26 September 2026 · Pendamping [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md), [TASKS.md](TASKS.md), [TRACEABILITY.md](TRACEABILITY.md).

Dokumen ini memisahkan **fakta** (hasil audit repository), **konflik spesifikasi** (antar dokumen sumber), **asumsi perencanaan** (ditetapkan di tahap planning ini), **usulan** (perlu disetujui/dikunci), dan **risiko**. Sesuai batasan tahap ini: dokumen sumber di `RANDesk_Planning/` tidak diubah dan tidak ada kode aplikasi yang dibuat; usulan perubahan kontrak dicatat di sini.

## 1. Fakta hasil audit repository (26 September 2026)

- Workspace `d:\laragon\www\ITHelpDesk` hanya berisi `RANDesk_Planning/` (13 dokumen Markdown, versi 1.0, tanggal 26 September 2026) dan `docs/planning/` (empat dokumen ini).
- Tidak ada `.git`, tidak ada `AGENTS.md`, tidak ada kode backend/frontend, tidak ada migrasi, tidak ada CI, tidak ada lockfile/go.mod, tidak ada README aplikasi, tidak ada OpenAPI.
- Workspace berada di dalam tree Laragon (stack PHP/MySQL Windows). **Toolchain diverifikasi 26 September 2026** (perintah `go version`, `node --version`, `Get-Command`): Go 1.26.4 windows/amd64 (GOROOT `C:\Program Files\Go`), Node v22.16.0, npm 10.9.2, git 2.37.3 **tersedia**; Docker, PostgreSQL, dan make **tidak terpasang** (`d:\laragon\bin` hanya memuat MySQL/MariaDB, tidak ada folder postgresql). Sisa prasyarat TASK-001/003: instal PostgreSQL dev (keputusan: installer native Windows — lihat §4.2).

## 2. Dokumen yang belum tersedia (deliverable, bukan gap spesifikasi)

| Dokumen | Direncanakan pada | Keterangan |
| --- | --- | --- |
| OpenAPI machine-readable | TASK-053 | API_SPEC §10: deliverable implementasi; selisih dengan API_SPEC.md harus direkonsiliasi |
| README aplikasi + panduan setup | TASK-054 | TEST_PLAN §8 release gate; harus dibedakan dari README paket perencanaan |
| Laporan hasil pengujian (per run) | Setiap task uji; format TEST_PLAN §9 | Belum ada satu pun hasil uji — tidak boleh ada klaim lulus saat ini |
| Laporan load test, a11y, UAT, restore drill | TASK-048, 046, 049, 051 | Template laporan mengikuti TEST_PLAN §9 |
| go.mod / lockfile berisi versi terkunci | TASK-001 | Keputusan D-04 |
| AGENTS.md / konvensi repo | TASK-001 | Belum ada; usulan isi: perintah build/test, konvensi commit, larangan secret |

## 3. Ketidaksesuaian antar dokumen sumber (konflik)

Sesuai README paket §5: konflik ditampilkan beserta file dan bagian; **tidak ada dokumen yang dianggap otomatis direvisi**. Tidak ditemukan konflik kontrak bisnis yang memblokir alur utama; dua butir berikut adalah ketidaksesuaian perencanaan yang perlu keputusan.

### GAP-01 — Lokasi paket dokumen vs struktur repo rencana
- **Fakta:** ARCHITECTURE.md §3 merencanakan `docs/` sebagai lokasi "Paket dokumen Markdown ini"; pada repository nyata paket berada di `RANDesk_Planning/`.
- **Dampak:** referensi silang dokumen vs path repo; posisi README aplikasi masa depan (root) agar tidak rancu dengan README paket.
- **Kategori:** dapat ditunda (tidak memblokir task teknis).
- **Usulan:** pertahankan `RANDesk_Planning/` sebagai sumber read-only pada tahap ini; saat TASK-001, putuskan apakah paket dipindah ke `docs/` (memperbarui ARCHITECTURE §3 lewat keputusan tercatat) atau path rencana direvisi. Rekomendasi: pindahkan ke `docs/planning/randesk/` **hanya bila** disetujui pemilik proyek, karena tahap ini melarang perubahan dokumen sumber.
- **Task terpengaruh:** TASK-001, TASK-054.

### GAP-02 — Estimasi usaha: top-down ROADMAP vs bottom-up backlog
- **Fakta:** ROADMAP.md §1 memperkirakan 82–122 jam inti (100–160 jam praktis, 7–11 minggu @15 jam/minggu). Penjumlahan rentang 54 task di TASKS.md menghasilkan ±190–295 jam — sekitar 1,8× batas atas ROADMAP.
- **Analisis:** ROADMAP adalah estimasi perencanaan awal; backlog bottom-up memasukkan eksplisit tes per task, CI, hardening, dan buffer belajar Go/React yang oleh ROADMAP dihitung sebagai tambahan 20–30% di luar "pekerjaan inti". Sebagian selisih adalah perbedaan basis hitung, sebagian kemungkinan realisme bottom-up.
- **Kategori:** dapat ditunda — tidak memblokir M0/M1.
- **Rekomendasi:** pakai rentang bottom-up sebagai baseline kerja; jalankan evaluasi ulang resmi setelah M2 sebagaimana diamanatkan ROADMAP §1 ("evaluasi setelah milestone kedua"), lalu revisi ROADMAP/DECISIONS lewat perubahan tercatat. Jangan merevisi angka ROADMAP sekarang tanpa bukti kecepatan nyata.
- **Task terpengaruh:** semua (perencanaan), secara formal none blocked.

## 4. Keputusan yang belum diputuskan

### 4.1 Dari DECISIONS.md §3 (sumber)

| ID | Item | Batas waktu sumber | Status planning | Task terpengaruh |
| --- | --- | --- | --- | --- |
| D-01 | Nama final produk & domain | Sebelum publikasi | Terbuka; "RANDesk" tetap dipakai sebagai nama kerja (tidak memblokir apa pun — DECISIONS §3) | TASK-050, TASK-054 |
| D-02 | Tool migrasi | Sebelum migrasi pertama | **DIPUTUSKAN 26 Sep 2026: golang-migrate** (disetujui pemilik proyek); eksekusi instalasi di TASK-003 | TASK-003 |
| D-03 | Reverse proxy & hosting | Sebelum milestone deployment | **Usulan default: Caddy pada satu VM Linux/container**; dikunci saat TASK-050. Catatan: ROADMAP §6 menyarankan staging sejak M2 — bila VM tersedia lebih awal, gunakan lebih awal | TASK-050, TASK-051 |
| D-04 | Versi dependency tepat | Milestone fondasi | Terbuka; dikunci di TASK-001 setelah verifikasi toolchain (GAP-06) | TASK-001…007 |
| D-05 | Retensi data organisasi | Sebelum data riil masuk | Terbuka; baseline demo OPERATIONS §7 dipakai (TASK-052). Tidak memblokir MVP (data sintetis — A-07) | TASK-052 |
| D-06 | Calon peserta UAT (5 orang) | Sebelum tahap UAT | **DIPUTUSKAN 26 Sep 2026: adopsi resmi alternatif walkthrough terstruktur** (TEST_PLAN §7 / P-07); varian UAT 5 peserta dilepas dari cakupan MVP. TASK-049 tidak lagi `blocked` (status `todo`); laporan wajib mencatat hasil belum mewakili pengguna nyata | TASK-049 |

### 4.2 Keputusan teknis rutin (**disetujui pemilik proyek 26 September 2026** — dikunci saat kickoff task terkait dan dicatat di AGENTS.md)

Tambahan keputusan lingkungan hasil verifikasi toolchain: PostgreSQL dev memakai **installer native Windows (EDB)** karena Docker tidak terpasang di mesin dev; CI memakai **service container PostgreSQL GitHub Actions** sehingga Docker lokal tidak diperlukan; `make` tidak tersedia — TASK-001 memakai skrip PowerShell atau Taskfile opsional sebagai pengganti Makefile.

| Usulan | Default yang disarankan | Alasan | Task |
| --- | --- | --- | --- |
| Platform CI | GitHub Actions | Gratis untuk repo publik portofolio, mendukung service container PostgreSQL & race detector | TASK-007 |
| DB tes integrasi | PostgreSQL nyata via Docker service container / instance lokal kedua | TEST_PLAN §1 melarang SQLite untuk uji lock; TEST_PLAN §6 memerlukan DB terisolasi per run | TASK-003, TASK-007 |
| Framework E2E | Playwright | Multi-browser, artefak screenshot (TEST_PLAN §6), webServer terintegrasi | TASK-047 |
| Tool load test | k6 (atau setara) | Skenario pacing 1 req/user/detik mudah diekspresikan (PRD §8) | TASK-048 |
| Audit a11y otomatis | axe-core (via Playwright/vitest) | NFR-04 mensyaratkan audit otomatis + manual | TASK-046 |
| Library Argon2id | Library Go terawat sesuai SECURITY §2 (pilihan tepat dikunci di D-04) | SECURITY §2 tidak menyebut nama library; parameter sudah ditetapkan | TASK-008 |
| Frontend unit/component test | Vitest + Testing Library | Selaras ekosistem Vite | TASK-005…007 |
| Lokasi grants runtime | Script SQL terpisah yang dijalankan operator (bukan migrasi schema) | Grants adalah konfigurasi keamanan lingkungan, bukan struktur; SCHEMA §5 memisah role migrasi vs runtime — bentuk final diputuskan saat TASK-027 | TASK-027 |

Tidak satu pun usulan di atas mengubah kontrak bisnis (API_SPEC/RULES/SCHEMA); semuanya berada dalam ruang pilihan yang sengaja dibiarkan terbuka oleh dokumen sumber.

## 5. Asumsi perencanaan yang ditetapkan pada tahap ini

| ID | Asumsi | Dasar | Dampak jika salah |
| --- | --- | --- | --- |
| P-01 | Nama kerja "RANDesk" dan cookie prefix `randesk_` dipertahankan sampai D-01 diputuskan | DECISIONS §3 ("jangan menunda alur utama karena nama") | Rename kosmetik di M6 (TASK-050/054) |
| P-02 | Struktur modul backend mengikuti ARCHITECTURE §3 apa adanya; modul `internal/activity` tidak dibuat — penulisan event/notifikasi hidup di `internal/notifications` + `internal/tickets` | ARCHITECTURE §3 tidak menyebut `internal/activity`; penamaan final diserahkan implementasi | Refactor penamaan kecil, tanpa dampak kontrak |
| P-03 | `GET /assignees` diimplementasikan di modul users (TASK-014), bukan master data | API_SPEC §5 mengelompokkannya bersama akun; data berasal dari tabel users | Hanya penempatan kode |
| P-04 | Satu instalasi demo memakai dataset sintetis; tidak ada data pribadi nyata sampai D-05 diputuskan | SECURITY §9, A-07 | TASK-051/052 harus direvisi sebelum data riil |
| P-05 | Load test (TASK-048) kemungkinan dijalankan pada lingkungan dev yang lebih kecil dari acuan NFR-03, dan hasilnya dicatat sebagai deviasi — bukan pemenuhan target | PRD §8 melarang mengubah target tanpa profil; GAP-04 | Klaim NFR-03 tidak dapat dibuat; release gate TEST_PLAN §8 butir performa dicatat apa adanya |
| P-06 | Deployment target M6 adalah satu VM Linux atau container lokal yang setara; bila belum tersedia, TASK-050/051 dijalankan pada lingkungan lokal terisolasi dan dicatat sebagai staging sementara | OPERATIONS §1 (satu VM Linux + persistent volume) | Bukti TC-33 perlu diulang di lingkungan target sebenarnya |
| P-07 | UAT penuh (5 peserta) tidak tersedia sebelum acceptance MVP; walkthrough terstruktur oleh pengembang menjadi bukti pengganti dengan catatan keterbatasan | TEST_PLAN §7 secara eksplisit mengizinkan; **diadopsi resmi oleh pemilik proyek 26 Sep 2026 (keputusan D-06)** | Target PRD §4.2 baris 1–2 tidak terukur; dicatat di laporan |

## 6. Risiko

| ID | Risiko | Kemungkinan/Dampak | Mitigasi di rencana |
| --- | --- | --- | --- |
| R-01 | Kurva belajar Go (goroutine/tx/context) menyebabkan task M2 molor, terutama TASK-021/025 | Sedang / Tinggi (jalur kritis) | Task kecil 2–6 jam; TASK-021 dijadwalkan sebelum mutasi apa pun; analogi Laravel di IMPLEMENTATION_PLAN §3; tes table-driven memaksa pemahaman |
| R-02 | ~~Toolchain belum terverifikasi~~ **Sebagian terselesaikan 26 Sep 2026**: Go/Node/git terbukti ada; Docker & PostgreSQL tidak terpasang | Menurun / Sedang | Sisa mitigasi: instal PostgreSQL native Windows di TASK-001 (keputusan §4.2); CI tidak bergantung Docker lokal (service container GitHub Actions) |
| R-03 | Tes konkurensi flaky memberi rasa aman palsu | Sedang / Tinggi | TEST_PLAN §5: verifikasi state akhir DB (bukan hanya kode respons), run berulang ≥20x, race detector di CI (TASK-007/045) |
| R-04 | Deviasi Windows (dev) vs Linux (deploy): path storage, cookie `__Host-`, proxy headers | Sedang / Sedang | ROADMAP §6: staging sejak dini (TASK-050 dapat dimajukan setelah M2 bila VM tersedia); smoke test produksi-like |
| R-05 | Scope creep P1 (SSE/SLA/email) masuk sebelum MVP stabil | Sedang / Sedang | Daftar lanjutan TASKS.md; RULES §9 mensyaratkan proposal tertulis; IMPLEMENTATION_PLAN §6 butir 6 |
| R-06 | Lingkungan load test tidak memenuhi acuan NFR-03 | Tinggi / Sedang | P-05: catat deviasi jujur; profilkan sebelum mengusulkan perubahan target (PRD §8) |
| R-07 | Kompensasi filesystem (orphan attachment) salah implementasi → data tak konsisten | Sedang / Tinggi | Failure injection di boundary file move/commit sejak TASK-034 (bukan ditunda ke M5); TASK-035 cleanup teruji; TASK-045 mengulang |
| R-08 | Klaim selesai tanpa bukti (status `done` prematur) | Sedang / Sedang | Aturan status IMPLEMENTATION_PLAN §6: `done` hanya dengan bukti; TEST_PLAN §9 template laporan |

## 7. Blocker vs dapat ditunda — ringkasan

| Item | Klasifikasi | Efek saat ini |
| --- | --- | --- |
| GAP-06/R-02: toolchain & PostgreSQL dev | **Terselesaikan sebagian (26 Sep 2026)** | Go/Node/npm/git terverifikasi ada; sisa pekerjaan: instal PostgreSQL native Windows di TASK-001 — bukan lagi blocker tak terduga |
| D-06: peserta UAT | **Selesai — diputuskan 26 Sep 2026** | Walkthrough terstruktur (P-07) diadopsi resmi; TASK-049 `todo`, tidak ada task yang blocked |
| GAP-04/P-05: lingkungan acuan NFR-03 | **needs_verification** | TASK-048 tetap dijadwalkan; klaim pemenuhan target ditahan sampai lingkungan terbukti |
| GAP-01: lokasi dokumen | Dapat ditunda | Tidak memblokir task teknis; diputuskan di TASK-001/054 |
| GAP-02: selisih estimasi | Dapat ditunda | Re-baseline resmi setelah M2 |
| D-01 nama/domain, D-05 retensi data riil | Dapat ditunda | Batas waktu sumber: sebelum publikasi / sebelum data riil |
| GAP-08 (di bawah): monitoring berkelanjutan | Dapat ditunda | Di luar cakupan MVP demo |

### GAP-04 — Lingkungan acuan performa (detail)
NFR-03/PRD §8 mensyaratkan host acuan 2 vCPU/4 GiB/SSD dengan generator beban di host terpisah. Tidak ada bukti lingkungan semacam itu tersedia (fakta audit §1). TASK-048 ditandai `needs_verification`; bila hanya lingkungan dev yang ada, jalankan pengukuran, catat spesifikasi host sebenarnya, dan tandai hasilnya sebagai **indikatif**, bukan pemenuhan NFR-03.

### GAP-08 — Monitoring/alerting berkelanjutan belum memiliki task khusus
OPERATIONS §8 mendefinisikan sinyal dan ambang (5xx, readiness, disk, pool DB, orphan, backup), tetapi tidak ada DEV/TC yang menuntut sistem monitoring aktif pada MVP; TEST_PLAN hanya mensyaratkan health endpoint + log (NFR-08, tercakup TASK-002/004/043). Untuk demo/pilot terbatas ini dapat ditunda; bila MVP dipromosikan ke pilot nyata, tambahkan backlog item monitoring (log shipping/alert sederhana) sebelum menerima pengguna riil. Tidak memblokir task mana pun.

### Observasi minor (bukan konflik, dicatat agar tidak hilang)
1. **Metrik "Menunggu konfirmasi" (PRD §9)** tidak memiliki field eksplisit di payload `GET /dashboard/summary` (API_SPEC §9), tetapi dapat diturunkan dari `status_counts.resolved`. TASK-040/041 memakai penurunan ini; jika pemilik proyek menginginkan field eksplisit, itu usulan perubahan API_SPEC (dampak kecil, non-breaking bila menambah field).
2. **RULES §2: teknisi bukan-assignee boleh membaca untuk koordinasi** — UI (TASK-030/036) harus memastikan composer komentar tidak dirender untuk teknisi non-assignee, karena backend menolak (403). Sudah dicakup acceptance criteria TASK-030 (render dari allowed_actions).
3. **`updated_at` default `CURRENT_TIMESTAMP` pada DDL (SCHEMA §4)** vs semantik "perubahan inti terakhir" (RULES §5/SCHEMA §3): service harus menyetel `updated_at` eksplisit pada mutasi inti dan **tidak** mengubahnya pada operasi child (komentar/lampiran). Dicatat sebagai perhatian implementasi TASK-024/031/034; bukan perubahan schema.
4. **Verifikasi privilege sequence/identity untuk runtime role** (SCHEMA §5) dijadwalkan eksplisit di TASK-003 langkah 6 dan TASK-027.

## 8. Usulan perubahan kontrak/arsitektur (belum diputuskan, memerlukan persetujuan)

Belum ada usulan perubahan kontrak yang diajukan pada tahap ini. Semua task dirancang agar patuh pada API_SPEC/SCHEMA/RULES apa adanya. Satu-satunya kandidat usulan (payload eksplisit "Menunggu konfirmasi", observasi 1 di atas) bersifat opsional dan tidak boleh dikerjakan sebelum disetujui.

## 9. Rekomendasi

1. **Mulai dari TASK-001 dengan verifikasi toolchain sebagai gate pertama** (mengubah risiko terbesar hari pertama menjadi langkah eksplisit).
2. **Kunci D-02 (golang-migrate) dan D-04 (versi dependency) pada minggu pertama** agar CI dan migrasi tidak berubah-ubah.
3. **Bangun TASK-021 (event/notification writer) segera setelah auth selesai**, sebelum mutasi tiket apa pun — ini prasyarat BR-12 yang paling sering terlupa dan paling mahal bila dirakit belakangan.
4. **Jalankan failure-injection upload sejak TASK-034**, jangan tunda ke M5 (mitigasi R-07, sejalan TEST_PLAN §5).
5. **Setelah M2, lakukan re-baseline estimasi resmi** (GAP-02) dengan kecepatan nyata yang terukur; perbarui ROADMAP lewat perubahan tercatat.
6. ~~Ajukan keputusan D-06 sekarang~~ **Selesai 26 Sep 2026**: pemilik proyek mengadopsi P-07 (walkthrough terstruktur) sebagai pengganti UAT penuh; seluruh usulan §4.2 dan D-02 disetujui.
7. **Pertimbangkan memajukan TASK-050 (staging) ke akhir M2** bila VM/container target tersedia — ROADMAP §6 secara eksplisit menyarankan staging dini untuk menangkap deviasi hosting (R-04).

## 10. Catatan deviasi eksekusi (ditambahkan saat implementasi)

### TASK-003 — deviasi & temuan (27 September 2026)

1. **Pembungkus `BEGIN/COMMIT` dilepas dari migrasi baseline.** DDL SCHEMA §4 disalin apa adanya ke `000001_baseline.up.sql` **kecuali** pembungkus `BEGIN;`/`COMMIT;` dokumen sumber, karena golang-migrate sudah menjalankan setiap migrasi dalam satu transaksi; `BEGIN` bersarang akan error di PostgreSQL. Isi DDL (tabel, kolom, CHECK, indeks, komentar) tidak diubah — bukan perubahan kontrak.
2. **golang-migrate wajib di-build dengan tag `postgres`.** `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest` tanpa tag menghasilkan binary tanpa driver ("unknown driver postgres"). Perintah terkunci (v4.20.1): `go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1` — sudah dicatat di AGENTS.md.
3. **Urutan evaluasi CHECK constraint.** Untuk status tiket invalid, PostgreSQL menolak lewat constraint tabel `tickets_lifecycle_ck` lebih dulu sebelum `tickets_status_check` (CHECK dievaluasi tanpa urutan terjamin). Enum status tetap ditegakkan; skrip uji (`check-constraints-task003.sql`) mengasertakan nama `tickets_lifecycle_ck` dan mendokumentasikan alasannya. Tidak ada perubahan schema.
4. **Kolom identity tidak butuh grant sequence terpisah.** `tickets.ticket_no` dan `ticket_events.seq` memakai `GENERATED ALWAYS AS IDENTITY`, sehingga hak INSERT pada tabel sudah cukup untuk runtime role — tidak ada grant `USAGE/UPDATE ON SEQUENCE` tambahan (catatan SCHEMA §5 terpenuhi; grants penuh tetap dijadwalkan di TASK-027).
5. **`migrate down` interaktif.** Perintah `migrate down` tanpa argumen meminta konfirmasi y/N; skrip memakai `migrate down 1`. Setelah down, tabel bookkeeping `schema_migrations` sengaja tersisa (perilaku bawaan golang-migrate) — bukan residu DDL aplikasi.
6. **Versi final D-02/D-04 yang dieksekusi:** PostgreSQL 17.11 (installer native EDB, service `postgresql-x64-17`), golang-migrate v4.20.1. Password role dev adalah placeholder pengembangan (`*_dev`) dan hanya berlaku lokal — bukan secret produksi.

### TASK-004 — deviasi & temuan (28 September 2026)

1. **Timeout handler dipetakan ke 500 INTERNAL_ERROR, bukan 504.** TASKS.md menyebut "500/504 terkontrol", tetapi API_SPEC §3 tidak mendefinisikan code untuk 504 dan tidak memakai status itu. Agar tetap patuh kontrak, deadline handler yang terlampaui dipetakan ke 500 INTERNAL_ERROR generik (`apierr.From(context.DeadlineExceeded)`). Timeout query 3 detik dan handler 10 detik diterapkan lewat **context turunan** (kooperatif), bukan goroutine yang membatalkan handler secara paksa — sesuai ARCHITECTURE §5 ("hindari goroutine yang hidup tanpa pembatalan"). Konsekuensi: handler wajib menghormati `ctx`; backstop middleware menulis 500 hanya bila handler selesai tanpa menulis respons setelah deadline.
2. **Request ID memakai UUID v4 dari `crypto/rand` stdlib** (tanpa dependensi `google/uuid`). Request_id selalu dibuat server (API_SPEC §1), tidak mengambil dari header masuk, dan dikorelasikan ke header `X-Request-ID` + `meta.request_id` + field log.
3. **Batas body diterapkan dua lapis:** `Content-Length` melebihi batas → 413 sebelum handler; body streaming dibungkus `http.MaxBytesReader` sehingga `DecodeJSON` memetakan `*http.MaxBytesError` → 413 PAYLOAD_TOO_LARGE (bukan 400). Batas multipart 6 MiB sudah terpasang sebagai hook untuk M3.
4. **Route template pada log** memakai `c.FullPath()` (pola rute, mis. `/tickets/:id`), bukan path mentah; rute tak cocok dicatat sebagai `(unmatched)`. Tidak ada body/cookie/token/header sensitif yang dicatat (SECURITY §8) — dibuktikan `TestLogging_RequiredFieldsAndNoBody`.

### TASK-005 — deviasi & temuan (29 September 2026)

1. **Versi frontend = rilis stabil terbaru saat scaffold (D-04).** `npx create-vite` tidak dapat dijalankan andal di shell dev (npm/PowerShell menelan argumen `--template` + `npx` gagal resolve Node), sehingga proyek dibuat manual lewat `npm install` agar registry memilihkan stabil terkini: React 19.3, react-router 7.18, TanStack Query 5.104, Vite 8.3, plugin-react 6.1, Tailwind 4.3, TS 6.0, ESLint 10, vitest 5. Semua terkunci di `package-lock.json`.
2. **Tailwind v4 config-first.** Tokens DESIGN §4 didefinisikan sebagai variabel `@theme` di `src/index.css` (bukan `tailwind.config.js`) mengikuti cara Tailwind v4; tidak ada dependensi font eksternal (stack sistem sesuai DESIGN §4.2).
3. **Typecheck memakai `tsc --noEmit` satu `tsconfig.json`**, bukan project-references `tsc -b` bawaan template (tidak perlu `composite`). Flag `noUncheckedSideEffectImports` sengaja tidak diaktifkan karena menolak import sisi-efek `./index.css`; strict penuh + `erasableSyntaxOnly` tetap aktif.
4. **Bukti `/api` mencapai backend memakai proxy `/health`.** Karena endpoint `/api/v1/*` baru ada mulai TASK-008, demo acceptance memakai `GET /health/live` yang ikut di-proxy Vite; badge "Backend hidup: ok" + `X-Request-ID` membuktikan satu-origin SPA↔API. Proxy `/api` sudah dikonfigurasi untuk task berikutnya.
5. **Keterbatasan screenshot otomatis.** Bukti visual 4 viewport (1440/1024/768/360) terhambat lingkungan otomasi browser (tab `visibilityState=hidden` → `requestAnimationFrame` tidak fire → `take_screenshot` timeout; lebar window terkunci 976px). Perilaku responsif sebaliknya dibuktikan keras lewat pengukuran DOM/computed-style (sidebar `static` 240px ≥1024 vs `x=-240` tersembunyi + hamburger 44px <1024; drawer `aria-expanded` + backdrop) dan 1 screenshot nyata mode drawer. Bukan cacat aplikasi; operator manusia tetap dianjurkan memverifikasi visual.

### TASK-007 — deviasi, temuan, & status needs_verification (29 September 2026)

1. **Platform CI = GitHub Actions** (keputusan rutin §4.2 disetujui pemilik). Workflow `.github/workflows/ci.yml`: job `backend` (service `postgres:17`, gate gofmt/vet/build/test/race), `frontend` (`npm ci` + lint/typecheck/vitest/build), `security-scan` non-blocking (`govulncheck` + `npm audit`, temuan dikaji manual SECURITY §10). Tidak ada secret — hanya kredensial service container sementara.
2. **Migrasi di test memakai golang-migrate sebagai library** (bukan CLI) agar self-contained: sumber `iofs` dari `//go:embed *.sql` di paket `backend/migrations`, instance DB `postgres.WithInstance` di atas pool pgx `database/sql`. Konsekuensi: `github.com/golang-migrate/migrate/v4 v4.20.1` menjadi dependensi langsung Go dan menarik `github.com/lib/pq` (indirect) sebagai driver `postgres`. Prod tetap memakai CLI migrate; tidak ada perubahan skema.
3. **Helper testdb sengaja TIDAK memakai `DATABASE_URL`** (itu DB dev) — hanya `TEST_DATABASE_URL`; kosong → `t.Skip`. Ini mencegah tes integrasi merusak `randesk_dev` dan menjaga `go test ./...` tetap hijau tanpa DB.
4. **Status `needs_verification`, bukan `done`:** dua acceptance belum terbukti karena batasan lingkungan — (a) repo **belum punya remote GitHub** (`git remote -v` kosong) sehingga "CI hijau + URL run" belum bisa dihasilkan; (b) **`go test -race` tidak dapat jalan di Windows dev ini** (tidak ada gcc/cgo; `-race requires cgo`) — akan jalan di runner Linux CI. Semua gerbang lain (build, unit, integrasi contoh `TestIntegrationMigrateAndCRUD` lulus terhadap `randesk_test`, gate frontend) sudah hijau lokal. Untuk menutup task: inisial remote + push + tautkan log run CI hijau.
