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

func (r *OTPRepository) Create(ctx context.Context, email, code string, expiry time.Duration) error {
	query := `INSERT INTO otp_codes (email, code, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, email, code, time.Now().Add(expiry))
	return err
}

func (r *OTPRepository) Verify(ctx context.Context, email, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM otp_codes WHERE email = $1 AND code = $2 AND expiry > $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, code, time.Now()).Scan(&exists)
	return exists, err
}

func (r *OTPRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM otp_codes WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}
