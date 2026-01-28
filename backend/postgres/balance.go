package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nrf24l01/go-web-utils/pgkit"
	"github.com/silaeder-labs/bank/backend/schemas"
)

type Balance struct {
	UserID     uuid.UUID
	InsertedAt time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time

	AmountCents int64
}

func (b *Balance) ToBalanceFull() schemas.BalanceFull {
	return schemas.BalanceFull{
		Id:      b.UserID.String(),
		Balance: b.AmountCents,
	}
}

func GetBalanceByUserID(db *pgkit.DB, ctx context.Context, userID uuid.UUID) (*Balance, error) {
	var balance Balance
	err := db.Pool.QueryRow(ctx, "SELECT amount_cents FROM balances WHERE user_id = $1 AND deleted_at IS NULL", userID).Scan(&balance.AmountCents)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func UpdateBalance(db *pgkit.DB, ctx context.Context, userID uuid.UUID, amountCents int64) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO balances (user_id, amount_cents)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET amount_cents = balances.amount_cents + EXCLUDED.amount_cents
	`, userID, amountCents)
	return err
}

func CheckUserCanPay(db *pgkit.DB, ctx context.Context, userID uuid.UUID, amountCents int64) (bool, error) {
	var currentBalance int64
	err := db.Pool.QueryRow(ctx, "SELECT amount_cents FROM balances WHERE user_id = $1 AND deleted_at IS NULL", userID).Scan(&currentBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return currentBalance >= amountCents, nil
}
