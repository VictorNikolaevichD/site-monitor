package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTx = errors.New("db: transaction failed")

type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// WithinTx runs fn inside a transaction.
// On success the transaction is committed; on error or panic it is rolled back.
// Begin/Commit failures are returned as ErrTx.
func WithinTx(ctx context.Context, beginner TxBeginner, fn func(ctx context.Context) error) error {
	tx, err := beginner.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: begin: %v", ErrTx, err)
	}
	defer tx.Rollback(ctx)

	if err := fn(WithConn(ctx, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: commit: %v", ErrTx, err)
	}

	return nil
}

// WithinTxResult is like WithinTx but returns a value from fn.
func WithinTxResult[T any](ctx context.Context, beginner TxBeginner, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T

	tx, err := beginner.Begin(ctx)
	if err != nil {
		return zero, fmt.Errorf("%w: begin: %v", ErrTx, err)
	}
	defer tx.Rollback(ctx)

	result, err := fn(WithConn(ctx, tx))
	if err != nil {
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return zero, fmt.Errorf("%w: commit: %v", ErrTx, err)
	}

	return result, nil
}

var _ TxBeginner = (*pgxpool.Pool)(nil)
