-- 000001_baseline.down.sql — hapus bersih DDL baseline (urutan dependensi terbalik).
-- Hanya untuk development/testing; di produksi down migration destruktif tidak
-- dijalankan otomatis (OPERATIONS §4).

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS ticket_events;
DROP TABLE IF EXISTS ticket_attachments;
DROP TABLE IF EXISTS ticket_comments;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS departments;
