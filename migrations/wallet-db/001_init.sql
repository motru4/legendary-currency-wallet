-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(64) NOT NULL UNIQUE,
    email         VARCHAR(256) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TYPE currency_code AS ENUM ('USD', 'RUB', 'EUR');

CREATE TABLE IF NOT EXISTS wallet_balances (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency  currency_code NOT NULL,
    balance   NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT wallet_balances_user_currency_uniq UNIQUE (user_id, currency)
);

CREATE TYPE operation_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'EXCHANGE');
CREATE TYPE operation_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED');

CREATE TABLE IF NOT EXISTS operations (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type                operation_type NOT NULL,
    status              operation_status NOT NULL DEFAULT 'PENDING',
    amount              NUMERIC(18,2) NOT NULL,
    from_currency       currency_code,
    to_currency         currency_code,
    rate                NUMERIC(18,6),
    exchanged_amount    NUMERIC(18,2),
    balance_after_usd   NUMERIC(18,2),
    balance_after_rub   NUMERIC(18,2),
    balance_after_eur   NUMERIC(18,2),
    error               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_user_id ON operations(user_id);
CREATE INDEX IF NOT EXISTS idx_operations_status ON operations(status);
CREATE INDEX IF NOT EXISTS idx_operations_type   ON operations(type);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS operations;
DROP TYPE IF EXISTS operation_status;
DROP TYPE IF EXISTS operation_type;
DROP TABLE IF EXISTS wallet_balances;
DROP TYPE IF EXISTS currency_code;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
