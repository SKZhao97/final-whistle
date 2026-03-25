-- Create player_ratings table
CREATE TABLE player_ratings (
  id BIGSERIAL PRIMARY KEY,
  check_in_id BIGINT NOT NULL,
  player_id BIGINT NOT NULL,
  rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 10),
  note VARCHAR(80),
  CONSTRAINT fk_player_ratings_check_in
    FOREIGN KEY (check_in_id)
    REFERENCES check_ins (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_player_ratings_player
    FOREIGN KEY (player_id)
    REFERENCES players (id)
    ON DELETE CASCADE
);

-- Create index on player_id for player lookups (required by spec)
CREATE INDEX player_ratings_player_id_idx ON player_ratings (player_id);

-- Create index on check_in_id for check-in lookups
CREATE INDEX player_ratings_check_in_id_idx ON player_ratings (check_in_id);

-- Add comment
COMMENT ON TABLE player_ratings IS 'Individual player ratings within a check-in';