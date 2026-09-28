# Traceability Matrix — RANDesk MVP

Versi 1.0 · 26 September 2026 · Pemetaan requirement → task → bukti uji.

Sumber ID: FR/NFR dari `RANDesk_Planning/PRD.md` §6/§8; BR dari `RANDesk_Planning/RULES.md` §4; TC dari `RANDesk_Planning/TEST_PLAN.md` §3; DEV dari `RANDesk_Planning/ROADMAP.md` §3; TASK dari [TASKS.md](TASKS.md). Kolom "TC pembuktian" mengikuti TEST_PLAN §4 dan ditambah TC lain yang relevan dengan task terkait. Tidak ada ID yang dikarang; seluruh ID di bawah diverifikasi ada di dokumen sumber.

## 1. Functional requirements (P0)

| FR | Ringkas | Task implementasi | BR terkait | TC pembuktian | Catatan |
| --- | --- | --- | --- | --- | --- |
| FR-01 | Autentikasi sesi | TASK-008, TASK-009, TASK-010, TASK-012, TASK-042 (rate limit login) | BR-17 | TC-01, TC-02, TC-03, TC-28 | Tercakup penuh |
| FR-02 | Kelola akun | TASK-011, TASK-014, TASK-015, TASK-018 | BR-15, BR-17 | TC-02, TC-24, TC-25, TC-28, TC-29 | Tercakup penuh |
| FR-03 | Kelola master data | TASK-016, TASK-019 | BR-14 | TC-26, TC-29 | Tercakup penuh |
| FR-04 | Buat dan baca tiket | TASK-022 (baca/detail), TASK-023, TASK-029 | BR-01, BR-02 | TC-04, TC-05, TC-06, TC-07 | Tercakup penuh |
| FR-05 | Daftar dan pencarian | TASK-022, TASK-028 | BR-16, BR-10 | TC-04, TC-15, TC-27 | Tercakup penuh |
| FR-06 | Assignment dan triage | TASK-020, TASK-025, TASK-030 | BR-03, BR-04, BR-05 | TC-08, TC-09, TC-10, TC-24 | Tercakup penuh |
| FR-07 | Lifecycle tiket | TASK-026, TASK-030 | BR-05, BR-07 | TC-11, TC-12, TC-13, TC-14, TC-35 | Tercakup penuh |
| FR-08 | Percakapan | TASK-031, TASK-036 | BR-09, BR-10 | TC-14, TC-15, TC-19 | Tercakup penuh |
| FR-09 | Lampiran | TASK-033, TASK-034, TASK-035, TASK-037 | BR-11 | TC-04, TC-16, TC-17, TC-18, TC-30 | Tercakup penuh |
| FR-10 | Riwayat aktivitas | TASK-021, TASK-032; event ditulis juga oleh TASK-023…026, 031, 034, 035 | BR-10, BR-12 | TC-07, TC-13, TC-15 | Tercakup penuh |
| FR-11 | Notifikasi | TASK-021 (tulis), TASK-038 (API), TASK-039 (UI polling) | BR-10, BR-12 | TC-20, TC-21, TC-22 | Tercakup penuh |
| FR-12 | Dashboard | TASK-040, TASK-041 | BR-13, BR-16 | TC-23 | Tercakup penuh |
| FR-13 | Audit administratif | TASK-013, TASK-019; writer dipakai TASK-011, 014, 016 | BR-18 | TC-29, TC-34 | Tercakup penuh |

## 2. Non-functional requirements

| NFR | Ringkas | Task implementasi | TC/bukti | Status ketercakupan |
| --- | --- | --- | --- | --- |
| NFR-01 | Otorisasi semua endpoint | TASK-009, TASK-020, TASK-022, TASK-044 | TC-03, TC-04, TC-05, TC-10, TC-15, TC-21 | Tercakup |
| NFR-02 | Konsistensi transaksi (mutasi+event+notifikasi) | TASK-021, TASK-023, TASK-024, TASK-025, TASK-026, TASK-031, TASK-034, TASK-045 | TC-07, TC-08, TC-09, TC-16, TC-18, TC-29 | Tercakup |
| NFR-03 | Performa p95 read ≤500 ms, write ≤800 ms | TASK-048 (+ TASK-017 dataset benchmark) | TC-32 | **Tercakup bersyarat** — lingkungan acuan 2 vCPU/4 GiB belum terverifikasi (GAP-04); TASK-048 `needs_verification` |
| NFR-04 | UX & aksesibilitas WCAG 2.2 AA | TASK-005 (tokens), TASK-046; praktik a11y di semua task UI (012, 018, 019, 028–030, 036, 037, 039, 041) | TC-31 | Tercakup (target, dibuktikan audit — bukan klaim) |
| NFR-05 | Ketahanan (timeout, batas body, rollback) | TASK-004, TASK-033, TASK-034, TASK-035, TASK-042, TASK-045 | TC-18, TC-19, TC-30 | Tercakup |
| NFR-06 | Operasi: backup/restore RPO 24 jam, RTO 4 jam | TASK-050, TASK-051 | TC-33 | Tercakup |
| NFR-07 | Maintainability (modul, lock, lint/build) | TASK-001, TASK-002, TASK-005, TASK-007 (gate CI); disiplin seluruh task | Build/lint/typecheck + hasil suite tiap milestone (TEST_PLAN §4) | Tercakup |
| NFR-08 | Observability (request ID, log, health) | TASK-002, TASK-004, TASK-043 | TC-34 + smoke health (TEST_PLAN §4) | Tercakup |

## 3. Business rules (RULES.md §4)

| BR | Ringkas | Task penegak utama | Task pendukung/verifikasi |
| --- | --- | --- | --- |
| BR-01 | Reporter dari sesi, satu assignee | TASK-023 | TASK-025, TASK-044 (TC-05 forged reporter) |
| BR-02 | ticket_no sequence, format HD- | TASK-023 | TASK-003 (identity column) |
| BR-03 | Assignment tidak mengubah status | TASK-025 | TASK-026 (CHECK tickets_assignee_state_ck) |
| BR-04 | Claim open kosong; reassign admin + reason | TASK-025 | TASK-020 (policy), TASK-030 (UI reason) |
| BR-05 | expected_version, row lock, transaksi | TASK-024, TASK-025, TASK-026 | TASK-045 (konkurensi) |
| BR-06 | 409 stale; 422 NO_CHANGE | TASK-024 | TASK-030 (conflict UX), TC-09 |
| BR-07 | Hanya transisi tabel lifecycle | TASK-026 | TASK-020, TC-11 |
| BR-08 | Prioritas default normal; hanya admin + reason | TASK-023, TASK-024 | TASK-029 (UI), TASK-030 |
| BR-09 | Komentar plain text append-only, visibility | TASK-031 | TASK-027 (grants), TASK-036 |
| BR-10 | Internal tidak bocor ke requester | TASK-021, TASK-022, TASK-031, TASK-032, TASK-038 | TASK-044 (sweep TC-15 menyeluruh) |
| BR-11 | 5 lampiran, 5 MiB, JPEG/PNG/PDF | TASK-033, TASK-034 | TASK-045 (TC-16), TASK-037 (batas klien) |
| BR-12 | Event+notifikasi dalam transaksi sama | TASK-021 | TASK-023…026, 031, 034, 035; TASK-045 (rollback) |
| BR-13 | UTC simpan, WITA tampil, durasi kalender | TASK-005 (lib tanggal), TASK-040 (metrik kalender) | TASK-028/030/036/041 (render WITA), TC-23 |
| BR-14 | Master dinonaktifkan bukan dihapus | TASK-016 | TASK-019, TASK-023 (MASTER_INACTIVE), TC-26 |
| BR-15 | Role immutable; jaga admin aktif; reassign sebelum nonaktif | TASK-014, TASK-015 | TASK-025 (TC-24), TASK-045 |
| BR-16 | Scope sama untuk pagination/search/aggregate/download | TASK-020 | TASK-022, TASK-035, TASK-038, TASK-040, TASK-044 |
| BR-17 | Sesi hilang saat nonaktif/reset/ganti password | TASK-008, TASK-009 | TASK-014, TASK-015, TC-02 |
| BR-18 | Tidak ada hapus tiket/komentar/event/audit via API | TASK-027 (grants) | TASK-016, TASK-031, TASK-052 (purge hanya via operator) |

## 4. Test cases (TEST_PLAN.md §3) → task yang mengimplementasikan verifikasi

| TC | Prioritas | Task pemilik verifikasi | Keterangan |
| --- | --- | --- | --- |
| TC-01 | Kritis | TASK-010 | Tes API contract auth |
| TC-02 | Kritis | TASK-009, TASK-010, TASK-012 | Sisi backend (revocation) + frontend (cache dibersihkan) |
| TC-03 | Kritis | TASK-009, TASK-010 | CSRF/Origin pada mutasi dan login |
| TC-04 | Kritis | TASK-022, TASK-035, TASK-044 | Detail/list/download lintas requester → 404 |
| TC-05 | Kritis | TASK-023, TASK-044 | Forged role/reporter/status |
| TC-06 | Tinggi | TASK-023 (+ validasi di 014, 016, 024, 026, 031) | Batas RULES §7, Unicode, password tidak ditrim (TASK-008) |
| TC-07 | Kritis | TASK-023, TASK-021, TASK-045 | Failure injection transaksi create |
| TC-08 | Kritis | TASK-025, TASK-045 | Claim bersamaan, satu pemenang |
| TC-09 | Kritis | TASK-024, TASK-030 | Stale version 409 + draft UI dipertahankan |
| TC-10 | Kritis | TASK-020, TASK-025, TASK-044 | Mutasi oleh teknisi non-assignee; izin ikut reassignment |
| TC-11 | Kritis | TASK-026 | Matriks penuh status × actor |
| TC-12 | Kritis | TASK-026 | Resolve tanpa/short summary → 422 |
| TC-13 | Kritis | TASK-026 | Reopen: state reset, reopen_count, event solusi lama |
| TC-14 | Kritis | TASK-026, TASK-031, TASK-034, TASK-036 | Closed read-only di semua jalur mutasi |
| TC-15 | Kritis | TASK-022, TASK-031, TASK-032, TASK-036, TASK-038, TASK-044 | Penanda internal nol kemunculan di semua jalur requester |
| TC-16 | Kritis | TASK-034, TASK-045 | Cap 5 lampiran saat konkurensi |
| TC-17 | Kritis | TASK-033 | MIME palsu, extension, traversal, kosong, oversize |
| TC-18 | Kritis | TASK-034, TASK-035, TASK-045 | Kompensasi commit gagal & byte soft-deleted; cleanup |
| TC-19 | Tinggi | TASK-031, TASK-045 | Komentar/upload vs close bersamaan |
| TC-20 | Tinggi | TASK-021 | Pemilihan penerima, dedup, nonaktif, aktor |
| TC-21 | Kritis | TASK-038 | Mark-read lintas pengguna → 404 |
| TC-22 | Tinggi | TASK-039 | Polling lifecycle (hidden/offline/resume) |
| TC-23 | Tinggi | TASK-040, TASK-041 | Metrik dashboard vs fixture, cohort, timezone |
| TC-24 | Kritis | TASK-015, TASK-025, TASK-045 | Deaktivasi vs assignment bersamaan |
| TC-25 | Kritis | TASK-015, TASK-045 | Dual-admin deactivation, admin aktif tersisa |
| TC-26 | Tinggi | TASK-016, TASK-023 | Master nonaktif: referensi baru ditolak, historis terbaca |
| TC-27 | Kritis | TASK-022, TASK-044 | Injeksi SQL/sort terlarang/filter staff dari requester |
| TC-28 | Tinggi | TASK-010, TASK-014, TASK-018 | must_change_password + temporary password sekali tampil |
| TC-29 | Tinggi | TASK-013, TASK-014 | Audit gagal → rollback; audit admin-only |
| TC-30 | Tinggi | TASK-037, TASK-034 | Create sukses + upload gagal: tiket tetap, retry per file |
| TC-31 | Tinggi | TASK-046 | Keyboard, focus dialog, zoom 200%, mobile |
| TC-32 | Tinggi | TASK-048 | **needs_verification** — lingkungan acuan NFR-03 (GAP-04) |
| TC-33 | Kritis | TASK-051 | Backup set + restore drill lingkungan kosong |
| TC-34 | Kritis | TASK-042, TASK-043 | Log redaction + header/cookie produksi |
| TC-35 | Tinggi | TASK-047 (+ TASK-030 manual) | E2E lintas peran + refresh browser |

## 5. DEV backlog (ROADMAP §3) → task

| DEV | Task terkait |
| --- | --- |
| DEV-01 | TASK-001, TASK-002, TASK-005, TASK-006 |
| DEV-02 | TASK-003, TASK-017, TASK-027 |
| DEV-03 | TASK-002, TASK-004 |
| DEV-04 | TASK-008, TASK-009, TASK-010, TASK-011, TASK-012 |
| DEV-05 | TASK-013, TASK-014, TASK-015, TASK-016, TASK-018, TASK-019 |
| DEV-06 | TASK-020, TASK-022 |
| DEV-07 | TASK-021, TASK-023, TASK-024, TASK-025, TASK-026 |
| DEV-08 | TASK-028, TASK-029, TASK-030 |
| DEV-09 | TASK-031, TASK-032, TASK-036 |
| DEV-10 | TASK-033, TASK-034, TASK-035, TASK-037 |
| DEV-11 | TASK-038, TASK-039 |
| DEV-12 | TASK-040, TASK-041 |
| DEV-13 | TASK-007, TASK-042, TASK-043, TASK-044, TASK-045, TASK-046, TASK-047, TASK-048, TASK-049 |
| DEV-14 | TASK-050, TASK-051, TASK-052, TASK-053, TASK-054 |

## 6. Requirement/target yang belum atau bersyarat tercakup

| Item | Status | Alasan & tindak lanjut |
| --- | --- | --- |
| NFR-03 / TC-32 (target performa terukur) | **Bersyarat** | Memerlukan lingkungan acuan 2 vCPU/4 GiB/SSD yang belum terbukti tersedia (GAP-04). TASK-048 `needs_verification`; hasil pada lingkungan berbeda wajib dicatat sebagai deviasi, bukan pemenuhan target. |
| PRD §4.2 target UAT (4/5 peserta, median ≤2 menit) | **Diputuskan — di luar cakupan MVP** | D-06 diputuskan 26 Sep 2026: walkthrough terstruktur sesuai TEST_PLAN §7 (asumsi P-07) diadopsi resmi sebagai bukti penerimaan. TASK-049 `todo`; target PRD §4.2 baris 1–2 dicatat tidak terukur (belum mewakili pengguna nyata) di laporan. |
| Release gate TEST_PLAN §8 (keseluruhan) | Direncanakan | Dibuktikan pada akhir M6 lewat TASK-054 (checklist dengan tautan bukti); belum ada item yang dapat dicentang saat ini. |
| NFR-04 klaim "WCAG 2.2 AA terpenuhi" | Direncanakan sebagai target | Hanya boleh diklaim setelah audit TASK-046; DESIGN §4.1/§8 melarang klaim tanpa pemeriksaan. |
| Operasional: monitoring/alert berkelanjutan (OPERATIONS §8) | **Di luar cakupan task MVP** | Sinyal ambang didefinisikan OPERATIONS §8 tetapi tidak ada TC/DEV khusus MVP; TASK-050 hanya smoke test readiness. Dicatat di GAPS_AND_DECISIONS.md GAP-08 sebagai keputusan yang dapat ditunda (demo/pilot terbatas). |

Seluruh FR P0 (13/13), NFR (8/8, satu bersyarat lingkungan), BR (18/18), dan TC (35/35; TC-32 bersyarat lingkungan) memiliki task yang bertanggung jawab. Tidak ada requirement P0 tanpa task; tidak ada task berstatus `blocked` per 26 September 2026 (D-06 diselesaikan lewat adopsi resmi P-07).
