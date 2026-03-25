-- Create players table
CREATE TABLE players (
  id BIGSERIAL PRIMARY KEY,
  team_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  position VARCHAR(50),
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_players_team
    FOREIGN KEY (team_id)
    REFERENCES teams (id)
    ON DELETE CASCADE
);

-- Create unique index on slug
CREATE UNIQUE INDEX players_slug_idx ON players (slug);

-- Create index on team_id for lookups
CREATE INDEX players_team_id_idx ON players (team_id);

-- Create index on name for search
CREATE INDEX players_name_idx ON players (name);

-- Add comment
COMMENT ON TABLE players IS 'Football players belonging to teams';