CREATE DATABASE IF NOT EXISTS url_shortener;

CREATE TABLE IF NOT EXISTS url_shortener.link_visits
(
  ts       DateTime,
  link_id  UInt32,
  user_id  UInt32,
  kind     LowCardinality(String),
  uuid     UUID DEFAULT generateUUIDv4()
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(ts)
ORDER BY (link_id, ts)
TTL ts + INTERVAL 30 DAY DELETE
SETTINGS index_granularity = 8192;
