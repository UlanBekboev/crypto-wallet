package models

import (
	"time"
	"github.com/google/uuid"
)

type Transaction struct {
	ID           uuid.UUID `json:"id" db:"id"`
	FromWalletID uuid.UUID `json:"from_wallet_id" db:"from_wallet_id"`
	ToWalletID   uuid.UUID `json:"to_wallet_id" db:"to_wallet_id"`
	Amount       int64     `json:"amount" db:"amount"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
