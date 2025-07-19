package models

import (
	"time"
	"github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WalletID  uuid.UUID `gorm:"type:uuid;not null"`
	Wallet    Wallet    `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amount    float64   `gorm:"not null"`
	Type      string    `gorm:"type:varchar(10);not null"` // например: "in" или "out"
	CreatedAt time.Time
	UpdatedAt time.Time
}
