package repositories

import (
	"context"
	"database/sql"
	"time"
)

type OTPRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(ctx context.Context, email, code, purpose string, expiry time.Duration) error {
	query := `INSERT INTO otp_codes (email, code, purpose, expiry) VALUES ($1, $2, $3, $4)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, email, code, purpose, time.Now().Add(expiry))
	return err
}

func (r *OTPRepository) Verify(ctx context.Context, email, code, purpose string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM otp_codes WHERE email = $1 AND code = $2 AND purpose = $3 AND expiry > $4)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, code, purpose, time.Now()).Scan(&exists)
	return exists, err
}

func (r *OTPRepository) DeleteByEmail(ctx context.Context, email, purpose string) error {
	query := `DELETE FROM otp_codes WHERE email = $1 AND purpose = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, email, purpose)
	return err
}
