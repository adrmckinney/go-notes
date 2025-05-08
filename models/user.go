package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint            `json:"id,omitempty" gorm:"primaryKey"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Username  string          `json:"username" gorm:"unique"`
	Password  string          `json:"password"`
	CreatedAt *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type UserWithToken struct {
	User
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	FirstName       *string `json:"firstName,omitempty"`
	LastName        *string `json:"lastName,omitempty"`
	Username        *string `json:"username,omitempty"`
	Password        *string `json:"password,omitempty"`
	ConfirmPassword *string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

var AllowedUserUpdateFields = map[string]bool{
	"firstName": true,
	"lastName":  true,
	"username":  true,
	"password":  true,
}
