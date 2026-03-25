-- Create tags table
CREATE TABLE tags (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  slug VARCHAR(50) NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create unique index on slug
CREATE UNIQUE INDEX tags_slug_idx ON tags (slug);

-- Create index on sort_order for ordering
CREATE INDEX tags_sort_order_idx ON tags (sort_order);

-- Create index on is_active for filtering
CREATE INDEX tags_is_active_idx ON tags (is_active);

-- Add comment
COMMENT ON TABLE tags IS 'Predefined tags for check-in emotions and impressions';