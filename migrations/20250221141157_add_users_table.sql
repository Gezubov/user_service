-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users
ADD COLUMN uuid UUID DEFAULT gen_random_uuid() NOT NULL;

UPDATE users
SET uuid = gen_random_uuid()
WHERE uuid IS NULL;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_pkey;

ALTER TABLE users
ADD PRIMARY KEY (uuid);

CREATE UNIQUE INDEX IF NOT EXISTS users_id_unique ON users(id);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS set_updated_at ON users;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd