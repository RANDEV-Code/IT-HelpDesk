-- check-constraints-task003.sql — bukti CHECK constraint baseline (TASK-003).
-- Dijalankan sebagai randesk_migrate pada randesk_dev; seluruh perubahan
-- di-ROLLBACK sehingga DB tetap bersih. Empat INSERT pertama HARUS gagal dengan
-- nama constraint yang dicantumkan di komentar; kontrol positif HARUS berhasil.
\set ON_ERROR_STOP off
BEGIN;

INSERT INTO departments (id, name) VALUES ('11111111-1111-1111-1111-111111111111', 'Dev Dept');
INSERT INTO categories (id, name) VALUES ('22222222-2222-2222-2222-222222222222', 'Jaringan');
INSERT INTO users (id, name, email, password_hash, role, must_change_password)
VALUES ('33333333-3333-3333-3333-333333333333', 'Requester Satu', 'r1@example.com', 'x', 'requester', false);
INSERT INTO users (id, name, email, password_hash, role, must_change_password)
VALUES ('44444444-4444-4444-4444-444444444444', 'Teknisi Satu', 't1@example.com', 'x', 'technician', false);

-- TES 1: status invalid → HARUS ditolak. Catatan: PostgreSQL mengevaluasi
-- tickets_lifecycle_ck sebelum tickets_status_check, sehingga status di luar
-- enum selalu ditolak lifecycle_ck lebih dulu; enum status tetap tegak.
SAVEPOINT s1;
INSERT INTO tickets (id, requester_id, assignee_id, category_id, title, description, status)
VALUES (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444',
        '22222222-2222-2222-2222-222222222222', 'judul cukup panjang',
        'deskripsi minimal dua puluh karakter xyz', 'bogus');
ROLLBACK TO SAVEPOINT s1;

-- TES 2: version 0 → HARUS ditolak tickets_version_check
SAVEPOINT s2;
INSERT INTO tickets (id, requester_id, assignee_id, category_id, title, description, version)
VALUES (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444',
        '22222222-2222-2222-2222-222222222222', 'judul cukup panjang',
        'deskripsi minimal dua puluh karakter xyz', 0);
ROLLBACK TO SAVEPOINT s2;

-- TES 3: email tidak lowercase → HARUS ditolak users_email_normalized_ck
SAVEPOINT s3;
INSERT INTO users (id, name, email, password_hash, role)
VALUES (gen_random_uuid(), 'Test User', 'Upper@Example.COM', 'x', 'requester');
ROLLBACK TO SAVEPOINT s3;

-- TES 4: description terlalu pendek → HARUS ditolak tickets_description_check
SAVEPOINT s4;
INSERT INTO tickets (id, requester_id, category_id, title, description)
VALUES (gen_random_uuid(), '33333333-3333-3333-3333-333333333333',
        '22222222-2222-2222-2222-222222222222', 'judul cukup panjang', 'pendek');
ROLLBACK TO SAVEPOINT s4;

-- KONTROL POSITIF: tiket open valid tanpa assignee → HARUS berhasil (INSERT 0 1)
SAVEPOINT s5;
INSERT INTO tickets (id, requester_id, category_id, title, description)
VALUES (gen_random_uuid(), '33333333-3333-3333-3333-333333333333',
        '22222222-2222-2222-2222-222222222222', 'judul cukup panjang',
        'deskripsi minimal dua puluh karakter xyz');
ROLLBACK TO SAVEPOINT s5;

ROLLBACK;
