# Architecture Decisions and Assumptions — RANDesk

Versi 1.0 · 26 September 2026

Keputusan berstatus “baseline” adalah pilihan perencanaan untuk mulai implementasi. Status tersebut tidak menyatakan adanya persetujuan organisasi, pengujian produksi, atau hasil benchmarking.

## 1. ADR register

| ID | Keputusan baseline | Alasan | Konsekuensi dan pemicu evaluasi ulang |
| --- | --- | --- | --- |
| ADR-001 | Go + Gin untuk REST API | Sesuai stack yang dipilih pemilik proyek; Go memiliki dokumentasi resmi tutorial Gin | Pengembang harus mempelajari error handling, context, dan SQL transaksi secara eksplisit |
| ADR-002 | React + TypeScript + Vite | Memperluas kemampuan frontend, cocok untuk dashboard dan komponen interaktif | Routing/server state dipilih eksplisit; SSR tidak menjadi kebutuhan MVP |
| ADR-003 | Modular monolith | Membatasi deployment dan memungkinkan debugging alur utuh | Pemisahan service baru dipertimbangkan jika ada kebutuhan scaling atau kepemilikan tim yang nyata |
| ADR-004 | PostgreSQL dan SQL eksplisit | Relasi, constraints, transaksi, dan row lock mendukung workflow | Schema/migration harus disiplin; query kompleks perlu benchmark |
| ADR-005 | Sesi server dengan cookie dan CSRF | Browser satu origin; mudah mencabut akses akun dan sesi | Memerlukan penyimpanan sesi dan pemeriksaan DB; klien mobile/SSO akan membutuhkan keputusan tambahan |
| ADR-006 | Empat status dan closed terminal | Menjaga lifecycle yang jelas dan dapat diuji | Waiting/paused/auto-close memerlukan aturan serta metrik baru |
| ADR-007 | Polling notifikasi pada MVP | Durabilitas berasal dari DB dengan implementasi sederhana | Latensi notifikasi normal sampai sekitar 30 detik ditambah waktu jaringan; SSE masuk P1 |
| ADR-008 | Lampiran di volume privat | Cukup untuk satu instance dan mudah dipelajari | Replika API memerlukan shared object storage; backup harus mencakup DB dan file |
| ADR-009 | Audit dan event berbeda | Event untuk timeline produk, audit untuk administrasi | Perlu allowlist payload dan kontrol akses terpisah |
| ADR-010 | Version pada mutasi inti | Mencegah edit menimpa perubahan pengguna lain | UI harus menyediakan conflict recovery; komentar/attachment tetap lock tanpa bump version |
| ADR-011 | Akun/master dinonaktifkan, tiket tidak dihapus | Referensi historis dan pelacakan tetap konsisten | Perlu prosedur operator untuk retensi/penghapusan data di luar MVP |
| ADR-012 | Single organization dan role immutable | Menjaga scope realistis bagi pengembang individual | Multi-tenancy dan perubahan role memerlukan migrasi serta review otorisasi tersendiri |

## 2. Asumsi perencanaan

| ID | Asumsi | Dampak jika berubah |
| --- | --- | --- |
| A-01 | Proyek awal untuk portofolio dan private pilot | Jika dipakai untuk produksi organisasi, validasi kebutuhan, kebijakan data, dan operasi kembali |
| A-02 | Satu organisasi, sekitar 200 akun dan 25 pengguna aktif bersamaan sebagai target uji | Beban lebih tinggi memerlukan pengukuran baru |
| A-03 | Semua pengguna dapat memakai aplikasi web | Kebutuhan offline/native memerlukan kontrak sinkronisasi tambahan |
| A-04 | Admin menyediakan akun dan password sementara lewat kanal tepercaya | Jika onboarding massal diperlukan, desain invitation dan email verification |
| A-05 | Seluruh teknisi boleh membaca semua tiket dan catatan staff | Organisasi dengan pembatasan departemen membutuhkan policy lebih rinci |
| A-06 | Durasi dihitung waktu kalender; tanpa business hours pada MVP | SLA jam kerja memerlukan kalender, timezone, hari libur, dan aturan pause |
| A-07 | Tidak ada file rahasia organisasi nyata dalam demo publik | Penggunaan data riil memerlukan review retensi dan pengamanan operasional |

## 3. Keputusan yang perlu dikunci saat kickoff

| ID | Item | Default yang disarankan | Waktu penyelesaian |
| --- | --- | --- | --- |
| D-01 | Nama final produk dan domain | RANDesk sebagai nama kerja; domain dipilih pemilik | Sebelum publikasi |
| D-02 | Tool migrasi | golang-migrate | Sebelum migrasi pertama |
| D-03 | Reverse proxy dan hosting | Caddy pada satu VM Linux/container | Sebelum milestone deployment |
| D-04 | Versi dependency tepat | Rilis stabil yang kompatibel saat inisialisasi | Milestone fondasi |
| D-05 | Retensi data organisasi | Baseline demo pada OPERATIONS; validasi untuk pilot nyata | Sebelum data riil masuk |
| D-06 | Calon peserta UAT | Target lima calon pengguna, peran terwakili | Sebelum tahap UAT |

Keputusan terbuka tersebut tidak mengubah kontrak bisnis inti. Jangan menunda penulisan alur utama hanya karena nama/domain belum dipilih.

## 4. Template ADR baru

Gunakan struktur: ID dan judul; tanggal; status; konteks masalah; opsi yang ditimbang; keputusan; alasan; konsekuensi; migrasi/rollback; dokumen yang harus diperbarui; bukti yang akan dipakai untuk mengevaluasi keputusan.

## 5. Catatan batas verifikasi

Paket ini dirancang agar requirement, izin, schema, dan API saling sesuai. Tidak ada klaim bahwa codebase, server, database, deployment, benchmark, atau aplikasi produksi telah dibangun. Hasil implementasi harus dicatat dalam laporan pengujian tersendiri.
