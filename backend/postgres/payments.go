package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nrf24l01/go-web-utils/pgkit"
	"github.com/silaeder-labs/bank/backend/schemas"
)

type Payment struct {
	ID         uuid.UUID
	InsertedAt time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time

	From        uuid.UUID
	To          uuid.UUID
	Amount      int64
	Status      string
	Description string
}

func (p *Payment) ToPaymentFull() schemas.PaymentFull {
	return schemas.PaymentFull{
		ID:          p.ID.String(),
		CreateAt:    p.InsertedAt.Format(time.RFC3339),
		From:        p.From.String(),
		To:          p.To.String(),
		Amount:      p.Amount,
		Status:      schemas.Status(p.Status),
		Description: p.Description,
	}
}

func (p *Payment) Insert(db *pgkit.DB, ctx context.Context) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO payments (from_id, to_id, amount, description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, inserted_at, updated_at
	`, p.From, p.To, p.Amount, p.Description, p.Status).Scan(&p.ID, &p.InsertedAt, &p.UpdatedAt)
	return err
}

func GetPaymentByID(db *pgkit.DB, ctx context.Context, paymentID uuid.UUID, userID uuid.UUID) (*Payment, error) {
	var payment Payment
	err := db.Pool.QueryRow(ctx, `
		SELECT id, from_id, to_id, amount, description, status, inserted_at, updated_at
		FROM payments
		WHERE id = $1 AND deleted_at IS NULL AND (from_id = $2 OR to_id = $2)
	`, paymentID, userID).Scan(&payment.ID, &payment.From, &payment.To, &payment.Amount, &payment.Description, &payment.Status, &payment.InsertedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
