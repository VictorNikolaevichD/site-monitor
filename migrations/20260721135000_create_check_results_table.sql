-- +goose Up
CREATE TABLE check_results (
    id           BIGSERIAL PRIMARY KEY,
    site_id      UUID NOT NULL REFERENCES sites (id) ON DELETE CASCADE,
    http_code    INT NOT NULL,
    duration_ns  BIGINT NOT NULL,
    availability BOOLEAN NOT NULL,
    error        TEXT NOT NULL DEFAULT '',
    checked_at   TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE check_results;
