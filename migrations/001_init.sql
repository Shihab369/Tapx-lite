-- accounts: stores user balances for the closed-loop payment system.
CREATE TABLE IF NOT EXISTS accounts (
    id         SERIAL PRIMARY KEY,
    balance    INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- transactions: immutable ledger of all processed payments.
CREATE TABLE IF NOT EXISTS transactions (
    id         TEXT PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES accounts(id),
    amount     INT NOT NULL CHECK (amount > 0),
    nonce      TEXT NOT NULL UNIQUE,
    status     TEXT NOT NULL DEFAULT 'accepted',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- seed: insert two test accounts for local development.
INSERT INTO accounts (id, balance) VALUES (1, 10000) ON CONFLICT DO NOTHING;
INSERT INTO accounts (id, balance) VALUES (2, 5000)  ON CONFLICT DO NOTHING;