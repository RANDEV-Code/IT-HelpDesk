# Product Requirements Document — RANDesk

Versi 1.0 · 26 September 2026 · Status: baseline perencanaan

## 1. Ringkasan produk

RANDesk membantu pegawai melaporkan masalah IT dan membantu tim IT mengelola penanganannya sampai dikonfirmasi selesai. Produk menyediakan satu tempat untuk melihat pemilik pekerjaan, progres, bukti pendukung, percakapan, dan riwayat keputusan.

Nilai portofolio yang dituju adalah kemampuan merancang dan membangun aplikasi bisnis secara menyeluruh: frontend berbasis komponen, API Go, database relasional, otorisasi per objek, perubahan data atomik, dan operasi aplikasi.

## 2. Masalah dan hipotesis

| Masalah yang diasumsikan | Dampak | Respons produk |
| --- | --- | --- |
| Laporan tersebar di chat atau disampaikan lisan | Laporan terlewat dan informasi kurang lengkap | Formulir tiket dengan kategori, prioritas, dan lampiran |
| Penanggung jawab tidak terlihat | Pekerjaan ganda atau tidak ditangani | Assignment dengan riwayat perubahan |
| Pelapor harus berulang kali meminta kabar | Interupsi kerja dan ketidakpastian | Status, percakapan publik, dan notifikasi |
| Penyelesaian tidak terdokumentasi | Solusi sulit ditelusuri | Ringkasan solusi dan timeline |
| Beban layanan tidak terukur | Sulit mengevaluasi proses | Dashboard volume dan durasi penanganan |

Hipotesis ini perlu dikonfirmasi melalui wawancara singkat dengan calon pengguna sebelum digunakan pada organisasi nyata.

## 3. Pengguna dan pekerjaan utama

| Persona asumsi | Peran sistem | Kebutuhan utama |
| --- | --- | --- |
| Pegawai nonteknis | requester | Membuat laporan dengan mudah dan mengetahui tindak lanjut |
| Staf dukungan IT | technician | Melihat antrean, mengambil tugas yang kosong, mencatat penanganan |
| Koordinator IT | admin | Memilah prioritas, membagi pekerjaan, mengelola akun, melihat hasil layanan |

Semua peran boleh membuat tiket untuk dirinya sendiri. Akun tidak dapat membuat tiket atas nama orang lain pada MVP. Hak staff tetap mengikuti peran akun.

## 4. Tujuan dan batas keberhasilan

### 4.1 Tujuan produk

- Pelapor dapat membuat tiket tanpa bantuan teknisi.
- Setiap tiket yang sedang dikerjakan memiliki satu teknisi penanggung jawab aktif.
- Perubahan assignment, prioritas, dan status dapat ditelusuri.
- Pelapor hanya memperoleh informasi yang memang boleh dilihatnya.
- Pengembang dapat mendemonstrasikan alur lintas peran dengan data sintetis.

### 4.2 Target evaluasi awal

| Indikator | Target awal | Cara ukur |
| --- | --- | --- |
| Keberhasilan membuat tiket | Minimal 4 dari 5 peserta menyelesaikan tugas tanpa bantuan | UAT setelah MVP, jika peserta tersedia |
| Waktu membuat tiket sederhana | Median ≤ 2 menit, tanpa waktu menyiapkan lampiran | Pengamatan UAT; bukan hasil yang telah dicapai |
| Integritas alur kritis | Seluruh kasus prioritas kritis pada TEST_PLAN lulus | Pengujian integrasi dan E2E |
| Kebocoran akses pada fixture pengujian | 0 kasus akses lintas pelapor atau catatan internal | Pengujian negatif API |
| Demo portofolio | Alur tiga peran dapat dijalankan ulang dari seed | Demonstrasi manual dan panduan setup |

## 5. Scope

### 5.1 P0 — MVP

Login, logout, perubahan password, manajemen akun dasar, kategori dan departemen, tiket, assignment, lifecycle, komentar publik dan internal, lampiran, timeline, notifikasi dalam aplikasi, dashboard, pencarian sederhana, serta dokumentasi dan pengujian.

### 5.2 P1 — setelah MVP stabil

SSE, pengingat SLA, email, mekanisme outbox dan worker, ekspor laporan, filter tersimpan, serta perbaikan berdasarkan UAT. Fitur P1 tidak termasuk syarat selesai MVP.

### 5.3 Di luar baseline

Multi-tenancy, billing, inventaris perangkat, kontrol jarak jauh, integrasi WhatsApp, AI diagnosis, native mobile, SSO, rich text, lampiran pada catatan internal, registrasi publik, dan reset password mandiri melalui email.

## 6. Functional requirements

| ID | Requirement P0 | Acceptance criteria utama |
| --- | --- | --- |
| FR-01 | Autentikasi sesi | Login akun aktif menghasilkan sesi; salah kredensial mendapat pesan generik; logout dan password change mencabut sesi sesuai SECURITY |
| FR-02 | Kelola akun | Admin membuat akun, melihat daftar, mereset password, mengubah nama/departemen, serta mengaktifkan/nonaktifkan akun sesuai BR-15 |
| FR-03 | Kelola master data | Admin membuat dan mengubah kategori/departemen; item nonaktif tidak dapat dipilih pada input baru tetapi riwayat tetap terbaca |
| FR-04 | Buat dan baca tiket | Judul 5–150 karakter, deskripsi 20–10.000 karakter, kategori aktif, prioritas valid; reporter berasal dari sesi; nomor tiket unik |
| FR-05 | Daftar dan pencarian | Scope akses diterapkan sebelum pencarian/count/pagination; filter status, prioritas, kategori, assignment dan teks tersedia sesuai peran |
| FR-06 | Assignment dan triage | Admin dapat menetapkan teknisi; teknisi mengambil tiket open yang belum ditugaskan; dua claim bersamaan hanya menghasilkan satu pemenang |
| FR-07 | Lifecycle tiket | Hanya transisi dalam RULES yang diterima; resolve memerlukan ringkasan; reopen mengembalikan tiket ke open tanpa assignee; closed terminal |
| FR-08 | Percakapan | Komentar publik tersedia untuk pengguna yang berhak; catatan internal hanya staff; komentar bersifat append-only; closed tidak menerima komentar |
| FR-09 | Lampiran | Maksimal 5 lampiran aktif per tiket, 5 MiB per file, JPEG/PNG/PDF; unduh memeriksa akses; kegagalan unggah tidak menghapus tiket |
| FR-10 | Riwayat aktivitas | Perubahan inti menyimpan aktor, waktu, tipe, dan informasi yang diizinkan; requester tidak menerima event internal |
| FR-11 | Notifikasi | Notifikasi persisten untuk penerima yang ditentukan; polling 30 detik saat tab aktif; tandai dibaca hanya milik sendiri |
| FR-12 | Dashboard | Jumlah status, backlog aktif, distribusi kategori, rata-rata respons pertama dan penyelesaian sesuai scope dan definisi metrik |
| FR-13 | Audit administratif | Perubahan akun, master data, dan reset password tercatat tanpa secret; hanya admin dapat membaca audit |

## 7. User stories prioritas

- Sebagai pelapor, saya ingin mengirim masalah beserta bukti agar teknisi memahami kondisinya.
- Sebagai pelapor, saya ingin melihat tiket saya dan balasan publik tanpa mengetahui percakapan internal tim IT.
- Sebagai teknisi, saya ingin mengambil tiket yang kosong agar tidak bekerja pada tugas yang sedang diambil orang lain.
- Sebagai admin, saya ingin mengubah prioritas dengan alasan agar keputusan penanganan dapat dipahami.
- Sebagai teknisi penanggung jawab, saya ingin mencatat solusi agar pelapor dapat memverifikasinya.
- Sebagai pelapor, saya ingin membuka kembali tiket resolved yang belum tuntas agar penanganan dilanjutkan.

### Contoh acceptance scenario FR-06

**Given:** tiket open, belum memiliki assignee, version 1. **When:** dua teknisi berbeda mengirim claim dengan expected_version 1. **Then:** satu request berhasil; lainnya mendapat 409; hanya satu assignment berlaku; event dan notifikasi konsisten dengan pemenang.

### Contoh acceptance scenario FR-08

**Given:** teknisi menulis komentar internal pada tiket milik seorang requester. **When:** requester mengambil detail, komentar, timeline, notifikasi, dan hasil pencarian. **Then:** teks internal dan event internal tidak muncul pada seluruh respons tersebut.

## 8. Non-functional requirements

Semua angka berikut adalah target verifikasi, bukan klaim performa aplikasi.

| ID | Area | Target dan kondisi |
| --- | --- | --- |
| NFR-01 | Otorisasi | Semua endpoint terlindungi mengecek sesi aktif, peran, dan cakupan objek; kasus negatif kritis lulus |
| NFR-02 | Konsistensi | Mutasi inti, event, dan notifikasi berada dalam satu transaksi database; tidak ada event sukses dari transaksi gagal |
| NFR-03 | Performa | Pada lingkungan acuan 2 vCPU/4 GiB/SSD, 10.000 tiket dan 50.000 komentar, 25 pengguna virtual: p95 read API ≤ 500 ms dan write JSON ≤ 800 ms selama 5 menit; login/upload dikecualikan |
| NFR-04 | UX dan aksesibilitas | Jalur utama dapat digunakan dengan keyboard; label/error terbaca; target WCAG 2.2 AA untuk layar MVP, dibuktikan audit manual dan otomatis |
| NFR-05 | Ketahanan | Timeout, batas body, dan rollback terpasang; kegagalan DB/upload menghasilkan error yang dapat ditindaklanjuti tanpa data setengah jadi |
| NFR-06 | Operasi | Backup dan restore diuji; target awal RPO 24 jam dan RTO 4 jam untuk lingkungan demo/private pilot |
| NFR-07 | Maintainability | Modul jelas, dependensi dikunci, lint/typecheck/build lulus, pengujian aturan kritis tersedia |
| NFR-08 | Observability | Request ID, log terstruktur, durasi/status request dan health endpoint tersedia; log tidak berisi password atau token |

Performa diukur menggunakan campuran 80% read dan 20% write, API dan database pada host acuan, generator beban di host terpisah. Gunakan pacing satu request per pengguna virtual per detik, input valid, dan total error tak terduga <1%; pengujian rate limit serta konflik dilakukan terpisah. Catat versi software, dataset, error rate, latency, CPU, dan memori. Jika target gagal, profilkan penyebab sebelum mengubah target atau infrastruktur.

## 9. Definisi metrik dashboard

Filter tanggal menggunakan rentang `created_at` tiket: `from` inklusif dan `to` eksklusif dalam UTC. UI mengonversi hari WITA menjadi rentang UTC. Tanpa filter, gunakan seluruh data yang boleh dilihat pengguna.

| Metrik | Definisi |
| --- | --- |
| Jumlah per status | Banyak tiket pada status saat ini dalam cohort filter |
| Backlog aktif | Jumlah tiket `open` + `in_progress` dalam cohort |
| Menunggu konfirmasi | Jumlah tiket `resolved` dalam cohort |
| Respons pertama | Rata-rata `first_response_at - created_at` untuk tiket yang memiliki first_response_at |
| Waktu penyelesaian | Rata-rata `resolved_at - created_at` untuk status resolved/closed yang memiliki resolved_at |
| Distribusi kategori | Jumlah tiket per kategori dalam cohort, termasuk kategori yang sudah nonaktif |

Respons pertama berarti aksi publik pertama staff selain reporter: komentar publik atau transisi ke in_progress/resolved. Assignment dan catatan internal tidak dihitung. Reopen tidak mereset respons pertama; resolved_at mengikuti penyelesaian terakhir. Nilai tanpa sampel ditampilkan “Belum ada data”, bukan nol. Ini adalah durasi kalender, bukan jam kerja atau metrik SLA.

## 10. Dependensi dan penerimaan rilis

MVP memerlukan server Go, PostgreSQL, penyimpanan lampiran privat, frontend, dan reverse proxy TLS. Layanan email dan broker tidak menjadi dependensi MVP.

Rilis diterima jika requirement P0 memiliki bukti uji, seluruh kasus kritis lulus, tidak ada cacat yang membuka akses data tanpa izin, setup dapat direproduksi, backup dapat dipulihkan, dan demo menggunakan data sintetis. Checklist lengkap berada di [TEST_PLAN.md](TEST_PLAN.md).
