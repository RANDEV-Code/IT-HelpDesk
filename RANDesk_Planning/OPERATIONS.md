# Deployment and Operations — RANDesk

Versi 1.0 · 26 September 2026 · Runbook rencana untuk satu instance

## 1. Lingkungan

| Environment | Tujuan | Data |
| --- | --- | --- |
| Development | Coding dan eksperimen | Seed sintetis, dapat di-reset |
| Test/CI | Pengujian otomatis | Database terisolasi per run |
| Staging | UAT, migrasi, restore, pemeriksaan proxy | Salinan data sintetis |
| Production/private demo | Demonstrasi atau pilot terbatas | Dataset yang telah ditinjau; tanpa secret organisasi pada demo publik |

Pisahkan DB, volume, akun, dan credential setiap lingkungan. Backend dapat dijalankan lokal atau container; deployment target satu VM Linux dengan persistent volume dan reverse proxy TLS.

## 2. Konfigurasi minimum

| Variabel rencana | Fungsi | Ketentuan |
| --- | --- | --- |
| APP_ENV | development/test/staging/production | Production gagal start jika konfigurasi tidak aman |
| APP_ORIGIN | Origin browser yang sah | Scheme + host + port, tanpa wildcard |
| HTTP_ADDR | Bind address API | Privat di belakang proxy pada deployment |
| DATABASE_URL | Koneksi PostgreSQL | Secret; tidak dilog |
| STORAGE_ROOT | Direktori lampiran | Path privat, writable oleh API, persistent |
| SESSION_TTL_HOURS | Masa aktif sesi | Default 8; perubahan memerlukan review UX/security |
| LOG_LEVEL | Verbosity | Info di production; tanpa body sensitif |
| TRUSTED_PROXIES | Proxy yang dipercaya | Daftar eksplisit, bukan seluruh internet |
| DB_MAX_OPEN_CONNS | Batas koneksi | Baseline 10; ukur pool wait sebelum mengubah |
| DB_MAX_IDLE_CONNS | Koneksi idle | Baseline 5 |

Batas upload, panjang field, dan enum bisnis merupakan aturan produk. Jika dijadikan konfigurasi, backend dan frontend harus menerima satu sumber konfigurasi yang selaras. Jangan memiliki angka berbeda pada proxy, UI, dan API.

## 3. Bootstrap

1. Siapkan DB dan credential migrasi/runtime terpisah.
2. Jalankan migrasi yang telah diuji; terapkan grants runtime.
3. Buat admin pertama melalui perintah operator interaktif atau secret environment sementara. Hash password memakai kode aplikasi dan paksa perubahan password awal.
4. Buat kategori/departemen awal dan akun sintetis bila environment demo.
5. Jalankan API, verifikasi readiness, lalu sajikan frontend melalui proxy.
6. Periksa login, cookie, CSRF, satu tiket, dan download lampiran.

Jangan menaruh password bootstrap pada argumen CLI yang tersimpan di shell history atau process list. Hapus secret bootstrap setelah digunakan.

## 4. Rilis dan rollback

Setiap rilis memiliki commit/tag, versi dependency terkunci, hasil CI, migrasi, catatan perubahan, dan rencana rollback. Build artifact yang sama dipromosikan dari staging ke deployment target.

Urutan rilis: backup terverifikasi → migrasi kompatibel → deploy API/frontend → readiness → smoke test tiga peran → pantau error. Gunakan migrasi expand/contract untuk perubahan yang tidak dapat di-rollback sederhana.

Rollback aplikasi memakai artifact sebelumnya jika schema masih kompatibel. Jangan otomatis menjalankan down migration yang menghapus data. Jika perubahan schema destruktif telah diterapkan, gunakan rencana pemulihan yang sudah diuji dan putuskan dampak data secara eksplisit.

## 5. Backup yang konsisten

Target awal demo/private pilot: RPO 24 jam, RTO 4 jam. Target ini harus diuji, bukan diasumsikan.

Pada arsitektur satu instance, prosedur paling sederhana adalah maintenance window: hentikan penerimaan request API, tunggu request aktif selesai, hentikan proses API serta maintenance writer, lalu ambil logical dump PostgreSQL dan salinan volume lampiran pada keadaan yang tidak berubah. Setelah keduanya selesai, hidupkan layanan dan verifikasi.

Simpan DB dump, file, manifest waktu/versi aplikasi/migrasi, dan checksum bersama sebagai satu backup set. Enkripsi dan simpan di lokasi berbeda dari host aplikasi. Credential backup memiliki akses terbatas. Backup DB tanpa volume lampiran tidak dianggap backup lengkap.

Baseline retensi backup: 7 backup harian dan 4 backup mingguan untuk demo. Jika nanti memakai backup online, desain mekanisme snapshot/koordinasi DB+storage terlebih dahulu.

## 6. Restore drill

1. Pilih backup set dan lingkungan kosong yang terisolasi.
2. Cocokkan versi PostgreSQL dan aplikasi yang mendukung schema backup.
3. Pulihkan DB dan lampiran dengan owner/permission yang benar.
4. Cabut sesi hasil restore sebelum lingkungan dibuka; jangan menghidupkan sesi lama tanpa review.
5. Verifikasi jumlah objek, FK, tiket lintas status, timeline, dan checksum sampel file.
6. Jalankan smoke test login, akses lintas peran, serta download terotorisasi.
7. Catat waktu pemulihan, selisih data terhadap RPO, kegagalan, dan perbaikannya.

Lakukan restore drill sebelum acceptance MVP dan ulangi sesudah perubahan storage atau mekanisme backup.

## 7. Maintenance dan retensi

| Data | Baseline demo/private pilot | Pelaksanaan |
| --- | --- | --- |
| Sesi expired/revoked | Hapus setelah 7 hari dari expiry/revocation | Perintah maintenance terjadwal, role terbatas |
| Log aplikasi | 14 hari | Rotasi dan pembatasan akses |
| Notifikasi yang sudah dibaca | 90 hari | Cleanup batch; tidak menghapus event tiket |
| Audit administratif | 180 hari | Review sebelum purge melalui operator |
| Tiket/komentar/event/lampiran aktif | Dipertahankan selama pilot | Penghapusan memerlukan kebijakan pemilik data |
| File staging/orphan | Cleanup setelah 24 jam dan verifikasi tidak direferensikan | Maintenance dengan grace period |
| Byte lampiran soft-delete | Hapus sesudah commit; retry jika gagal | Metadata/timeline tetap disimpan |

Angka ini adalah asumsi produk demo, bukan nasihat kepatuhan hukum atau kebijakan organisasi. Tinjau ulang sebelum menerima data riil. Pertimbangkan dampak backup saat memenuhi permintaan penghapusan.

Gunakan maintenance role terpisah untuk purge; runtime tidak diberi kemampuan menghapus audit/event/comment. Cleanup harus memiliki dry-run, log jumlah item, dan batas batch. Jangan menghapus file hanya karena tidak ditemukan dalam query yang gagal.

## 8. Monitoring awal

| Sinyal | Ambang investigasi awal |
| --- | --- |
| API 5xx | >1% selama 5 menit, minimum 100 request sebagai konteks |
| Readiness | Gagal pada beberapa probe berturut-turut |
| Disk | Terpakai >80% atau pertumbuhan tak wajar |
| DB pool | Waktu tunggu meningkat atau koneksi mendekati batas terus-menerus |
| Upload cleanup | Ada file orphan yang bertahan setelah maintenance |
| Backup | Backup set terakhir lebih tua dari 24 jam atau restore drill gagal |

Endpoint health tidak menampilkan connection string, detail filesystem, atau stack trace. Log menggunakan route template agar UUID/isi input tidak menjadi dimensi tak terbatas. Simpan request_id untuk penelusuran.

## 9. Respons gangguan

- DB tidak tersedia: readiness gagal, tampilkan pesan gangguan sementara, hentikan retry write otomatis dari frontend.
- Storage penuh: tolak upload dengan error yang jelas, pertahankan tiket, periksa kapasitas dan orphan; jangan hapus file aktif sembarangan.
- Lonjakan login gagal: periksa rate limit dan proxy IP, cabut sesi jika ada indikasi kompromi, hindari mempublikasikan identitas akun.
- Bug otorisasi: batasi fitur terdampak, simpan bukti, perbaiki policy dan regression test, tinjau kemungkinan akses yang sudah terjadi.
- Deploy bermasalah: kembalikan artifact kompatibel atau jalankan rencana pemulihan; catat dampak dan hasil verifikasi.

## 10. Batas kesiapan operasional

Runbook ini belum membuktikan deployment, backup, alert, atau recovery telah dijalankan. Implementasi harus melampirkan bukti nyata dan nama penanggung jawab sebelum aplikasi dipakai sebagai layanan operasional.
