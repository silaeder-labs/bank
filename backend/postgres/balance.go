package postgres

import (
	"github.com/google/uuid"
)

type Balance struct {
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	AmountCents int64     `gorm:"not null"`
}