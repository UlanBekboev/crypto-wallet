package models

import (
	"time"
)

type PasswordReset struct {
	Email     string    `gorm:"primaryKey"`
	Token     string
	ExpiresAt time.Time
}
