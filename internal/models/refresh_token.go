package models

import "time"

type RefreshToken struct {
	ID         uint       `gorm:"primaryKey"`
	UserID     uint       `gorm:"index;not null"`
	TokenID    string     `gorm:"uniqueIndex;size:64;not null"` // jti
	TokenHash  string     `gorm:"size:255;not null"`
	ExpiresAt  time.Time  `gorm:"index;not null"`
	RevokedAt  *time.Time `gorm:"index"`
	CreatedAt  time.Time
}
