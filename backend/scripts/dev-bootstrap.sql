-- dev-bootstrap.sql — disiapkan oleh TASK-003; jalankan sebagai superuser postgres:
--   psql -U postgres -h 127.0.0.1 -f backend/scripts/dev-bootstrap.sql
-- Membuat role dan database development/testing. Kredensial di bawah adalah
-- placeholder DEVELOPMENT SAJA (bukan secret produksi); nilai sebenarnya untuk
-- mesin dev disimpan di .env (tidak di-commit). Grants runtime penuh menyusul
-- di TASK-027 (SCHEMA §5: role runtime terpisah dari role migrasi).

-- Role migrasi: pemilik objek schema di DB dev/test.
DO $$ BEGIN
    CREATE ROLE randesk_migrate LOGIN PASSWORD 'randesk_migrate_dev';
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Role runtime: dipakai proses API (grants dibatasi di TASK-027).
DO $$ BEGIN
    CREATE ROLE randesk_runtime LOGIN PASSWORD 'randesk_runtime_dev';
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Role operator retention/purge (credential terpisah, OPERATIONS §7; dipakai TASK-052).
DO $$ BEGIN
    CREATE ROLE randesk_operator NOLOGIN;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
