package models

import (
	"time"
	"github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WalletID  uuid.UUID `gorm:"type:uuid;not null"`
	Wallet    Wallet    `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amount   float64    `json:"amount"`   
	Type      string    `gorm:"type:varchar(10);"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
