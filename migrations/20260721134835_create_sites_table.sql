-- +goose Up
CREATE TABLE sites (
    id         UUID PRIMARY KEY,
    url        TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE sites;
