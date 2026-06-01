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
	CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error
	Activate(ctx context.Context, token string) error
	Delete(ctx context.Context, userID int64) error
}

type Storage struct {
	Users Users
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: repositories.NewUserRepository(db),
	}
}
