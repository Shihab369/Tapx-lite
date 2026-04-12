
-- Initial Schema
-- PostgreSQL 15
-- ROLES
-- Defines access levels for system users (RBAC foundation).
CREATE TABLE roles (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO roles (name) VALUES
    ('super_admin'),
    ('admin'),
    ('operator')
ON CONFLICT DO NOTHING;

-- USERS
-- System and admin users. Separate from financial accounts.

CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    role_id    INT NOT NULL REFERENCES roles(id),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users (name, email, role_id) VALUES
    ('Shihab', 'shihab@tapx.com', 1),
    ('Admin User', 'admin@tapx.com', 2)
ON CONFLICT DO NOTHING;



-- ACCOUNTS
-- Financial entity. One account per user.
-- Separates identity (users) from financial state (accounts).

CREATE TABLE accounts (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance    INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO accounts (user_id, balance) VALUES
    (1, 10000),
    (2, 5000)
ON CONFLICT DO NOTHING;



-- TRANSACTIONS (partitioned by month)
-- Core payment ledger. Immutable after insert.
-- Partitioned by created_at for performance at scale.

CREATE TABLE transactions (
    id         TEXT NOT NULL,
    user_id    INT NOT NULL REFERENCES users(id),
    amount     INT NOT NULL CHECK (amount > 0),
    nonce      TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'accepted',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE transactions_2026_04
    PARTITION OF transactions
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

CREATE TABLE transactions_2026_05
    PARTITION OF transactions
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');



-- SYNC BATCHES
-- Tracks each device sync request.
-- Enables observability: how many accepted/rejected per sync.

CREATE TABLE sync_batches (
    id              SERIAL PRIMARY KEY,
    device_id       TEXT NOT NULL,
    total_submitted INT NOT NULL DEFAULT 0,
    total_accepted  INT NOT NULL DEFAULT 0,
    total_rejected  INT NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'processing',
    started_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMP
);



-- SYNC TRANSACTIONS
-- Granular per-transaction record inside a sync batch.
-- Enables debugging individual failures.

CREATE TABLE sync_transactions (
    id             SERIAL PRIMARY KEY,
    batch_id       INT NOT NULL REFERENCES sync_batches(id) ON DELETE CASCADE,
    transaction_id TEXT,
    status         TEXT,
    error_message  TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);



-- AUDIT LOG
-- Immutable record of all data changes across tables.
-- Required for security, compliance, and debugging.

CREATE TABLE audit_log (
    id          SERIAL PRIMARY KEY,
    table_name  TEXT NOT NULL,
    operation   TEXT NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    row_id      TEXT NOT NULL,
    changed_by  INT REFERENCES users(id),
    changed_at  TIMESTAMP NOT NULL DEFAULT NOW()
);



-- INDEXES
-- Applied on frequently queried columns.


CREATE INDEX idx_users_role_id
    ON users(role_id);

CREATE INDEX idx_accounts_user_id
    ON accounts(user_id);

CREATE INDEX idx_transactions_user_id
    ON transactions(user_id);

CREATE INDEX idx_transactions_nonce
    ON transactions(nonce);

CREATE INDEX idx_transactions_created_at
    ON transactions(created_at);

CREATE INDEX idx_sync_batches_device_id
    ON sync_batches(device_id);

CREATE INDEX idx_sync_transactions_batch_id
    ON sync_transactions(batch_id);

CREATE INDEX idx_audit_log_table_operation
    ON audit_log(table_name, operation);