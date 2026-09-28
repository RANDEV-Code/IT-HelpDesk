# UI/UX Design Specification — RANDesk

Versi 1.0 · 26 September 2026 · Spesifikasi desain untuk implementasi frontend

## 1. Arah desain

Karakter visual: profesional, tenang, jelas, dan berorientasi pekerjaan. Halaman menonjolkan status, penanggung jawab, serta tindakan berikutnya. Gunakan ruang kosong dan hierarki teks; hindari grafik dekoratif yang tidak membantu keputusan.

Dokumen ini mendefinisikan tata letak dan perilaku. Ini belum merupakan hasil riset pengguna, mockup Figma, atau antarmuka yang telah diimplementasikan.

## 2. Prinsip pengalaman pengguna

1. Satu tindakan utama terlihat pada setiap konteks, misalnya “Buat tiket”, “Mulai penanganan”, atau “Konfirmasi selesai”.
2. Label menggunakan bahasa pengguna. Identifier `resolved` ditampilkan “Selesai ditangani”; `closed` ditampilkan “Ditutup”.
3. Form mempertahankan input ketika request gagal. Pesan menjelaskan field atau tindakan yang perlu diperbaiki.
4. Tindakan sensitif menampilkan konteks dan akibat yang konkret.
5. Status penting tidak disampaikan dengan warna saja; selalu ada label.
6. Catatan internal dipisahkan secara visual dan diberi teks “Hanya terlihat oleh tim IT”.

## 3. Information architecture

| Route frontend | Halaman | Akses |
| --- | --- | --- |
| `/login` | Login | Pengunjung; pengguna bersesi diarahkan ke dashboard |
| `/change-password` | Ganti password | Semua akun bersesi, termasuk yang wajib mengganti password |
| `/dashboard` | Ringkasan layanan | Semua peran; scope mengikuti RULES |
| `/tickets` | Daftar tiket | Semua peran |
| `/tickets/new` | Buat tiket | Semua peran |
| `/tickets/:id` | Detail, percakapan, dan riwayat | Sesuai akses objek |
| `/notifications` | Semua notifikasi | Milik pengguna sendiri |
| `/admin/users` | Akun pengguna | Admin |
| `/admin/categories` | Kategori tiket | Admin |
| `/admin/departments` | Departemen | Admin |
| `/admin/audit` | Audit administratif | Admin |

Route guard frontend mengarahkan navigasi, sementara API tetap memeriksa izin. URL objek yang tidak boleh dibaca menggunakan halaman “Tiket tidak ditemukan” agar tidak mengonfirmasi keberadaan objek.

## 4. Design tokens

### 4.1 Warna

| Token | Nilai awal | Pemakaian |
| --- | --- | --- |
| `color.background` | `#F8FAFC` | Latar aplikasi |
| `color.surface` | `#FFFFFF` | Kartu, tabel, form |
| `color.text.primary` | `#0F172A` | Judul dan teks utama |
| `color.text.secondary` | `#475569` | Keterangan dan metadata |
| `color.primary` | `#1D4ED8` | Tombol utama, link, focus |
| `color.border` | `#CBD5E1` | Pemisah dekoratif dan tabel |
| `color.control.border` | `#64748B` | Batas input/checkbox yang perlu terlihat jelas |
| `color.danger` | `#B91C1C` | Error dan tindakan berisiko |
| `color.success` | `#166534` | Keberhasilan |

| Status | Label | Teks | Latar |
| --- | --- | --- | --- |
| open | Baru | `#1E40AF` | `#DBEAFE` |
| in_progress | Ditangani | `#92400E` | `#FEF3C7` |
| resolved | Selesai ditangani | `#166534` | `#DCFCE7` |
| closed | Ditutup | `#334155` | `#E2E8F0` |

Verifikasi contrast pada komponen final, termasuk hover, disabled, dan focus. Target teks normal minimal 4,5:1; teks besar dan batas kontrol relevan minimal 3:1. Nilai tersebut merujuk pada WCAG, bukan klaim bahwa seluruh antarmuka telah diaudit; lihat [REFERENCES.md](REFERENCES.md).

### 4.2 Tipografi dan ruang

- Font awal: sistem `ui-sans-serif, system-ui, sans-serif`; tidak bergantung pada pemuatan font eksternal.
- Judul halaman 28/36 px desktop, 24/32 px mobile; heading bagian 20/28 px.
- Body 16/24 px; label dan metadata 14/20 px. Nomor tiket menggunakan tabular numerals.
- Skala spacing 4, 8, 12, 16, 24, 32 px; kartu radius 12 px; input/tombol radius 8 px.
- Tombol utama dan field tinggi minimal 44 px sebagai target kenyamanan produk.
- Focus ring 2 px dengan offset 2 px; tidak boleh tertutup sticky header.
- Animasi pendek 120–180 ms dan mengikuti `prefers-reduced-motion`.

## 5. Struktur layout responsif

| Viewport | Shell dan layout |
| --- | --- |
| ≥ 1024 px | Sidebar 240 px, topbar 64 px, padding konten 24 px; detail dua kolom dengan panel informasi sekitar 320 px |
| 768–1023 px | Sidebar menjadi drawer; konten padding 20 px; detail menumpuk jika lebar kolom tidak cukup |
| 360–767 px | Drawer navigasi, konten padding 16 px, tombol utama full-width bila diperlukan; daftar tiket menjadi kartu |

Konten utama memiliki lebar maksimum 1440 px. Form pembuatan tiket maksimum 760 px. Jangan menyembunyikan status, judul, atau aksi utama pada mobile. Tabel admin boleh horizontal scroll dengan label yang tetap terbaca.

## 6. Spesifikasi layar

### D-01 — Login dan password

Form berisi email, password, tombol tampil/sembunyikan password, dan tombol “Masuk”. Error kredensial menggunakan satu pesan: “Email atau password tidak sesuai, atau akun tidak aktif.”

Saat submit, tombol loading dan request ganda dicegah. Login berhasil ke dashboard; jika `must_change_password=true`, arahkan ke change-password. Halaman tersebut menjelaskan aturan panjang password dan meminta password lama, password baru, serta konfirmasi lokal. Berhasil mengganti password mencabut sesi dan mengarahkan kembali ke login.

Tidak tersedia link registrasi publik atau “Kirim reset email” pada MVP. Informasi bantuan: “Hubungi admin untuk bantuan akun.”

### D-02 — Dashboard

Bagian atas: judul, keterangan scope, filter tanggal, dan tombol “Buat tiket”. Requester melihat “Ringkasan tiket saya”; staff melihat “Ringkasan layanan”.

Empat kartu menunjukkan status saat ini. Area berikutnya menampilkan backlog, waktu respons pertama, dan waktu penyelesaian. Distribusi kategori menggunakan bar horizontal sederhana dengan label jumlah; sediakan padanan tabel untuk aksesibilitas.

Klik kartu status menuju daftar dengan filter yang sesuai. Filter tanggal menggunakan tanggal pembuatan tiket dan dijelaskan pada tooltip. Tampilkan jumlah sampel di bawah metrik durasi. Empty metric menggunakan “Belum ada data”.

### D-03 — Daftar tiket

Toolbar berisi pencarian nomor/judul/deskripsi, filter status, prioritas, kategori, dan reset filter. Staff mendapat filter assignee serta shortcut “Ditugaskan kepada saya” dan “Belum ditugaskan”. Filter disimpan di URL agar navigasi kembali konsisten.

| Kolom desktop | Perilaku |
| --- | --- |
| Nomor dan judul | Link ke detail, maksimal dua baris judul |
| Status | Badge dengan label |
| Prioritas | Label Rendah/Normal/Tinggi/Mendesak |
| Kategori | Nama saat ini; label nonaktif jika relevan |
| Pelapor | Ditampilkan untuk staff |
| Penanggung jawab | Nama atau “Belum ditugaskan” |
| Dibuat | Tanggal WITA; waktu absolut tersedia |

Default sort terbaru berdasarkan created_at, lalu ID sebagai tie-breaker. Pagination default 20 item; pilihan 20, 50, 100. Pencarian dikirim setelah debounce sekitar 300 ms; request lama dibatalkan. UI membedakan “Belum ada tiket” dan “Tidak ada hasil yang sesuai filter”.

### D-04 — Form tiket baru

Urutan field: judul, kategori, prioritas, deskripsi, lampiran. Beri contoh deskripsi singkat: masalah, waktu kejadian, dan langkah yang sudah dicoba. Prioritas memiliki penjelasan dampak; urgent tidak menyatakan adanya jaminan SLA.

Lampiran dipilih lokal, diperiksa awal ukuran/jenis, kemudian tiket dibuat melalui JSON. Setelah tiket berhasil dibuat, unggah tiap file melalui endpoint lampiran. Jika beberapa unggahan gagal, arahkan ke detail dengan pesan “Tiket berhasil dibuat. Dua lampiran belum terunggah” dan tombol coba lagi untuk file terkait.

Jangan mengulang pembuatan tiket ketika unggahan gagal. Bila hasil create-ticket tidak diketahui akibat jaringan putus, minta pengguna memeriksa daftar tiket terbaru sebelum mengirim ulang; MVP belum menyediakan idempotency key.

### D-05 — Detail tiket

Header: nomor, judul, status, prioritas, waktu dibuat. Kolom utama: deskripsi, lampiran, ringkasan solusi jika ada, percakapan, dan timeline. Panel kanan: pelapor, kategori, assignee, tindakan sesuai izin.

Staff memiliki tab percakapan “Publik” dan “Internal”. Requester hanya menerima dan melihat publik. Composer internal memiliki label permanen; pergantian tab tidak langsung mengirim draft.

| Kondisi | Tindakan utama |
| --- | --- |
| Open, kosong, technician | “Ambil tiket” |
| Open, ditugaskan ke pengguna | “Mulai penanganan” |
| In progress, pengguna assignee | “Tandai selesai ditangani” |
| Resolved, pengguna reporter | “Konfirmasi selesai” dan aksi sekunder “Buka kembali” |
| Closed | Banner “Tiket telah ditutup” dan konten read-only |

Admin mendapat panel assignment/triage. Reassignment, perubahan prioritas, dan reopen meminta alasan. Resolve membuka form ringkasan solusi. Closed bukan state yang dapat diedit lewat dropdown bebas.

Timeline menampilkan label aktor dan waktu absolut; event internal memiliki label staff. Lampiran diunduh melalui API, tanpa preview inline pada MVP. Ikon file disertai nama dan ukuran.

Jika API mengembalikan konflik version, tampilkan: “Tiket berubah sejak halaman ini dibuka. Muat data terbaru sebelum menyimpan.” Pertahankan draft pengguna dan minta mereka meninjau ulang sebelum submit.

### D-06 — Notifikasi

Bell menampilkan unread_count. Item berisi ringkasan aman, nomor tiket, waktu, serta status dibaca. Klik item melakukan mark-read dan menuju detail; kegagalan mark-read tidak menghalangi navigasi.

Polling hanya berjalan ketika tab terlihat dan pengguna bersesi. Jeda saat offline, gunakan backoff setelah kegagalan, dan refetch saat tab kembali aktif. Pembaruan yang sudah terjadi tidak boleh memindahkan focus keyboard.

### D-07 — Administrasi

Akun: tabel nama/email/role/departemen/status, form buat akun, edit nama/departemen, aktivasi, dan reset password. Role ditampilkan read-only setelah akun dibuat. Password sementara hanya ditampilkan sekali setelah pembuatan/reset dan tidak disimpan oleh browser aplikasi.

Kategori/departemen: nama, deskripsi opsional, status aktif. Gunakan “Nonaktifkan”, dengan penjelasan bahwa data historis tetap ada. Tidak ada tombol hapus.

Audit: tabel waktu/aktor/aksi/target/request ID; filter dan pagination server. Jangan menampilkan credential, token, atau body sensitif.

## 7. Komponen reusable

| Komponen | Tanggung jawab |
| --- | --- |
| AppShell, Sidebar, Topbar | Navigasi dan layout |
| Button, Input, Select, Textarea | Interaksi dasar dan state standar |
| FormField, ErrorSummary | Label, bantuan, validasi, fokus pada error |
| TicketStatusBadge, PriorityBadge | Representasi konsisten status/prioritas |
| TicketTable, TicketCard | List desktop dan mobile |
| TicketFilters, Pagination | Query state dan navigasi hasil |
| CommentComposer, CommentList | Draft, visibility, submit dan daftar |
| AttachmentPicker, AttachmentList | Batas berkas dan hasil per file |
| ActivityTimeline | Event terotorisasi |
| ConfirmDialog, ResolutionDialog | Konfirmasi dan input keputusan |
| LoadingState, EmptyState, ErrorState | Kondisi sistem |

Mulai dengan komponen yang digunakan nyata. Pengelompokan “primitives / shared / features” cukup; hierarki atomic design lengkap bukan syarat proyek.

## 8. State dan aksesibilitas

Semua halaman memiliki loading, empty, error, forbidden/not-found, dan success sesuai konteks. Form memiliki dirty, submitting, invalid, dan conflict state. Pesan sukses singkat melalui aria-live; error form tetap terlihat sampai diperbaiki.

Gunakan elemen semantik, label eksplisit, urutan tab logis, dialog dengan pengelolaan focus, dan fokus kembali ke pemicu setelah dialog ditutup. Toast tidak menjadi satu-satunya lokasi error. Icon-only button memiliki accessible name.

Verifikasi pada lebar 360, 768, 1024, dan 1440 px; zoom 200%; keyboard-only; reduced motion; serta screen reader untuk login, create-ticket, dan detail. Catat temuan pada TEST_PLAN, jangan menyatakan kepatuhan WCAG tanpa pemeriksaan.
