package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nrf24l01/go-web-utils/pgkit"
	"github.com/silaeder-labs/bank/backend/schemas"
)

type Transaction struct {
	LineID      uuid.UUID
	InsertedAt  time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time

	From        uuid.UUID
	To          uuid.UUID
	AmountCents int64
	Description string
}

func (t *Transaction) Insert(db *pgkit.DB, ctx context.Context) error {
	if err := db.Pool.QueryRow(ctx, "INSERT INTO transactions (from_user_id, to_user_id, amount_cents, description) VALUES ($1, $2, $3, $4) RETURNING line_id, inserted_at, updated_at, deleted_at",
		t.From, t.To, t.AmountCents, t.Description).Scan(&t.LineID, &t.InsertedAt, &t.UpdatedAt, &t.DeletedAt); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) ToTransactionFull() schemas.TransactionFull {
	return schemas.TransactionFull{
		ID:      t.LineID.String(),
		CreatedAt:  t.InsertedAt.Format(time.RFC3339),
		Source:        t.From.String(),
		Target:          t.To.String(),
		Amount:      t.AmountCents,
		Comment: t.Description,
	}
}