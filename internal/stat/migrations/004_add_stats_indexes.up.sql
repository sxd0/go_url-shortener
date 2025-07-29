ALTER TABLE stats
    ADD CONSTRAINT IF NOT EXISTS uniq_link_date UNIQUE (link_id, date);
CREATE INDEX IF NOT EXISTS idx_stats_user_date  ON stats (user_id, date);
CREATE INDEX IF NOT EXISTS idx_stats_link_date  ON stats (link_id, date);