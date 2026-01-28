package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nrf24l01/go-web-utils/pgkit"
	"github.com/silaeder-labs/bank/backend/schemas"
)

type Transaction struct {
	LineID     uuid.UUID
	InsertedAt time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time

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
		ID:        t.LineID.String(),
		CreatedAt: t.InsertedAt.Format(time.RFC3339),
		Source:    t.From.String(),
		Target:    t.To.String(),
		Amount:    t.AmountCents,
		Comment:   t.Description,
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

func MakeTransaction(db *pgkit.DB, ctx context.Context, from uuid.UUID, to uuid.UUID, amount int64, description string) (*Transaction, error) {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	balances, err := getBalancesForUpdate(tx, ctx, from, to)
	if err != nil {
		return nil, err
	}

	fromBalance := balances[from]
	if fromBalance < amount {
		return nil, ErrCantPay
	}

	transaction := Transaction{
		From:        from,
		To:          to,
		AmountCents: amount,
		Description: description,
	}

	if err := insertTransactionTx(tx, ctx, &transaction); err != nil {
		return nil, err
	}
	if err := updateBalancesTx(tx, ctx, from, to, amount); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	committed = true

	return &transaction, nil
}

func getBalancesForUpdate(tx pgx.Tx, ctx context.Context, from uuid.UUID, to uuid.UUID) (map[uuid.UUID]int64, error) {
	balances := map[uuid.UUID]int64{
		from: 0,
		to:   0,
	}
	ids := []uuid.UUID{from}
	if from != to {
		ids = append(ids, to)
	}

	rows, err := tx.Query(ctx, "SELECT user_id, amount_cents FROM balances WHERE user_id = ANY($1) AND deleted_at IS NULL ORDER BY user_id FOR UPDATE", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID uuid.UUID
		var amountCents int64
		if err := rows.Scan(&userID, &amountCents); err != nil {
			return nil, err
		}
		balances[userID] = amountCents
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}

func updateBalancesTx(tx pgx.Tx, ctx context.Context, from uuid.UUID, to uuid.UUID, amountCents int64) error {
	if from == to {
		return nil
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO balances (user_id, amount_cents)
		VALUES ($1, $2), ($3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET amount_cents = balances.amount_cents + EXCLUDED.amount_cents
	`, from, -amountCents, to, amountCents)
	return err
}

func insertTransactionTx(tx pgx.Tx, ctx context.Context, t *Transaction) error {
	return tx.QueryRow(ctx, "INSERT INTO transactions (from_user_id, to_user_id, amount_cents, description) VALUES ($1, $2, $3, $4) RETURNING line_id, inserted_at, updated_at",
		t.From, t.To, t.AmountCents, t.Description).Scan(&t.LineID, &t.InsertedAt, &t.UpdatedAt)
}