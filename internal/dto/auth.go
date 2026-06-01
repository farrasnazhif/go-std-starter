package dto

import "github.com/farrasnazhif/go-std-starter/internal/store/models"

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=225"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*models.User
	Token string `json:"token"`
}
