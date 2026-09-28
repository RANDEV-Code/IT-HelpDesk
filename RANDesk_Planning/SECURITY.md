# Security and Data Protection — RANDesk

Versi 1.0 · 26 September 2026 · Baseline kontrol teknis MVP

## 1. Model perlindungan

Aset yang dilindungi: akun, kredensial, isi tiket, catatan internal, lampiran, riwayat tindakan, dan backup. Batas kepercayaan berada antara browser/API, API/database, dan API/storage. Browser tidak dipercaya untuk menetapkan role, reporter, visibility yang boleh digunakan, atau hasil validasi.

| Risiko utama | Kontrol produk |
| --- | --- |
| Requester membaca tiket pengguna lain | Scope query dan policy per objek pada seluruh endpoint |
| Catatan internal bocor melalui search/count/notifikasi | Filtering sebelum pagination serta DTO allowlist |
| Sesi dicuri atau tidak dicabut | Cookie HttpOnly/Secure, hash token, expiry, revocation |
| Aksi pengguna dipicu situs lain | Origin validation dan token CSRF terikat sesi |
| File berbahaya atau akses file langsung | Jenis/ukuran terbatas, storage privat, otorisasi download |
| Dua aktor membuat perubahan bertentangan | Row lock, expected_version, constraints, transaksi |
| Password/isi tiket masuk log | Log terstruktur tanpa request body atau header sensitif |

## 2. Password dan akun

Gunakan Argon2id dengan library terawat. Baseline parameter proyek: memory 64 MiB, iterations 3, parallelism 1, salt acak minimal 16 byte, output 32 byte, format PHC. Benchmark pada host target dan catat perubahan parameter; rujukan password storage berada di REFERENCES.

Password 12–128 karakter, menerima spasi, tidak ditrim, tidak dipotong, dan tidak dibatasi pola komposisi yang arbitrer. Batas panjang diperiksa sebelum hashing. Gunakan pembandingan hash dari library.

Akun baru/reset memperoleh password sementara acak minimal 16 karakter dan must_change_password. Hanya admin yang dapat membuat akun; registrasi publik tidak tersedia. UI tidak menyimpan password sementara setelah ditampilkan. Jangan merekam password dalam event, audit, trace, atau error.

## 3. Sesi dan pencabutan

Session token terdiri dari 32 byte acak kriptografis yang di-encode base64url. Database hanya menyimpan SHA-256 token; token berentropi tinggi tidak diperlakukan seperti password manusia. Masa aktif absolut delapan jam, tidak diperpanjang oleh aktivitas.

| Lingkungan | Cookie |
| --- | --- |
| Produksi HTTPS | `__Host-randesk_session`; HttpOnly; Secure; SameSite=Lax; Path=/; tanpa Domain |
| Development localhost HTTP | `randesk_session_dev`; HttpOnly; SameSite=Lax; Path=/; Secure=false hanya pada mode development |

Sesi baru dibuat setelah login; token lama pada browser yang sama dicabut. Logout mencabut sesi saat ini. Password change/reset dan deactivation mencabut semua sesi user dalam transaksi yang sama dengan perubahan akun.

Setiap request mengecek expiry, revoked_at, dan is_active user. Jika must_change_password aktif, hanya endpoint me/logout/change-password yang lolos. Request yang sudah berjalan sebelum pencabutan bisa menyelesaikan transaksi sah; request berikutnya wajib ditolak.

Jangan log Cookie, Set-Cookie, Authorization, password, atau csrf_token. Frontend membersihkan cache dan draft sensitif ketika sesi berakhir.

## 4. CSRF dan origin

Server membuat CSRF token acak 32 byte terpisah per sesi; simpan dalam sessions.csrf_token dan kirim melalui JSON login/me. Frontend menyimpannya di memori, mengirim `X-CSRF-Token` untuk POST/PUT/PATCH/DELETE, dan mengambil ulang lewat me setelah reload. Bandingkan token secara constant-time.

API memeriksa Origin terhadap APP_ORIGIN pada mutasi browser. Jika Origin tidak ada, validasi origin dari Referer; jika keduanya tidak ada atau tidak sesuai, tolak. SameSite menambah lapisan perlindungan dan bukan pengganti pemeriksaan token.

Login tidak memiliki CSRF token sesi, tetapi wajib memakai JSON dan Origin/Referer yang sah, menolak form content type, serta tidak mengaktifkan credentialed cross-origin access. Endpoint GET tidak boleh mengubah state bisnis.

Development memakai proxy Vite agar browser tetap satu origin. Gunakan allowlist eksplisit jika arsitektur origin berubah, disertai ADR dan pengujian baru.

## 5. Otorisasi dan input

- Akses objek yang tidak berada dalam scope baca menghasilkan 404; aksi terlarang atas objek yang boleh dilihat menghasilkan 403.
- Batas scope juga berlaku pada count, aggregate, daftar file, timeline, dan notifikasi.
- Frontend tidak menentukan actor_id, requester_id, role efektif, timestamps, atau version baru.
- DTO input menggunakan allowlist; unknown fields ditolak.
- SQL selalu memakai parameter untuk nilai; sort dan nama kolom berasal dari allowlist internal.
- Teks user dirender sebagai text node. Tidak ada raw HTML atau rich text pada MVP.
- Runtime database user berbeda dari pemilik migrasi; komentar/event/audit tidak memperoleh hak update/delete.

## 6. Lampiran

Jenis yang diperbolehkan: .jpg/.jpeg dengan image/jpeg, .png dengan image/png, dan .pdf dengan application/pdf. Cocokkan extension yang dinormalisasi, MIME yang dideteksi server, dan signature; jangan percaya header browser saja. Tolak file kosong, file >5 MiB, nama/path berbahaya, dan request lebih dari satu file.

Simpan dengan key acak yang dibuat server di luar webroot. Nama asli hanya metadata tampilan setelah sanitasi. Download selalu melewati otorisasi tiket, menggunakan Content-Disposition attachment dan X-Content-Type-Options nosniff. File yang sudah soft-delete tidak dapat diunduh.

Pemeriksaan tipe bukan jaminan file bebas malware. MVP private demo memakai berkas sintetis/tepercaya; sebelum menerima unggahan tak tepercaya secara publik, tambahkan pemindaian malware atau nonaktifkan upload publik. Rancangan upload mengikuti rujukan OWASP pada REFERENCES.

## 7. Batas trafik dan resource

Nilai awal untuk diuji pada deployment satu instance:

| Jalur | Batas awal | Respons |
| --- | --- | --- |
| Login per IP | 20 percobaan/15 menit | 429 dan Retry-After |
| Login per email ternormalisasi + IP | 5 kegagalan/15 menit | 429 generik; tidak mengungkap akun |
| API umum per user | 120 request/menit dengan burst 30 | 429 |
| Upload per user | 10 file/menit | 429 |
| JSON body | 64 KiB | 413 |
| Multipart body | 6 MiB | 413 |

Limiter memory pada satu instance boleh dipakai untuk MVP dan akan reset saat restart; catat keterbatasannya. Jangan menerapkan lockout permanen yang memudahkan penyerang memblokir akun. Jika aplikasi memakai beberapa replika, limiter harus dikoordinasikan sebelum rilis.

Hanya percaya proxy forwarding header dari proxy yang dikonfigurasi. Hindari membaca seluruh file besar ke memory; gunakan streaming dengan batas byte. Batasi pool dan timeout database.

## 8. Header dan secret

Produksi menggunakan HTTPS. Terapkan nosniff, Referrer-Policy yang membatasi kebocoran URL, serta Content-Security-Policy sesuai aset self-hosted: default-src 'self', object-src 'none', frame-ancestors 'none', base-uri 'self', dan connect-src 'self'. Tambahkan pengecualian hanya jika kebutuhan fitur tercatat. Aktifkan HSTS setelah domain/HTTPS stabil.

Database URL, password, dan credential backup berada pada environment atau secret store, bukan repository. File `.env.example` pada implementasi hanya berisi nama dan placeholder. Gunakan secret berbeda untuk development, staging, dan production; rotasi jika terekspos.

## 9. Data demo dan incident response

Gunakan data sintetis, nama contoh, dan domain reserved seperti example.test. Jangan menaruh data pribadi organisasi dalam screenshot, README, atau demo publik. Demo publik tidak boleh memberikan akun admin dengan akses ke pengaturan server atau data riil; isolasikan dataset dan pertimbangkan menonaktifkan upload/admin mutation untuk pengunjung.

Jika terjadi insiden: batasi akses fitur/akun terdampak, simpan log yang relevan, cabut sesi/credential terkait, pulihkan keadaan aman, verifikasi akses dan integritas, lalu catat akar masalah serta perbaikannya. Prosedur operasional berada pada OPERATIONS.

## 10. Gate keamanan MVP

Wajib lulus pengujian lintas requester, internal note isolation, CSRF, revocation, forged role/reporter, file path/size/type, SQL input, stale version, serta log redaction. Pindai dependency dengan alat yang tersedia dan kaji temuan kritis sebelum deployment. Kontrol dalam dokumen ini belum berarti aplikasi telah diaudit atau tersertifikasi.
