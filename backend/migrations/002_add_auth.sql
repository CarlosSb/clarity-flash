-- Migration 002: add auth system to users table

-- Add password_hash column (nullable for backwards compatibility)
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;

-- Add updated_at column
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Add unique constraint on email
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'users_email_key'
  ) THEN
    ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
  END IF;
END $$;

-- Index for email lookups on login
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
