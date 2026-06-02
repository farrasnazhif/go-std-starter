package service

import (
	"context"
	"time"

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

	// Create user with is_active = false (default)
	if err := s.store.Users.CreateAndInvite(ctx, user, "", time.Hour); err != nil {
		return nil, err
	}

	// Send OTP for verification
	if err := s.otpService.Send(ctx, email); err != nil {
		_ = s.store.Users.Delete(ctx, user.ID)
		return nil, err
	}

	return &RegisterResult{User: user}, nil
}

func (s *AuthService) Activate(ctx context.Context, token string) error {
	return s.store.Users.Activate(ctx, token)
}
