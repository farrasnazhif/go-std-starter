package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ResetTokenRepository struct {
	db *sql.DB
}

func NewResetTokenRepository(db *sql.DB) *ResetTokenRepository {
	return &ResetTokenRepository{db: db}
}

func (r *ResetTokenRepository) Create(ctx context.Context, token, email string, expiry time.Duration) error {
	query := `INSERT INTO reset_tokens (token, email, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, token, email, time.Now().Add(expiry))
	return err
}

func (r *ResetTokenRepository) Validate(ctx context.Context, token string) (string, error) {
	query := `SELECT email FROM reset_tokens WHERE token = $1 AND expiry > $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var email string
	err := r.db.QueryRowContext(ctx, query, token, time.Now()).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return email, nil
}

func (r *ResetTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM reset_tokens WHERE token = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
