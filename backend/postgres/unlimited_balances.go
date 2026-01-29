package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nrf24l01/go-web-utils/pgkit"
)

func HasUnlimitedBalance(db *pgkit.DB, ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM unlimited_balances
			WHERE user_id = $1 AND deleted_at IS NULL
		)
	`, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GrantUnlimitedBalance(db *pgkit.DB, ctx context.Context, userID uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO unlimited_balances (user_id, deleted_at)
		VALUES ($1, NULL)
		ON CONFLICT (user_id) DO UPDATE
		SET deleted_at = NULL
	`, userID)
	return err
}

func RevokeUnlimitedBalance(db *pgkit.DB, ctx context.Context, userID uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE unlimited_balances
		SET deleted_at = now()
		WHERE user_id = $1 AND deleted_at IS NULL
	`, userID)
	return err
}

func hasUnlimitedBalanceTx(tx pgx.Tx, ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM unlimited_balances
			WHERE user_id = $1 AND deleted_at IS NULL
		)
	`, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
