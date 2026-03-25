-- Create sessions table
CREATE TABLE sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  token VARCHAR(255) NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_sessions_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

-- Create unique index on token for lookups
CREATE UNIQUE INDEX sessions_token_idx ON sessions (token);

-- Create index on user_id for user lookups
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

-- Create index on expired_at for cleanup queries
CREATE INDEX sessions_expired_at_idx ON sessions (expired_at);

-- Add comment
COMMENT ON TABLE sessions IS 'User sessions for authentication (Cookie-based)';