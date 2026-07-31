package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
func WithinTx(ctx context.Context, beginner TxBeginner, logger *slog.Logger, fn func(ctx context.Context) error) error {
	return execTx(ctx, beginner, logger, fn)
}

// WithinTxResult is like WithinTx but returns a value from fn.
func WithinTxResult[T any](
	ctx context.Context,
	beginner TxBeginner,
	logger *slog.Logger,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var result T
	err := execTx(ctx, beginner, logger, func(ctx context.Context) error {
		var err error
		result, err = fn(ctx)
		return err
	})
	return result, err
}

func execTx(ctx context.Context, beginner TxBeginner, logger *slog.Logger, fn func(ctx context.Context) error) error {
	logger.Info("transaction started")

	tx, err := beginner.Begin(ctx)
	if err != nil {
		logger.Error("transaction begin failed", "error", err)
		return fmt.Errorf("%w: begin: %v", ErrTx, err)
	}

	committed := false
	defer func() {
		if committed {
			return
		}

		if err := tx.Rollback(ctx); err != nil {
			logger.Error("transaction rollback failed", "error", err)
			return
		}

		logger.Info("transaction rolled back")
	}()

	if err := fn(WithConn(ctx, tx)); err != nil {
		logger.Info("transaction operation failed, rolling back", "error", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("transaction commit failed", "error", err)
		return fmt.Errorf("%w: commit: %v", ErrTx, err)
	}

	committed = true
	logger.Info("transaction committed")
	return nil
}

var _ TxBeginner = (*pgxpool.Pool)(nil)
