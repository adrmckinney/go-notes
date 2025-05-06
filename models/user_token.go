package models

import "time"

type UserToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserId    uint      `gorm:"not null"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
