-- Create teams table
CREATE TABLE teams (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  short_name VARCHAR(50),
  slug VARCHAR(100) NOT NULL,
  logo_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create unique index on slug
CREATE UNIQUE INDEX teams_slug_idx ON teams (slug);

-- Create index on name for search
CREATE INDEX teams_name_idx ON teams (name);

-- Add comment
COMMENT ON TABLE teams IS 'Football teams that participate in matches';