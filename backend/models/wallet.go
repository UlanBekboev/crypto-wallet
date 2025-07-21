package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique"` // 1-1 связь
	Balance   int64     `gorm:"default:0"` // В центах (например, 100 = 1.00)
	CreatedAt time.Time
	UpdatedAt time.Time
}
