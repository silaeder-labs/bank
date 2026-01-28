package schemas

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	FromID      uuid.UUID `json:"from_id" validate:"required,uuid4"`
	ToID        uuid.UUID `json:"to_id" validate:"required,uuid4"`
	Amount      int64     `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description,omitempty" validate:"max=120"`
}

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "UNPAID"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusCancelled PaymentStatus = "CANCELLED"
)

type PaymentFull struct {
	ID          string        `json:"id"`
	CreateAt    string        `json:"created_at"`
	From        string        `json:"from" validate:"required,uuid4"`
	To          string        `json:"to" validate:"required,uuid4"`
	Amount      int64         `json:"amount"`
	Status      PaymentStatus `json:"status"`
	Description string        `json:"description,omitempty"`
}
