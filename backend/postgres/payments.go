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
	Creator     uuid.UUID
	Amount      int64
	Status      schemas.PaymentStatus
	Description string
}

func (p *Payment) ToPaymentFull() schemas.PaymentFull {
	return schemas.PaymentFull{
		ID:          p.ID.String(),
		CreateAt:    p.InsertedAt.Format(time.RFC3339),
		From:        p.From.String(),
		To:          p.To.String(),
		Amount:      p.Amount,
		Status:      p.Status,
		Description: p.Description,
	}
}

func (p *Payment) Insert(db *pgkit.DB, ctx context.Context) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO payments (from_id, to_id, amount, description, status, creator_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, inserted_at, updated_at
	`, p.From, p.To, p.Amount, p.Description, p.Status, p.Creator).Scan(&p.ID, &p.InsertedAt, &p.UpdatedAt)
	return err
}

func GetPaymentByID(db *pgkit.DB, ctx context.Context, paymentID uuid.UUID, userID uuid.UUID) (*Payment, error) {
	var payment Payment
	err := db.Pool.QueryRow(ctx, `
		SELECT id, from_id, to_id, amount, description, status, creator_id, inserted_at, updated_at
		FROM payments
		WHERE id = $1 AND deleted_at IS NULL AND (from_id = $2 OR to_id = $2 OR creator_id = $2) AND status='UNPAID'
	`, paymentID, userID).Scan(&payment.ID, &payment.From, &payment.To, &payment.Amount, &payment.Description, &payment.Status, &payment.Creator, &payment.InsertedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *Payment) ChangeStatus(db *pgkit.DB, ctx context.Context, newStatus schemas.PaymentStatus) error {
	err := db.Pool.QueryRow(ctx, "UPDATE payments SET status=$1 WHERE id=$2 RETURNING updated_at, status", newStatus, p.ID).Scan(&p.UpdatedAt, &p.Status)
	return err
}
