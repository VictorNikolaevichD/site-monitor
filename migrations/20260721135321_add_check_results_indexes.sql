-- +goose Up
CREATE INDEX idx_check_results_site_id_checked_at ON check_results (site_id, checked_at DESC);
CREATE INDEX idx_check_results_checked_at ON check_results (checked_at DESC);

-- +goose Down
DROP INDEX idx_check_results_checked_at;
DROP INDEX idx_check_results_site_id_checked_at;
