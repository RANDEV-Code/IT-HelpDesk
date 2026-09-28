-- 000001_baseline.up.sql — DDL baseline MVP RANDesk
-- Sumber: RANDesk_Planning/SCHEMA.md §4, disalin apa adanya.
-- Deviasi tercatat (GAPS_AND_DECISIONS.md): pembungkus BEGIN/COMMIT dari
-- dokumen sumber dilepas karena golang-migrate menjalankan setiap migrasi
-- dalam satu transaksi; isi DDL tidak diubah.

CREATE TABLE departments (
    id UUID PRIMARY KEY,
    name VARCHAR(80) NOT NULL CHECK (char_length(btrim(name)) BETWEEN 2 AND 80),
    description VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX departments_name_ci_uq ON departments (lower(name));

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(80) NOT NULL CHECK (char_length(btrim(name)) BETWEEN 2 AND 80),
    description VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX categories_name_ci_uq ON categories (lower(name));

CREATE TABLE users (
    id UUID PRIMARY KEY,
    department_id UUID REFERENCES departments(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL CHECK (char_length(btrim(name)) BETWEEN 2 AND 100),
    email VARCHAR(254) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('requester', 'technician', 'admin')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    must_change_password BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_email_normalized_ck CHECK (
        email = lower(btrim(email)) AND char_length(email) BETWEEN 3 AND 254
    )
);
CREATE INDEX users_department_idx ON users (department_id);
CREATE INDEX users_role_active_idx ON users (role, is_active);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    token_hash CHAR(64) NOT NULL UNIQUE CHECK (token_hash ~ '^[0-9a-f]{64}$'),
    csrf_token VARCHAR(43) NOT NULL CHECK (csrf_token ~ '^[A-Za-z0-9_-]{43}$'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    CONSTRAINT sessions_expiry_ck CHECK (expires_at > created_at),
    CONSTRAINT sessions_revoked_ck CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);
CREATE INDEX sessions_user_idx ON sessions (user_id);
CREATE INDEX sessions_expiry_idx ON sessions (expires_at);

CREATE TABLE tickets (
    id UUID PRIMARY KEY,
    ticket_no BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    requester_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assignee_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(150) NOT NULL CHECK (char_length(btrim(title)) BETWEEN 5 AND 150),
    description TEXT NOT NULL CHECK (char_length(btrim(description)) BETWEEN 20 AND 10000),
    priority TEXT NOT NULL DEFAULT 'normal'
        CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status TEXT NOT NULL DEFAULT 'open'
        CHECK (status IN ('open', 'in_progress', 'resolved', 'closed')),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    first_response_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES users(id) ON DELETE RESTRICT,
    resolution_summary TEXT,
    closed_at TIMESTAMPTZ,
    closed_by UUID REFERENCES users(id) ON DELETE RESTRICT,
    reopen_count INTEGER NOT NULL DEFAULT 0 CHECK (reopen_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tickets_assignee_state_ck CHECK (
        status = 'open' OR assignee_id IS NOT NULL
    ),
    CONSTRAINT tickets_resolution_length_ck CHECK (
        resolution_summary IS NULL OR char_length(btrim(resolution_summary)) BETWEEN 20 AND 4000
    ),
    CONSTRAINT tickets_lifecycle_ck CHECK (
        (status IN ('open', 'in_progress')
            AND resolved_at IS NULL AND resolved_by IS NULL
            AND resolution_summary IS NULL AND closed_at IS NULL AND closed_by IS NULL)
        OR
        (status = 'resolved'
            AND resolved_at IS NOT NULL AND resolved_by IS NOT NULL
            AND resolution_summary IS NOT NULL AND closed_at IS NULL AND closed_by IS NULL)
        OR
        (status = 'closed'
            AND resolved_at IS NOT NULL AND resolved_by IS NOT NULL
            AND resolution_summary IS NOT NULL AND closed_at IS NOT NULL AND closed_by IS NOT NULL)
    ),
    CONSTRAINT tickets_time_order_ck CHECK (
        updated_at >= created_at
        AND (first_response_at IS NULL OR first_response_at >= created_at)
        AND (resolved_at IS NULL OR resolved_at >= created_at)
        AND (closed_at IS NULL OR closed_at >= resolved_at)
    )
);
CREATE INDEX tickets_requester_created_idx ON tickets (requester_id, created_at DESC, id DESC);
CREATE INDEX tickets_assignee_status_idx ON tickets (assignee_id, status, created_at DESC, id DESC);
CREATE INDEX tickets_status_created_idx ON tickets (status, created_at DESC, id DESC);
CREATE INDEX tickets_category_created_idx ON tickets (category_id, created_at DESC, id DESC);
CREATE INDEX tickets_created_idx ON tickets (created_at DESC, id DESC);

CREATE TABLE ticket_comments (
    id UUID PRIMARY KEY,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'internal')),
    body TEXT NOT NULL CHECK (char_length(btrim(body)) BETWEEN 1 AND 5000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX comments_ticket_visibility_idx
    ON ticket_comments (ticket_id, visibility, created_at DESC, id DESC);
CREATE INDEX comments_author_idx ON ticket_comments (author_id);

CREATE TABLE ticket_attachments (
    id UUID PRIMARY KEY,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    original_name VARCHAR(180) NOT NULL CHECK (char_length(original_name) BETWEEN 1 AND 180),
    storage_key TEXT NOT NULL UNIQUE,
    content_type TEXT NOT NULL CHECK (content_type IN ('image/jpeg', 'image/png', 'application/pdf')),
    size_bytes BIGINT NOT NULL CHECK (size_bytes BETWEEN 1 AND 5242880),
    sha256 CHAR(64) NOT NULL CHECK (sha256 ~ '^[0-9a-f]{64}$'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT attachments_delete_pair_ck CHECK (
        (deleted_at IS NULL AND deleted_by IS NULL)
        OR (deleted_at IS NOT NULL AND deleted_by IS NOT NULL AND deleted_at >= created_at)
    )
);
CREATE INDEX attachments_ticket_active_idx ON ticket_attachments (ticket_id, created_at)
    WHERE deleted_at IS NULL;
CREATE INDEX attachments_uploader_idx ON ticket_attachments (uploaded_by);

CREATE TABLE ticket_events (
    id UUID PRIMARY KEY,
    seq BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    actor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    event_type TEXT NOT NULL CHECK (event_type IN (
        'ticket.created', 'ticket.updated', 'ticket.assigned', 'ticket.reassigned',
        'ticket.unassigned', 'ticket.started', 'ticket.resolved', 'ticket.closed',
        'ticket.reopened', 'ticket.comment_added', 'ticket.attachment_added',
        'ticket.attachment_deleted'
    )),
    visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'internal')),
    payload JSONB NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(payload) = 'object'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX events_ticket_seq_idx ON ticket_events (ticket_id, seq DESC);
CREATE INDEX events_actor_idx ON ticket_events (actor_id);

CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    recipient_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    event_id UUID NOT NULL REFERENCES ticket_events(id) ON DELETE RESTRICT,
    message VARCHAR(200) NOT NULL CHECK (char_length(message) BETWEEN 1 AND 200),
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT notifications_recipient_event_uq UNIQUE (recipient_id, event_id),
    CONSTRAINT notifications_read_time_ck CHECK (read_at IS NULL OR read_at >= created_at)
);
CREATE INDEX notifications_recipient_created_idx
    ON notifications (recipient_id, created_at DESC, id DESC);
CREATE INDEX notifications_unread_idx ON notifications (recipient_id)
    WHERE read_at IS NULL;
CREATE INDEX notifications_event_idx ON notifications (event_id);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    actor_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    action VARCHAR(80) NOT NULL,
    target_type VARCHAR(40) NOT NULL,
    target_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(metadata) = 'object'),
    request_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX audit_created_idx ON audit_logs (created_at DESC, id DESC);
CREATE INDEX audit_actor_created_idx ON audit_logs (actor_id, created_at DESC);
CREATE INDEX audit_target_idx ON audit_logs (target_type, target_id);
