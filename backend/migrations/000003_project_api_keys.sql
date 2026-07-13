ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS project_id UUID REFERENCES projects(id) ON DELETE CASCADE;
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS key_prefix TEXT;
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS api_keys_project_id_idx ON api_keys (project_id);
CREATE INDEX IF NOT EXISTS api_keys_active_project_id_idx ON api_keys (project_id) WHERE revoked_at IS NULL;
