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
	if err := db.Pool.QueryRow(ctx, "INSERT INTO transactions (from_user_id, to_user_id, amount_cents, description) VALUES ($1, $2, $3, $4) RETURNING line_id, inserted_at, updated_at",
		t.From, t.To, t.AmountCents, t.Description).Scan(&t.LineID, &t.InsertedAt, &t.UpdatedAt); err != nil {
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

func GetTransactionsByUserID(db *pgkit.DB, ctx context.Context, userID uuid.UUID, limit, offset int) ([]Transaction, error) {
	rows, err := db.Pool.Query(ctx, "SELECT line_id, inserted_at, from_user_id, to_user_id, amount_cents, description FROM transactions WHERE (from_user_id = $1 OR to_user_id = $1) AND (deleted_at IS NULL) ORDER BY inserted_at DESC LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.LineID, &t.InsertedAt, &t.From, &t.To, &t.AmountCents, &t.Description); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func GetTransactionByID(db *pgkit.DB, ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) (*Transaction, error) {
	var t Transaction
	if err := db.Pool.QueryRow(ctx, "SELECT line_id, inserted_at, from_user_id, to_user_id, amount_cents, description FROM transactions WHERE line_id = $1 AND deleted_at IS NULL AND (from_user_id = $2 OR to_user_id = $2)", transactionID, userID).
		Scan(&t.LineID, &t.InsertedAt, &t.From, &t.To, &t.AmountCents, &t.Description); err != nil {
		return nil, err
	}
	return &t, nil
}