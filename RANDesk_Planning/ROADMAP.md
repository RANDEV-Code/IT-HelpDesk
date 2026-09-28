# Implementation Roadmap — RANDesk

Versi 1.0 · 26 September 2026 · Rencana bagi pengembang individual

## 1. Strategi pengerjaan

Bangun satu alur vertikal yang dapat digunakan dari browser hingga database, lalu perluas sesuai prioritas. Pelajari Go dan React melalui modul yang sedang dikerjakan. Hindari menyelesaikan seluruh backend tanpa pernah mencoba pengalaman pengguna.

Estimasi awal pekerjaan inti 82–122 jam. Tambahkan ruang belajar, debugging, dan revisi sekitar 20–30%, sehingga rencana praktis sekitar 100–160 jam. Pada alokasi 15 jam per minggu, kisaran awal sekitar 7–11 minggu. Ini estimasi perencanaan, bukan komitmen tenggat; evaluasi setelah milestone kedua.

## 2. Milestone dan exit criteria

| Milestone | Estimasi inti | Pekerjaan | Bukti selesai |
| --- | --- | --- | --- |
| M0 — Fondasi | 10–14 jam | Repo, versi dependency, struktur modul, PostgreSQL dev, migrasi, React shell, error envelope | DB kosong dapat dimigrasi; UI memanggil health/API; lint/build dasar berjalan |
| M1 — Akun dan sesi | 12–18 jam | Bootstrap admin, login/logout, password, CSRF, akun, kategori/departemen | Login lintas peran dan perubahan password bekerja; TC-01/02/03/28 lulus |
| M2 — Tiket inti | 16–24 jam | Create/list/detail, scope, assignment, lifecycle, event, version | Alur tiket sampai closed; TC-04/05/07/08/09/10/11/12/13 lulus |
| M3 — Kolaborasi dan berkas | 12–18 jam | Komentar publik/internal, upload/download/delete, compensation | Catatan internal terisolasi; batas file dan rollback diuji |
| M4 — Notifikasi dan dashboard | 8–12 jam | Persistence, polling, read-state, agregasi, filter tanggal | Penerima benar; metrik cocok fixture; TC-20/21/22/23 lulus |
| M5 — Hardening dan UX | 16–24 jam | Negative tests, concurrency, responsive, a11y, load test, UAT | Temuan kritis ditutup; hasil uji dicatat; scope MVP stabil |
| M6 — Deploy dan portofolio | 8–12 jam | TLS, backup/restore, CI release, README, OpenAPI, seed, demo video | Deployment dapat dipulihkan; orang lain dapat mencoba setup dan demo |

M0 → M1 → M2 merupakan jalur utama. Infrastruktur notifikasi tabel/service disiapkan saat M2 karena perubahan tiket dan notifikasi harus satu transaksi; UI polling dan dashboard diselesaikan pada M4. Pengujian mengikuti setiap milestone, bukan menunggu M5.

## 3. Backlog implementasi awal

| ID | Prioritas | Deliverable | Dependensi |
| --- | --- | --- | --- |
| DEV-01 | P0 | Inisialisasi backend/frontend dan konfigurasi lingkungan | Keputusan versi dependency |
| DEV-02 | P0 | Migrasi baseline, seed sintetis, runtime DB role | DEV-01, SCHEMA |
| DEV-03 | P0 | Error mapping, logging, request ID, health | DEV-01 |
| DEV-04 | P0 | Session auth, CSRF, password, account policy | DEV-02, DEV-03 |
| DEV-05 | P0 | Master data dan admin accounts/audit | DEV-04 |
| DEV-06 | P0 | Policy tiket dan query scope | DEV-04, DEV-05 |
| DEV-07 | P0 | Ticket service, event/notification writer, concurrency | DEV-06 |
| DEV-08 | P0 | UI list/create/detail dan lifecycle | DEV-07 |
| DEV-09 | P0 | Comments dan internal-note filtering | DEV-07, DEV-08 |
| DEV-10 | P0 | Storage adapter dan attachment workflow | DEV-07, DEV-08 |
| DEV-11 | P0 | Notification API/UI dan polling lifecycle | DEV-07 |
| DEV-12 | P0 | Dashboard query dan presentasi metrik | DEV-07, DEV-08 |
| DEV-13 | P0 | Integration/E2E/operational verification | Bertahap sejak DEV-04 |
| DEV-14 | P0 | Deployment, dokumentasi setup, OpenAPI, demo | Seluruh gate MVP |

OpenAPI dihasilkan/ditulis pada implementasi dan divalidasi terhadap API_SPEC. Endpoint yang belum tersedia harus ditandai sebagai backlog, bukan dipresentasikan seolah sudah berjalan.

## 4. Definition of Ready dan Done

Sebuah pekerjaan siap dimulai ketika requirement, role yang boleh bertindak, input/output, validasi, serta skenario gagal sudah jelas. Jika tidak, perbarui dokumen pemilik keputusan terlebih dahulu.

Pekerjaan selesai jika acceptance criteria terbukti, perubahan data/izin konsisten, UI menangani state yang relevan, tes bermakna lulus, log tidak membocorkan secret, dan dokumentasi yang terpengaruh diperbarui. Deployment bukan tanda selesai jika restore atau akses dasarnya belum diverifikasi.

## 5. Fitur lanjutan

| P1/P2 | Manfaat | Prasyarat desain |
| --- | --- | --- |
| SSE | Status/notifikasi tampil lebih cepat | Reconnect, refetch, stream authorization, proxy timeout |
| SLA dan pengingat | Batas waktu penanganan terpantau | Definisi response/resolution SLA, business hours, timezone, pause, worker/outbox |
| Email | Pengguna menerima pemberitahuan di luar aplikasi | Provider, retry, idempotency, preference, delivery status |
| Ekspor laporan | Hasil dapat dianalisis di luar aplikasi | Izin, filter yang konsisten, sanitasi spreadsheet formula |
| Inventaris aset | Tiket terhubung ke perangkat | Lifecycle aset dan izin lokasi/pemilik |
| Multi-organisasi | Satu layanan melayani beberapa organisasi | Tenant isolation end-to-end, unique keys, migrasi dan audit keamanan baru |

Prioritaskan satu peningkatan setelah MVP berdasarkan masalah pengguna atau nilai pembelajaran yang terukur.

## 6. Risiko dan mitigasi

| Risiko | Sinyal awal | Respons |
| --- | --- | --- |
| Belajar Go dan React sekaligus | Milestone fondasi melebar tanpa alur selesai | Gunakan satu form/tabel sederhana, selesaikan vertical slice |
| Scope bertambah | Fitur SLA/chat/AI masuk sebelum tiket inti | Masukkan backlog P1/P2 dan review sesudah MVP |
| Otorisasi tersebar | Hasil detail dan list berbeda | Satukan policy, tambah tes negatif lintas endpoint |
| Upload memakan waktu | File tersisa atau metadata tak konsisten | Implementasi storage adapter dan failure-injection lebih awal |
| UI dipoles sebelum alur utuh | Banyak komponen tetapi status belum berfungsi | Selesaikan journey tiga peran sebelum polish lanjutan |
| Hosting berbeda dari lokal | Cookie/proxy/storage rusak di staging | Deploy staging sejak M2 dengan data sintetis |

## 7. Paket presentasi portofolio

Siapkan README aplikasi yang menjelaskan masalah, pengguna, keputusan stack, fitur yang benar-benar selesai, cara menjalankan, akun demo yang aman, batasan, dan bukti pengujian. Tambahkan diagram arsitektur, ERD, dokumentasi API, screenshot responsif, serta video alur sekitar 3–5 menit.

Cerita demo: pegawai melaporkan Wi-Fi bermasalah; teknisi mengambil tiket; admin melihat distribusi pekerjaan; teknisi menulis catatan internal dan balasan publik; solusi dikonfirmasi pelapor. Tunjukkan satu kasus konflik claim atau akses ditolak untuk menjelaskan kedalaman backend.

Klaim CV harus mengikuti hasil implementasi. Contoh setelah fitur terbukti: “Membangun aplikasi IT Helpdesk menggunakan Go, React, dan PostgreSQL dengan kontrol akses, workflow tiket, serta pengujian konkurensi assignment.” Jangan menambahkan angka performa atau jumlah pengguna tanpa data.
