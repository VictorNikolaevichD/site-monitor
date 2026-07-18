package db

import (
	"context"
	"database/sql"
)

type PostgresChecker struct {
	db *sql.DB
}

func NewPostgresChecker(db *sql.DB) *PostgresChecker {
	return &PostgresChecker{db: db}
}

func (c *PostgresChecker) Name() string {
	return "postgres"
}

func (c *PostgresChecker) Check(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
