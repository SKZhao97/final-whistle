-- Create users table
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on email for lookups
CREATE UNIQUE INDEX users_email_idx ON users (email);

-- Create index on name for search
CREATE INDEX users_name_idx ON users (name);

-- Add comment
COMMENT ON TABLE users IS 'Registered users of the application';