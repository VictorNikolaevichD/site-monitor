package checkresult

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type PostgresCheckResultRepository struct {
	logger *slog.Logger
}

func NewPostgresCheckResultRepository(logger *slog.Logger) *PostgresCheckResultRepository {
	return &PostgresCheckResultRepository{
		logger: logger,
	}
}

func (r *PostgresCheckResultRepository) Create(
	ctx context.Context,
	siteID uuid.UUID,
	result domain.CheckStatus,
) error {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("create check result: connection not in context", "error", err)
		return domain.ErrStorage
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO check_results (site_id, http_code, duration_ns, availability, error, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		siteID,
		result.Code,
		result.Duration.Nanoseconds(),
		result.Availability == domain.Available,
		result.Error,
		result.CheckedAt,
	)
	if err != nil {
		r.logger.Error("create check result: insert failed", "error", err, "site_id", siteID)
		return domain.ErrStorage
	}

	return nil
}

func (r *PostgresCheckResultRepository) GetBySiteID(
	ctx context.Context,
	siteID uuid.UUID,
	limit, offset int,
) ([]domain.CheckResult, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("get check results by site id: connection not in context", "error", err)
		return nil, domain.ErrStorage
	}

	rows, err := conn.Query(ctx, `
		SELECT id, http_code, duration_ns, availability, error, checked_at
		FROM check_results
		WHERE site_id = $1
		ORDER BY checked_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, siteID, limit, offset)
	if err != nil {
		r.logger.Error("get check results by site id: query failed", "error", err, "site_id", siteID)
		return nil, domain.ErrStorage
	}
	defer rows.Close()

	results := make([]domain.CheckResult, 0)
	for rows.Next() {
		result, err := scanCheckResult(rows)
		if err != nil {
			r.logger.Error("get check results by site id: scan failed", "error", err, "site_id", siteID)
			return nil, domain.ErrStorage
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("get check results by site id: rows iteration failed", "error", err, "site_id", siteID)
		return nil, domain.ErrStorage
	}

	return results, nil
}

func (r *PostgresCheckResultRepository) CountBySiteID(ctx context.Context, siteID uuid.UUID) (int, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("count check results by site id: connection not in context", "error", err)
		return 0, domain.ErrStorage
	}

	var total int
	err = conn.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM check_results
		WHERE site_id = $1
	`, siteID).Scan(&total)
	if err != nil {
		r.logger.Error("count check results by site id: query failed", "error", err, "site_id", siteID)
		return 0, domain.ErrStorage
	}

	return total, nil
}

func (r *PostgresCheckResultRepository) GetLatestBySiteID(
	ctx context.Context,
	siteID uuid.UUID,
) (domain.CheckStatus, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("get latest check result: connection not in context", "error", err)
		return domain.CheckStatus{}, domain.ErrStorage
	}

	row := conn.QueryRow(ctx, `
		SELECT http_code, duration_ns, availability, error, checked_at
		FROM check_results
		WHERE site_id = $1
		ORDER BY checked_at DESC, id DESC
		LIMIT 1
	`, siteID)

	status, err := scanCheckStatus(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CheckStatus{}, domain.ErrCheckResultNotFound
		}
		r.logger.Error("get latest check result: query failed", "error", err, "site_id", siteID)
		return domain.CheckStatus{}, domain.ErrStorage
	}

	return status, nil
}

func (r *PostgresCheckResultRepository) DeleteBySiteID(ctx context.Context, siteID uuid.UUID) error {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("delete check results by site id: connection not in context", "error", err)
		return domain.ErrStorage
	}

	_, err = conn.Exec(ctx, `
		DELETE FROM check_results
		WHERE site_id = $1
	`, siteID)
	if err != nil {
		r.logger.Error("delete check results by site id: query failed", "error", err, "site_id", siteID)
		return domain.ErrStorage
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCheckResult(row scanner) (domain.CheckResult, error) {
	var (
		id           int64
		httpCode     int
		durationNs   int64
		availability bool
		checkError   string
		checkedAt    time.Time
	)

	if err := row.Scan(&id, &httpCode, &durationNs, &availability, &checkError, &checkedAt); err != nil {
		return domain.CheckResult{}, err
	}

	return domain.CheckResult{
		ID:           id,
		Availability: availabilityFromBool(availability),
		Code:         httpCode,
		CheckedAt:    checkedAt,
		Duration:     time.Duration(durationNs),
		Error:        checkError,
	}, nil
}

func scanCheckStatus(row scanner) (domain.CheckStatus, error) {
	var (
		httpCode     int
		durationNs   int64
		availability bool
		checkError   string
		checkedAt    time.Time
	)

	if err := row.Scan(&httpCode, &durationNs, &availability, &checkError, &checkedAt); err != nil {
		return domain.CheckStatus{}, err
	}

	return domain.CheckStatus{
		Availability: availabilityFromBool(availability),
		Code:         httpCode,
		CheckedAt:    checkedAt,
		Duration:     time.Duration(durationNs),
		Error:        checkError,
	}, nil
}

func availabilityFromBool(availability bool) domain.Availability {
	if availability {
		return domain.Available
	}
	return domain.Unavailable
}
