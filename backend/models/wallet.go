package models

import (
    "github.com/google/uuid"
    "time"
)

type Wallet struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID    uuid.UUID
    User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
    Address   string    `gorm:"not null"`
    Balance   float64   `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
