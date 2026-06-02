package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/store/models"
	"github.com/farrasnazhif/go-std-starter/internal/store/repositories"
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error
	Activate(ctx context.Context, token string) error
	ActivateByEmail(ctx context.Context, email string) error
	Delete(ctx context.Context, userID int64) error
}

type OTPs interface {
	Create(ctx context.Context, email, code string, expiry time.Duration) error
	Verify(ctx context.Context, email, code string) (bool, error)
	DeleteByEmail(ctx context.Context, email string) error
}

type Storage struct {
	Users Users
	OTPs  OTPs
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: repositories.NewUserRepository(db),
		OTPs:  repositories.NewOTPRepository(db),
	}
}
