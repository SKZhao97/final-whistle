-- Create check_ins table
CREATE TABLE check_ins (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  match_id BIGINT NOT NULL,
  watched_type VARCHAR(20) NOT NULL CHECK (watched_type IN ('FULL', 'PARTIAL', 'HIGHLIGHTS')),
  supporter_side VARCHAR(20) NOT NULL CHECK (supporter_side IN ('HOME', 'AWAY', 'NEUTRAL')),
  match_rating INTEGER NOT NULL CHECK (match_rating BETWEEN 1 AND 10),
  home_team_rating INTEGER NOT NULL CHECK (home_team_rating BETWEEN 1 AND 10),
  away_team_rating INTEGER NOT NULL CHECK (away_team_rating BETWEEN 1 AND 10),
  short_review VARCHAR(280),
  watched_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_check_ins_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_check_ins_match
    FOREIGN KEY (match_id)
    REFERENCES matches (id)
    ON DELETE CASCADE
);

-- Create unique index for one check-in per user per match
CREATE UNIQUE INDEX check_ins_user_match_idx ON check_ins (user_id, match_id);

-- Create index on created_at for chronological access
CREATE INDEX check_ins_created_at_idx ON check_ins (created_at);

-- Create index on match_id for match lookups
CREATE INDEX check_ins_match_id_idx ON check_ins (match_id);

-- Create index on user_id for user lookups
CREATE INDEX check_ins_user_id_idx ON check_ins (user_id);

-- Add comment
COMMENT ON TABLE check_ins IS 'User check-ins for watched matches with ratings and reviews';