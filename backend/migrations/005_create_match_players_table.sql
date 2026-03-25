-- Create match_players table
CREATE TABLE match_players (
  id BIGSERIAL PRIMARY KEY,
  match_id BIGINT NOT NULL,
  player_id BIGINT NOT NULL,
  team_id BIGINT NOT NULL,
  CONSTRAINT fk_match_players_match
    FOREIGN KEY (match_id)
    REFERENCES matches (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_match_players_player
    FOREIGN KEY (player_id)
    REFERENCES players (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_match_players_team
    FOREIGN KEY (team_id)
    REFERENCES teams (id)
    ON DELETE CASCADE
);

-- Create unique index to prevent duplicate player entries in same match
CREATE UNIQUE INDEX match_players_match_player_idx ON match_players (match_id, player_id);

-- Create index on match_id for match lookups
CREATE INDEX match_players_match_id_idx ON match_players (match_id);

-- Create index on player_id for player lookups
CREATE INDEX match_players_player_id_idx ON match_players (player_id);

-- Create index on team_id for team lookups
CREATE INDEX match_players_team_id_idx ON match_players (team_id);

-- Add comment
COMMENT ON TABLE match_players IS 'Players participating in specific matches';