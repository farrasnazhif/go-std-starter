package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/store"
)

var (
	ErrInvalidOTP = errors.New("invalid or expired OTP code")
	ErrUserActive = errors.New("user is already active")
)

const otpExpiry = 5 * time.Minute

type OTPService struct {
	store  store.Storage
	mailer mailer.Client
	env    string
}

func NewOTPService(store store.Storage, mailer mailer.Client, env string) *OTPService {
	return &OTPService{store: store, mailer: mailer, env: env}
}

func (s *OTPService) Send(ctx context.Context, email string) error {
	// Delete old OTPs for this email
	_ = s.store.OTPs.DeleteByEmail(ctx, email)

	code, err := generateOTP()
	if err != nil {
		return err
	}

	if err := s.store.OTPs.Create(ctx, email, code, otpExpiry); err != nil {
		return err
	}

	vars := struct{ Code string }{Code: code}
	isProd := s.env == "production"
	return s.mailer.Send(mailer.OTPTemplate, "", email, vars, !isProd)
}

func (s *OTPService) Verify(ctx context.Context, email, code string) error {
	valid, err := s.store.OTPs.Verify(ctx, email, code)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidOTP
	}

	if err := s.store.Users.ActivateByEmail(ctx, email); err != nil {
		return err
	}

	return s.store.OTPs.DeleteByEmail(ctx, email)
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
