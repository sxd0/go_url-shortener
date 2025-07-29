ALTER TABLE stats DROP CONSTRAINT IF EXISTS uniq_link_date;
DROP INDEX IF EXISTS idx_stats_user_date;
DROP INDEX IF EXISTS idx_stats_link_date;