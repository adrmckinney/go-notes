package models

type SignUpRequest struct {
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type SignInRequest struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}
