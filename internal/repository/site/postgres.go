package site

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/db"
	domain "github.com/ViktorNikolaevichD/site-monitor/internal/domain/site"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresSiteRepository struct {
	logger *slog.Logger
}

func NewPostgresSiteRepository(logger *slog.Logger) *PostgresSiteRepository {
	return &PostgresSiteRepository{
		logger: logger,
	}
}

func (r *PostgresSiteRepository) GetAll(ctx context.Context) ([]domain.Site, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("get all sites: connection not in context", "error", err)
		return nil, domain.ErrStorage
	}

	rows, err := conn.Query(ctx, `
		SELECT id, url, name
		FROM sites
		WHERE is_active
		ORDER BY created_at
	`)
	if err != nil {
		r.logger.Error("get all sites: query failed", "error", err)
		return nil, domain.ErrStorage
	}
	defer rows.Close()

	sites := make([]domain.Site, 0)
	for rows.Next() {
		var (
			id   uuid.UUID
			url  string
			name string
		)

		if err := rows.Scan(&id, &url, &name); err != nil {
			r.logger.Error("get all sites: scan failed", "error", err)
			return nil, domain.ErrStorage
		}

		sites = append(sites, domain.Site{
			ID:   id,
			URL:  url,
			Name: name,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("get all sites: rows iteration failed", "error", err)
		return nil, domain.ErrStorage
	}

	return sites, nil
}

func (r *PostgresSiteRepository) AddIfAbsent(ctx context.Context, site domain.Site) (domain.Site, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("add site: connection not in context", "error", err)
		return domain.Site{}, domain.ErrStorage
	}

	now := time.Now().UTC()
	tag, err := conn.Exec(ctx, `
		INSERT INTO sites (id, url, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO NOTHING
	`, site.ID, site.URL, site.Name, now, now)
	if err != nil {
		r.logger.Error("add site: insert failed", "error", err, "url", site.URL)
		return domain.Site{}, domain.ErrStorage
	}

	if tag.RowsAffected() == 0 {
		return domain.Site{}, domain.ErrSiteAlreadyExists
	}

	return site, nil
}

func (r *PostgresSiteRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("delete site: connection not in context", "error", err)
		return domain.ErrStorage
	}

	tag, err := conn.Exec(ctx, `
		DELETE FROM sites
		WHERE id = $1
	`, id)
	if err != nil {
		r.logger.Error("delete site: query failed", "error", err, "site_id", id)
		return domain.ErrStorage
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrSiteNotFound
	}

	return nil
}

func (r *PostgresSiteRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (domain.Site, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("get site by id for update: connection not in context", "error", err)
		return domain.Site{}, domain.ErrStorage
	}

	var (
		url  string
		name string
	)

	err = conn.QueryRow(ctx, `
		SELECT url, name
		FROM sites
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(&url, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Site{}, domain.ErrSiteNotFound
		}
		r.logger.Error("get site by id for update: query failed", "error", err, "site_id", id)
		return domain.Site{}, domain.ErrStorage
	}

	return domain.Site{
		ID:   id,
		URL:  url,
		Name: name,
	}, nil
}

func (r *PostgresSiteRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Site, error) {
	conn, err := db.ConnFromContext(ctx)
	if err != nil {
		r.logger.Error("get site by id: connection not in context", "error", err)
		return domain.Site{}, domain.ErrStorage
	}

	var (
		url  string
		name string
	)

	err = conn.QueryRow(ctx, `
		SELECT url, name
		FROM sites
		WHERE id = $1
	`, id).Scan(&url, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Site{}, domain.ErrSiteNotFound
		}
		r.logger.Error("get site by id: query failed", "error", err, "site_id", id)
		return domain.Site{}, domain.ErrStorage
	}

	return domain.Site{
		ID:   id,
		URL:  url,
		Name: name,
	}, nil
}
