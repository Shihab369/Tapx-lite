# Database Schema - TapX Lite

PostgreSQL 15. 9 tables. Partitioned transactions.

---

## Tables

### roles

Defines system access levels.

```sql
id, name (unique), created_at
```

Seed: `super_admin`, `admin`, `operator`

---

### users

System/admin users. Linked to a role.

```sql
id, name, email (unique), role_id → roles, is_active, created_at, updated_at
```

---

### accounts

Financial entity. One account per user.

```sql
id, user_id → users (unique), balance (>= 0), created_at
```

---

### transactions _(partitioned by month)_

Core payment ledger. Immutable after insert.

```sql
id, user_id → users, amount (> 0), nonce, status, created_at
PRIMARY KEY (id, created_at)
```

Partitions: `transactions_2026_04`, `transactions_2026_05`

---

### sync_batches

Tracks each device sync request.

```sql
id, device_id, total_submitted, total_accepted, total_rejected, status, started_at, finished_at
```

---

### sync_transactions

Per-transaction record inside a batch.

```sql
id, batch_id → sync_batches, transaction_id, status, error_message, created_at
```

---

### audit_log

Records all data changes across tables.

```sql
id, table_name, operation (INSERT/UPDATE/DELETE), row_id, changed_by → users, changed_at
```

---

## Indexes

| Index                          | Table             | Column(s)             |
| ------------------------------ | ----------------- | --------------------- |
| idx_users_role_id              | users             | role_id               |
| idx_accounts_user_id           | accounts          | user_id               |
| idx_transactions_user_id       | transactions      | user_id               |
| idx_transactions_nonce         | transactions      | nonce                 |
| idx_transactions_created_at    | transactions      | created_at            |
| idx_sync_batches_device_id     | sync_batches      | device_id             |
| idx_sync_transactions_batch_id | sync_transactions | batch_id              |
| idx_audit_log_table_operation  | audit_log         | table_name, operation |

---

## Key Concepts Applied

- **RBAC** roles + users
- **Partitioning** transactions split by month
- **Constraints** CHECK, UNIQUE, FOREIGN KEY
- **Audit trail** audit_log tracks all changes
- **Observability** sync_batches tracks device sync behavior
- **Indexes** on all frequently queried columns

---

## Run Migration

```bash
docker exec -i tapx-postgres psql -U postgres -d tapx < migrations/001_init.sql
```
