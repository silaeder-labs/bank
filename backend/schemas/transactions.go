package schemas

import "github.com/google/uuid"

type CreateTransactionRequest struct {
	TargetID uuid.UUID `json:"target_id" validate:"required,uuid4"`
	Amount   int64  `json:"amount" validate:"required,gt=0"`
	Comment  string `json:"comment,omitempty" validate:"max=100"`
}

type TransactionFull struct {
	ID        string `json:"transaction_id"`
	CreatedAt string `json:"created_at"`
	Amount    int64  `json:"amount"`
	Source    string  `json:"source" validate:"required,uuid4"`
	Target    string  `json:"target" validate:"required,uuid4"`
	Comment   string  `json:"comment,omitempty"`
}

type GetTransactionsRequest struct {
	Page    int `query:"page" validate:"gte=1"`
	Size	int `query:"size" validate:"gte=1,lte=100"`
}