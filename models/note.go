package models

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID        uint            `json:"id,omitempty" gorm:"primaryKey"`
	UserID    uint            `json:"userId" gorm:"not null"`
	Title     string          `json:"title,omitempty"`
	Content   string          `json:"content,omitempty"`
	CreatedAt *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

var AllowedNoteUpdateFields = map[string]bool{
	"title":   true,
	"content": true,
}
