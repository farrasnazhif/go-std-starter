package service

import (
	"context"

	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"github.com/farrasnazhif/go-std-starter/internal/store/models"
)

type AuthService struct {
	store      store.Storage
	mailer     mailer.Client
	env        string
	otpService *OTPService
}

func NewAuthService(store store.Storage, mailer mailer.Client, env string) *AuthService {
	return &AuthService{
		store:      store,
		mailer:     mailer,
		env:        env,
		otpService: NewOTPService(store, mailer, env),
	}
}

type RegisterResult struct {
	User *models.User
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*RegisterResult, error) {
	user := &models.User{
		Username: username,
		Email:    email,
	}

	if err := user.Password.Set(password); err != nil {
		return nil, err
	}

	if err := s.store.Users.Create(ctx, user); err != nil {
		return nil, err
	}

	if err := s.otpService.Send(ctx, email, PurposeActivation); err != nil {
		_ = s.store.Users.Delete(ctx, user.ID)
		return nil, err
	}

	return &RegisterResult{User: user}, nil
}
