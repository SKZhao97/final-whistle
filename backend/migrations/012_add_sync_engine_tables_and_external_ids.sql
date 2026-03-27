ALTER TABLE teams
  ADD COLUMN IF NOT EXISTS external_source VARCHAR(50),
  ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
  ADD COLUMN IF NOT EXISTS external_updated_at TIMESTAMPTZ;

ALTER TABLE players
  ADD COLUMN IF NOT EXISTS external_source VARCHAR(50),
  ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
  ADD COLUMN IF NOT EXISTS external_updated_at TIMESTAMPTZ;

ALTER TABLE matches
  ADD COLUMN IF NOT EXISTS external_source VARCHAR(50),
  ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
  ADD COLUMN IF NOT EXISTS external_updated_at TIMESTAMPTZ;

CREATE UNIQUE INDEX IF NOT EXISTS teams_external_unique_idx
ON teams (external_source, external_id);

CREATE UNIQUE INDEX IF NOT EXISTS players_external_unique_idx
ON players (external_source, external_id);

CREATE UNIQUE INDEX IF NOT EXISTS matches_external_unique_idx
ON matches (external_source, external_id);

CREATE TABLE IF NOT EXISTS sync_jobs (
  id BIGSERIAL PRIMARY KEY,
  job_type VARCHAR(50) NOT NULL,
  scope_type VARCHAR(50) NOT NULL,
  scope_key VARCHAR(200) NOT NULL,
  dedupe_key VARCHAR(255) NOT NULL,
  trigger_mode VARCHAR(20) NOT NULL CHECK (trigger_mode IN ('automatic', 'manual')),
  priority INT NOT NULL DEFAULT 100,
  status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'running', 'succeeded', 'failed', 'canceled')),
  scheduled_at TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  attempt INT NOT NULL DEFAULT 0,
  max_attempts INT NOT NULL DEFAULT 3,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS sync_jobs_dedupe_active_idx
ON sync_jobs (dedupe_key)
WHERE status IN ('pending', 'running');

CREATE INDEX IF NOT EXISTS sync_jobs_status_scheduled_idx
ON sync_jobs (status, scheduled_at, priority, id);

CREATE TABLE IF NOT EXISTS sync_cursors (
  id BIGSERIAL PRIMARY KEY,
  provider VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  scope_key VARCHAR(200) NOT NULL,
  last_success_at TIMESTAMPTZ,
  last_attempt_at TIMESTAMPTZ,
  last_error_at TIMESTAMPTZ,
  last_error TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (provider, resource_type, scope_key)
);
