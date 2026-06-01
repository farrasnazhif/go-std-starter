package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"github.com/farrasnazhif/go-std-starter/internal/store/models"
	"github.com/google/uuid"
)

type AuthService struct {
	store       store.Storage
	mailer      mailer.Client
	frontendURL string
	env         string
	mailExp     time.Duration
}

func NewAuthService(store store.Storage, mailer mailer.Client, frontendURL, env string, mailExp time.Duration) *AuthService {
	return &AuthService{
		store:       store,
		mailer:      mailer,
		frontendURL: frontendURL,
		env:         env,
		mailExp:     mailExp,
	}
}

type RegisterResult struct {
	User  *models.User
	Token string
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*RegisterResult, error) {
	user := &models.User{
		Username: username,
		Email:    email,
	}

	if err := user.Password.Set(password); err != nil {
		return nil, err
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := s.store.Users.CreateAndInvite(ctx, user, hashToken, s.mailExp); err != nil {
		return nil, err
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", s.frontendURL, plainToken)
	isProd := s.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	if err := s.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProd); err != nil {
		// rollback user creation if email fails (SAGA pattern)
		_ = s.store.Users.Delete(ctx, user.ID)
		return nil, err
	}

	return &RegisterResult{User: user, Token: plainToken}, nil
}

func (s *AuthService) Activate(ctx context.Context, token string) error {
	return s.store.Users.Activate(ctx, token)
}
