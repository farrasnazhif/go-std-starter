package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"github.com/farrasnazhif/go-std-starter/internal/store/models"
)

const (
	otpExpiry        = 5 * time.Minute
	resetTokenExpiry = 10 * time.Minute
	PurposeActivation    = "activation"
	PurposePasswordReset = "password_reset"
)

var (
	ErrInvalidOTP        = errors.New("invalid or expired OTP code")
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	ErrUserActive        = errors.New("user is already active")
)

type OTPService struct {
	store  store.Storage
	mailer mailer.Client
	env    string
}

func NewOTPService(store store.Storage, mailer mailer.Client, env string) *OTPService {
	return &OTPService{store: store, mailer: mailer, env: env}
}

func (s *OTPService) Send(ctx context.Context, email, purpose string) error {
	_ = s.store.OTPs.DeleteByEmail(ctx, email, purpose)

	code, err := generateOTP()
	if err != nil {
		return err
	}

	if err := s.store.OTPs.Create(ctx, email, code, purpose, otpExpiry); err != nil {
		return err
	}

	vars := struct{ Code string }{Code: code}
	isProd := s.env == "production"
	return s.mailer.Send(mailer.OTPTemplate, "", email, vars, !isProd)
}

func (s *OTPService) Verify(ctx context.Context, email, code, purpose string) error {
	valid, err := s.store.OTPs.Verify(ctx, email, code, purpose)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidOTP
	}

	if purpose == PurposeActivation {
		if err := s.store.Users.ActivateByEmail(ctx, email); err != nil {
			return err
		}
	}

	return s.store.OTPs.DeleteByEmail(ctx, email, purpose)
}

func (s *OTPService) VerifyForReset(ctx context.Context, email, code string) (string, error) {
	valid, err := s.store.OTPs.Verify(ctx, email, code, PurposePasswordReset)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", ErrInvalidOTP
	}

	_ = s.store.OTPs.DeleteByEmail(ctx, email, PurposePasswordReset)

	token, err := generateResetToken()
	if err != nil {
		return "", err
	}

	if err := s.store.ResetTokens.Create(ctx, token, email, resetTokenExpiry); err != nil {
		return "", err
	}

	return token, nil
}

func (s *OTPService) ResetPassword(ctx context.Context, token, newPassword string) error {
	email, err := s.store.ResetTokens.Validate(ctx, token)
	if err != nil {
		return ErrInvalidResetToken
	}

	var p models.Password
	if err := p.Set(newPassword); err != nil {
		return err
	}

	if err := s.store.Users.ResetPassword(ctx, email, p.Hash); err != nil {
		return err
	}

	return s.store.ResetTokens.Delete(ctx, token)
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
