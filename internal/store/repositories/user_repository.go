package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/farrasnazhif/go-std-starter/internal/store/models"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource conflicted")
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
	QueryTimeoutDuration = time.Second * 5
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, email, username, created_at FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, username, is_active, created_at FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ActivateByEmail(ctx context.Context, email string) error {
	query := `UPDATE users SET is_active = true WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

func (r *UserRepository) ResetPassword(ctx context.Context, email string, hashedPassword []byte) error {
	query := `UPDATE users SET password = $1 WHERE email = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, hashedPassword, email)
	return err
}

func (r *UserRepository) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		if err := r.Create(ctx, tx, user); err != nil {
			return err
		}
		return r.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
	})
}

func (r *UserRepository) Activate(ctx context.Context, token string) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		user, err := r.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		if err := r.update(ctx, tx, user); err != nil {
			return err
		}
		return r.deleteUserInvitations(ctx, tx, user.ID)
	})
}

func (r *UserRepository) Delete(ctx context.Context, userID int64) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		if err := r.delete(ctx, tx, userID); err != nil {
			return err
		}
		return r.deleteUserInvitations(ctx, tx, userID)
	})
}

func (r *UserRepository) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &models.User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	return err
}

func (r *UserRepository) update(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	return err
}

func (r *UserRepository) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *UserRepository) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
