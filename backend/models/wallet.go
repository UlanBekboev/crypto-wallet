package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique"` // 1-1 связь
	Balance   float64   `gorm:"type:numeric(10,4);default:0.0000"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
