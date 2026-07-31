package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrConnNotInContext = errors.New("db: connection not found in context")

type connKey struct{}

// Conn is a database connection or transaction already opened for the request.
type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func WithConn(ctx context.Context, conn Conn) context.Context {
	return context.WithValue(ctx, connKey{}, conn)
}

func ConnFromContext(ctx context.Context) (Conn, error) {
	conn, ok := ctx.Value(connKey{}).(Conn)
	if !ok || conn == nil {
		return nil, ErrConnNotInContext
	}

	return conn, nil
}
