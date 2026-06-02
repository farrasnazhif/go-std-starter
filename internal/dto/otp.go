package dto

type SendOTPPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTPPayload struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type ResetPasswordPayload struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=3,max=72"`
}
