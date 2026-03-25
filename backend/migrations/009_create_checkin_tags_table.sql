-- Create checkin_tags table
CREATE TABLE checkin_tags (
  id BIGSERIAL PRIMARY KEY,
  check_in_id BIGINT NOT NULL,
  tag_id BIGINT NOT NULL,
  CONSTRAINT fk_checkin_tags_check_in
    FOREIGN KEY (check_in_id)
    REFERENCES check_ins (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_checkin_tags_tag
    FOREIGN KEY (tag_id)
    REFERENCES tags (id)
    ON DELETE CASCADE
);

-- Create unique index to prevent duplicate tag assignments
CREATE UNIQUE INDEX checkin_tags_check_in_tag_idx ON checkin_tags (check_in_id, tag_id);

-- Create index on check_in_id for check-in lookups
CREATE INDEX checkin_tags_check_in_id_idx ON checkin_tags (check_in_id);

-- Create index on tag_id for tag lookups
CREATE INDEX checkin_tags_tag_id_idx ON checkin_tags (tag_id);

-- Add comment
COMMENT ON TABLE checkin_tags IS 'Many-to-many relationship between check-ins and tags';