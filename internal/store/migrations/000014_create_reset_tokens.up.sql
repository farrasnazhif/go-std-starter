CREATE TABLE IF NOT EXISTS reset_tokens (
  token varchar(64) PRIMARY KEY,
  email citext NOT NULL,
  expiry timestamp(0) with time zone NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
