-- Create matches table
CREATE TABLE matches (
  id BIGSERIAL PRIMARY KEY,
  competition VARCHAR(100) NOT NULL,
  season VARCHAR(50) NOT NULL,
  round VARCHAR(50),
  status VARCHAR(20) NOT NULL CHECK (status IN ('SCHEDULED', 'FINISHED')),
  kickoff_at TIMESTAMPTZ NOT NULL,
  home_team_id BIGINT NOT NULL,
  away_team_id BIGINT NOT NULL,
  home_score INTEGER,
  away_score INTEGER,
  venue VARCHAR(200),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_matches_home_team
    FOREIGN KEY (home_team_id)
    REFERENCES teams (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_matches_away_team
    FOREIGN KEY (away_team_id)
    REFERENCES teams (id)
    ON DELETE CASCADE
);

-- Create index on kickoff_at for time-based queries
CREATE INDEX matches_kickoff_at_idx ON matches (kickoff_at);

-- Create index on status for filtering
CREATE INDEX matches_status_idx ON matches (status);

-- Create index on competition and season for filtering
CREATE INDEX matches_competition_season_idx ON matches (competition, season);

-- Create index on home_team_id and away_team_id for team lookups
CREATE INDEX matches_home_team_id_idx ON matches (home_team_id);
CREATE INDEX matches_away_team_id_idx ON matches (away_team_id);

-- Add comment
COMMENT ON TABLE matches IS 'Football matches between teams';