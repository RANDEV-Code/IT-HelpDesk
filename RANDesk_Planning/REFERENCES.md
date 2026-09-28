# Technical References — RANDesk

Versi 1.0 · Tanggal peninjauan sumber: 26 September 2026

## 1. Cara menggunakan sumber

Sumber berikut merupakan dokumentasi primer untuk memverifikasi kemampuan teknologi dan prinsip implementasi. Requirement bisnis, jumlah peran, status tiket, target performa, batas file, estimasi waktu, dan keputusan scope adalah rancangan proyek ini, bukan rekomendasi wajib dari seluruh sumber.

Link `current` atau `latest` dapat berubah. Saat implementasi, catat versi yang dipakai dan gunakan dokumentasi yang sesuai versi tersebut. Tidak ada versi dependency yang diasumsikan sudah terpasang oleh paket ini.

## 2. Produk dan frontend

| Sumber | Penggunaan dalam rancangan |
| --- | --- |
| [React — Thinking in React](https://react.dev/learn/thinking-in-react) | Pemecahan UI menjadi komponen dan aliran data |
| [React — Build a React app from Scratch](https://react.dev/learn/build-a-react-app-from-scratch) | Jalur setup Vite dan TypeScript untuk frontend terpisah |
| [TanStack Query — Caching Examples](https://tanstack.com/query/latest/docs/framework/react/guides/caching) | Pemahaman cache dan refetch server state |
| [Tailwind CSS — Installation with Vite](https://tailwindcss.com/docs) | Integrasi styling pada toolchain Vite |

## 3. Backend dan database

| Sumber | Penggunaan dalam rancangan |
| --- | --- |
| [Go — Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin) | Dasar API HTTP dengan Go/Gin |
| [Go — Executing transactions](https://go.dev/doc/database/execute-transactions) | Batas transaksi dan penggunaan sql.Tx |
| [Go — Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency) | Dasar concurrency untuk pengembangan lanjutan |
| [pgx stdlib](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib) | Adapter pgx untuk database/sql |
| [PostgreSQL — Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html) | FK, uniqueness, CHECK, dan batas pemeriksaan lintas tabel |
| [PostgreSQL — Explicit Locking](https://www.postgresql.org/docs/17/explicit-locking.html) | Semantik row lock; cocokkan kembali dengan versi deployment |

## 4. Keamanan dan aksesibilitas

| Sumber | Penggunaan dalam rancangan |
| --- | --- |
| [OWASP — Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html) | Sesi, cookie, dan perlindungan identifier sesi |
| [OWASP — Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html) | Password hashing dengan algoritma adaptif |
| [OWASP — CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html) | Token CSRF, origin validation, dan perlindungan browser |
| [OWASP — File Upload Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html) | Validasi berkas, storage privat, dan download terotorisasi |
| [W3C — Understanding Contrast Minimum](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html) | Target kontras teks antarmuka |
| [W3C — Understanding Target Size Minimum](https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum) | Pertimbangan ukuran dan jarak target interaksi |

## 5. Batas bukti

Referensi teknis tidak membuktikan bahwa aplikasi sudah aman, cepat, mudah digunakan, atau siap produksi. Semua hasil tersebut memerlukan implementasi dan verifikasi yang dijelaskan pada TEST_PLAN dan OPERATIONS.
