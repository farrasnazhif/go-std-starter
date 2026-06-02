package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/store/models"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

type Users interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ActivateByEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email string, hashedPassword []byte) error
	Delete(ctx context.Context, userID int64) error
}

type OTPs interface {
	Create(ctx context.Context, email, code, purpose string, expiry time.Duration) error
	Verify(ctx context.Context, email, code, purpose string) (bool, error)
	DeleteByEmail(ctx context.Context, email, purpose string) error
}

type ResetTokens interface {
	Create(ctx context.Context, token, email string, expiry time.Duration) error
	Validate(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}

type Storage struct {
	Users       Users
	OTPs        OTPs
	ResetTokens ResetTokens
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:       repositories.NewUserRepository(db),
		OTPs:        repositories.NewOTPRepository(db),
		ResetTokens: repositories.NewResetTokenRepository(db),
	}
}
