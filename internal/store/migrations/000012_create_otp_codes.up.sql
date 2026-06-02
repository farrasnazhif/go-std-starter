CREATE TABLE IF NOT EXISTS otp_codes (
  id bigserial PRIMARY KEY,
  email citext NOT NULL,
  code varchar(6) NOT NULL,
  expiry timestamp(0) with time zone NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
