package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Username  string         `json:"username" gorm:"unique"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

var AllowedUserCreateFields = map[string]bool{
	"firstName": true,
	"lastName":  true,
	"username":  true,
	"password":  true,
}

var AllowedUserUpdateFields = map[string]bool{
	"firstName": true,
	"lastName":  true,
	"username":  true,
	"password":  true,
}
