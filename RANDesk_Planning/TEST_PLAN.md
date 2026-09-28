# Quality Assurance and Acceptance Plan — RANDesk

Versi 1.0 · 26 September 2026 · Seluruh skenario di bawah adalah rencana uji

## 1. Tujuan dan strategi

Verifikasi terfokus pada aturan yang paling berisiko: akses lintas pengguna, status dan assignment, konkurensi, isolasi catatan internal, pencabutan sesi, serta konsistensi database/file. Tes tidak boleh hanya mengulang bentuk implementasi tanpa membuktikan perilaku pengguna.

| Level | Area | Pendekatan |
| --- | --- | --- |
| Unit | Policy, validasi, transisi, pemilihan penerima, metrik | Kasus tabel dengan batas input dan state yang relevan |
| Integration | Query, transaksi, constraints, sesi, file compensation | PostgreSQL nyata yang dapat di-reset; jangan mengganti dengan SQLite untuk uji lock |
| API contract | Status HTTP, DTO, error, scope, pagination | Request melalui HTTP test server dan database test |
| Frontend component | Form, error, conflict, internal composer | Uji interaksi dan accessible labels |
| E2E | Alur tiga peran | Browser otomatis dengan akun/data sintetis |
| Manual | Responsif, keyboard, screen reader, demo | Checklist dan bukti screenshot/catatan |
| Operasional | Backup, restore, rollback | Lingkungan staging terisolasi |

## 2. Fixture minimum

Gunakan dua requester (R1/R2), dua technician (T1/T2), dan dua admin (A1/A2). Sediakan akun nonaktif, satu kategori nonaktif, beberapa departemen, serta tiket open kosong, open assigned, in_progress, resolved, closed, dan reopened.

Setiap requester memiliki minimal dua tiket. Salah satu tiket mempunyai komentar internal dengan teks penanda unik untuk mendeteksi kebocoran di seluruh endpoint. Lampiran menggunakan file sintetis valid dan file uji tidak valid. Waktu bisnis memakai clock yang dapat dikontrol dalam service.

## 3. Kasus uji utama

Prioritas Kritis berarti kegagalan memblokir rilis MVP. Tinggi berarti wajib selesai sebelum acceptance akhir; Sedang dapat dicatat sebagai defect UX dengan alasan dan rencana perbaikan.

| ID | Prioritas | Skenario | Hasil yang diwajibkan |
| --- | --- | --- | --- |
| TC-01 | Kritis | Login benar/salah/nonaktif | Hanya akun aktif berkredensial benar memperoleh sesi; error gagal generik |
| TC-02 | Kritis | Expiry, logout, change-password, reset, deactivation | Sesi terdampak tidak dapat digunakan lagi; cache frontend dibersihkan |
| TC-03 | Kritis | CSRF hilang/salah dan Origin asing pada mutasi; login form lintas origin | Ditolak tanpa perubahan data; login same-origin JSON tetap berfungsi |
| TC-04 | Kritis | R1 meminta tiket/file/comment/event R2 | 404; data dan keberadaan objek tidak bocor |
| TC-05 | Kritis | Requester mengirim role/reporter/status buatan | Field ditolak; actor dan reporter tetap berasal dari sesi |
| TC-06 | Tinggi | Validasi field kosong, batas minimum/maksimum, Unicode | Batas RULES berlaku konsisten; password tidak ditrim |
| TC-07 | Kritis | R1 membuat tiket dengan transaksi event gagal | Tidak ada tiket atau notifikasi setengah jadi; sequence boleh melompat |
| TC-08 | Kritis | T1 dan T2 claim bersamaan dengan version sama | Satu sukses, satu 409, satu assignee dan satu event claim berhasil |
| TC-09 | Kritis | Edit metadata memakai version stale | 409; perubahan terbaru tetap utuh; draft UI dipertahankan |
| TC-10 | Kritis | T1 mencoba mengubah tiket milik assignment T2 | 403; setelah reassignment, izin mutasi mengikuti assignee terbaru |
| TC-11 | Kritis | Seluruh pasangan status dan actor | Hanya empat transisi RULES lolos; invalid state 409 |
| TC-12 | Kritis | Resolve tanpa ringkasan atau ringkasan terlalu pendek | 422, status tetap in_progress |
| TC-13 | Kritis | Reopen resolved oleh reporter/admin | Open, assignee null, field resolusi null, reopen_count +1, solusi lama ada pada event |
| TC-14 | Kritis | Mutasi metadata/comment/file pada closed | Ditolak, timeline dan byte lampiran tidak berubah |
| TC-15 | Kritis | Catatan internal pada semua jalur requester | Penanda internal tidak muncul di list/detail/comment/event/search/count/notifikasi |
| TC-16 | Kritis | Dua upload bersamaan ketika ada 4 lampiran aktif | Hanya satu menjadi lampiran kelima; lainnya 422; file sisa dibersihkan |
| TC-17 | Kritis | MIME palsu, extension terlarang, path traversal, kosong, oversize | Ditolak sesuai error contract; tidak ada file executable/public |
| TC-18 | Kritis | Gagal commit setelah file dipindahkan; gagal hapus byte setelah soft delete | Kompensasi/maintenance membersihkan orphan; deleted attachment tetap tidak dapat diunduh |
| TC-19 | Tinggi | Komentar dan upload pada tiket yang ditutup bersamaan | Lock menentukan urutan sah; operasi yang melihat closed ditolak |
| TC-20 | Tinggi | Notifikasi aktor, duplikat recipient, akun nonaktif, internal note | Penerima sesuai RULES, tanpa duplikat per event dan tanpa informasi internal ke requester |
| TC-21 | Kritis | R1 mark-read notifikasi R2 | 404; notifikasi R2 tidak berubah |
| TC-22 | Tinggi | Polling, offline, tab tersembunyi, resume | Tidak polling pada hidden/offline; refetch saat aktif; tidak ada request loop |
| TC-23 | Tinggi | Dashboard dengan cohort, NULL metric, reopen, timezone | Perhitungan sesuai PRD; nilai kosong null; scope requester konsisten |
| TC-24 | Kritis | Deactivate technician saat assignment bersamaan | Tidak pernah menghasilkan tiket aktif yang ditugaskan ke teknisi nonaktif |
| TC-25 | Kritis | Dua admin mencoba menonaktifkan satu sama lain | Sistem tetap memiliki admin aktif; aturan self-deactivation diterapkan |
| TC-26 | Tinggi | Pemilihan master nonaktif dan deactivation setelah tiket ada | Referensi baru nonaktif ditolak; data historis tetap terbaca |
| TC-27 | Kritis | Cari/count/sort dengan input SQL atau filter staff dari requester | Query aman; kolom sort terlarang ditolak; scope tidak melebar |
| TC-28 | Tinggi | Akun baru must_change_password dan temporary password | Endpoint lain 403 sampai password diganti; secret hanya sekali tampil dan tidak dilog |
| TC-29 | Tinggi | Akun/master diubah dan insert audit gagal | Perubahan ikut rollback; audit hanya dapat dibaca admin |
| TC-30 | Tinggi | Form create berhasil tetapi upload gagal | Satu tiket tetap ada; retry hanya file gagal; pesan UX jelas |
| TC-31 | Tinggi | Akses keyboard, focus dialog, 200% zoom, mobile | Alur utama dapat dijalankan; tidak ada kontrol/teks utama terpotong |
| TC-32 | Tinggi | Load test target NFR-03 | Catat p95, error rate, CPU/RAM/DB; bandingkan dengan target, tanpa mengarang hasil |
| TC-33 | Kritis | Backup DB+file dan restore ke lingkungan kosong | Tiket, akun, event, lampiran, checksum, serta unduh terotorisasi dapat diverifikasi |
| TC-34 | Kritis | Request body/log redaction dan header produksi | Tidak ada password/token/body sensitif dalam log; HTTPS/cookie/headers benar |
| TC-35 | Tinggi | Alur E2E lintas peran dan refresh browser | Login → create → claim → start → comment → resolve → close berhasil tanpa kehilangan state |

## 4. Requirement traceability

| Requirement | Bukti uji utama |
| --- | --- |
| FR-01 | TC-01, TC-02, TC-03, TC-28 |
| FR-02 | TC-02, TC-24, TC-25, TC-28, TC-29 |
| FR-03 | TC-26, TC-29 |
| FR-04 | TC-04, TC-05, TC-06, TC-07 |
| FR-05 | TC-04, TC-15, TC-27 |
| FR-06 | TC-08, TC-09, TC-10, TC-24 |
| FR-07 | TC-11, TC-12, TC-13, TC-14, TC-35 |
| FR-08 | TC-14, TC-15, TC-19 |
| FR-09 | TC-04, TC-16, TC-17, TC-18, TC-30 |
| FR-10 | TC-07, TC-13, TC-15 |
| FR-11 | TC-20, TC-21, TC-22 |
| FR-12 | TC-23 |
| FR-13 | TC-29, TC-34 |
| NFR-01 | TC-03, TC-04, TC-05, TC-10, TC-15, TC-21 |
| NFR-02 | TC-07, TC-08, TC-09, TC-16, TC-18, TC-29 |
| NFR-03 | TC-32 |
| NFR-04 | TC-31 |
| NFR-05 | TC-18, TC-19, TC-30 |
| NFR-06 | TC-33 |
| NFR-07 | Build/lint/typecheck dan hasil suite tiap milestone |
| NFR-08 | TC-34 dan pemeriksaan health/log pada smoke test |

## 5. Pengujian konkurensi

Gunakan dua koneksi DB atau HTTP clients independen dan barrier untuk melepaskan request berdekatan. Verifikasi respons dan state akhir database, bukan hanya jumlah respons 200. Kasus claim, attachment cap, dan deactivation harus dijalankan berulang secukupnya untuk menemukan interleaving yang salah.

Tambahkan failure injection di boundary event insert, notification insert, file move, dan commit. Pastikan rollback serta cleanup terlihat pada bukti uji. Jangan memakai mock untuk mengklaim bahwa PostgreSQL row locking telah diuji.

## 6. CI dan perintah target

Perintah ini adalah kontrak yang harus disiapkan saat repo dibuat; belum dapat dijalankan dari paket dokumen ini.

| Area | Gate target |
| --- | --- |
| Backend | gofmt check, `go vet ./...`, `go test ./...`, build binary |
| Race-sensitive tests | `go test -race ./...` pada runner yang mendukung race detector |
| Frontend | `npm ci`, lint, typecheck, unit/component test, production build |
| Integrasi | Migrasi ke PostgreSQL kosong, fixture, suite transaksi/otorisasi |
| E2E | Alur TC-35 dengan browser otomatis, artefak screenshot saat gagal |
| Dependency | Pemeriksaan kerentanan dependensi Go/npm dan kajian temuan |

Tidak ada target coverage persen yang menggantikan bukti aturan kritis. Simpan hasil nyata, versi environment, dan tanggal eksekusi.

## 7. UAT

Target awal lima peserta dengan perwakilan pelapor dan staff. Bila belum tersedia, lakukan walkthrough terstruktur dan catat bahwa hasilnya belum mewakili pengguna nyata.

Tugas peserta: membuat tiket; menemukan progres; menambahkan informasi; mengambil dan menangani tiket untuk peserta staff; mengonfirmasi solusi atau membuka kembali tiket. Catat berhasil/gagal, waktu, bantuan yang diperlukan, kesalahan, dan komentar. Jangan mengarahkan peserta ke tombol sebelum mengamati kesulitannya.

## 8. Release gate

- [ ] Semua FR P0 memiliki bukti uji atau catatan penerimaan yang dapat ditelusuri.
- [ ] Seluruh test Kritis lulus; tidak ada celah akses atau kehilangan data yang diketahui.
- [ ] Kasus berprioritas Tinggi selesai atau memiliki keputusan risiko tertulis sebelum acceptance.
- [ ] Lint, typecheck, build, integration, dan E2E utama lulus.
- [ ] Cookie, CSRF, batas upload, log redaction, dan data demo telah diverifikasi.
- [ ] Migrasi pada DB kosong dan restore backup telah dicoba.
- [ ] README aplikasi, OpenAPI, petunjuk setup, seed, dan screenshot sesuai implementasi.
- [ ] Hasil performa dan aksesibilitas dicatat sebagai hasil pengukuran nyata.

## 9. Template laporan hasil

Untuk setiap run, catat tanggal, commit, environment, dataset, test IDs, expected result, actual result, status, bukti, defect, owner perbaikan, dan hasil retest. Dokumen ini tidak mengisi kolom hasil sebelum pengujian dilaksanakan.
