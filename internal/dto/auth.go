package dto

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=225"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
