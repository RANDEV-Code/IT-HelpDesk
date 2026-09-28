# RANDesk — Paket Perencanaan Produk dan Engineering

| Atribut | Nilai |
| --- | --- |
| Produk | RANDesk — IT Helpdesk untuk satu organisasi |
| Pemilik proyek | Rico Adrian Naibaho / RANDEV |
| Versi dokumen | 1.0 |
| Tanggal | 26 September 2026 |
| Bahasa produk | Bahasa Indonesia; identifier teknis berbahasa Inggris |
| Status | Baseline rancangan untuk implementasi; belum divalidasi pengguna |
| Stack | Go + Gin, PostgreSQL, React + TypeScript + Vite |

## 1. Tujuan paket

Paket ini menjadi acuan membangun aplikasi pelaporan dan penanganan masalah IT yang dapat didemonstrasikan dalam portofolio. Dokumen menghubungkan kebutuhan produk, pengalaman pengguna, aturan bisnis, kontrak API, penyimpanan data, dan cara memverifikasi implementasinya.

Dokumen disusun dengan sudut pandang product planner, business analyst, UX planner, dan software architect. Semua target, estimasi, persona, dan skenario penggunaan adalah asumsi perencanaan, bukan hasil penelitian lapangan atau hasil pengujian aplikasi. RANDesk adalah nama kerja yang dapat diganti.

## 2. Keputusan dasar

- Satu organisasi, satu instalasi, satu backend Go dengan modul yang terpisah secara jelas.
- Tiga peran: `requester`, `technician`, dan `admin`; masing-masing akun memiliki satu peran.
- Empat status: `open`, `in_progress`, `resolved`, dan `closed`.
- Autentikasi menggunakan sesi server dan cookie HttpOnly, dengan perlindungan CSRF.
- Semua akses objek diperiksa backend; frontend menampilkan tindakan sesuai izin.
- Notifikasi dalam aplikasi tersedia pada MVP melalui polling 30 detik ketika tab aktif.
- Pembaruan melalui SSE, SLA, email, dan pekerjaan latar belakang masuk fase lanjutan.
- Semua tanggal disimpan dalam UTC; antarmuka menggunakan `Asia/Makassar` dan label WITA.

## 3. Isi dokumen dan urutan membaca

| Urutan | File | Fungsi |
| --- | --- | --- |
| 1 | [PRD.md](PRD.md) | Masalah, pengguna, cakupan, requirement, acceptance criteria, target kualitas |
| 2 | [RULES.md](RULES.md) | Aturan bisnis, matriks izin, status, validasi, dan aturan pengembangan |
| 3 | [DESIGN.md](DESIGN.md) | Struktur halaman, alur pengguna, spesifikasi layar, komponen, aksesibilitas |
| 4 | [ARCHITECTURE.md](ARCHITECTURE.md) | Batas sistem, modul, transaksi, autentikasi, deployment, perkembangan arsitektur |
| 5 | [SCHEMA.md](SCHEMA.md) | Kamus data, ERD, DDL PostgreSQL, indeks, integritas, dan migrasi |
| 6 | [API_SPEC.md](API_SPEC.md) | Endpoint, payload, otorisasi, error, pagination, dan konflik perubahan |
| 7 | [SECURITY.md](SECURITY.md) | Kontrol keamanan, perlindungan data, dan penanganan penyalahgunaan |
| 8 | [TEST_PLAN.md](TEST_PLAN.md) | Skenario pengujian, keterlacakan requirement, UAT, dan release gate |
| 9 | [ROADMAP.md](ROADMAP.md) | Urutan pengerjaan, dependensi, milestone, risiko, dan paket demo |
| 10 | [OPERATIONS.md](OPERATIONS.md) | Konfigurasi, rilis, backup, restore, log, dan respons gangguan |
| 11 | [DECISIONS.md](DECISIONS.md) | Catatan alasan keputusan dan asumsi yang perlu ditinjau |
| 12 | [REFERENCES.md](REFERENCES.md) | Rujukan teknis primer dan batas penggunaan sumber |

## 4. Prioritas implementasi

**P0 / MVP:** autentikasi, akun dan master data, pembuatan tiket, assignment, lifecycle, komentar, lampiran, riwayat, notifikasi polling, dashboard dasar, serta pengujian akses dan transaksi.

**P1:** SSE untuk notifikasi, pengingat SLA dengan worker dan outbox, email, pencarian yang lebih lanjut, serta peningkatan operasional berdasarkan kebutuhan nyata.

**P2:** SSO, multi-organisasi, inventaris aset, integrasi chat, dan aplikasi mobile. Fitur ini memerlukan desain tambahan sebelum dikembangkan.

## 5. Konvensi dan sumber kebenaran

| Jenis keputusan | Dokumen pemilik |
| --- | --- |
| Tujuan, scope, dan acceptance criteria | PRD |
| Izin dan aturan perubahan data | RULES |
| Struktur penyimpanan | SCHEMA |
| Kontrak HTTP dan bentuk data | API_SPEC |
| Perilaku layar | DESIGN |
| Cara komponen sistem berinteraksi | ARCHITECTURE |
| Kontrol keamanan | SECURITY |

Jika ditemukan konflik, hentikan implementasi bagian yang terpengaruh, catat keputusan di DECISIONS, lalu perbarui seluruh dokumen terkait dalam perubahan yang sama. Jangan memilih aturan yang paling mudah secara diam-diam.

Requirement menggunakan `FR-xx`, kualitas menggunakan `NFR-xx`, aturan bisnis menggunakan `BR-xx`, keputusan menggunakan `ADR-xxx`, dan kasus uji menggunakan `TC-xx`. Istilah “wajib” berarti persyaratan baseline; “target” berarti hasil yang masih harus diukur.

## 6. Penggunaan saat mulai coding

1. Baca PRD dan RULES, lalu buat backlog sesuai ROADMAP.
2. Catat versi dependensi stabil yang benar-benar dipakai dalam lockfile dan go.mod.
3. Implementasikan satu alur utuh: login → membuat tiket → assignment → penanganan → penyelesaian → penutupan.
4. Tambahkan pengujian yang membuktikan aturan akses dan integritas pada alur tersebut.
5. Perluas fitur sesuai urutan milestone; perbarui dokumentasi ketika kontrak berubah.

Paket ini berisi spesifikasi Markdown. DDL dan contoh JSON di dalamnya adalah rancangan implementasi, bukan aplikasi yang sudah berjalan. Pengujian runtime, migrasi pada PostgreSQL, aksesibilitas, dan performa tetap menjadi pekerjaan implementasi.
