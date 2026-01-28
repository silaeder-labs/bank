-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE transactions (
    line_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    amount_cents BIGINT NOT NULL,
    description VARCHAR(100)
);

CREATE TABLE balances (
    user_id UUID PRIMARY KEY,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    amount_cents BIGINT NOT NULL
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS
'BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;'
LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at_balance
BEFORE UPDATE ON balances
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at_balance ON balances;
DROP TRIGGER IF EXISTS set_updated_at ON transactions;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS transactions;