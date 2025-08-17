CREATE TABLE IF NOT EXISTS url_shortener.kafka_link_events
(
  version Int32,
  kind    String,
  link_id UInt32,
  user_id UInt32,
  ts      String
)
ENGINE = Kafka
SETTINGS
  kafka_broker_list = 'kafka:9092',
  kafka_topic_list = 'link.events',
  kafka_group_name = 'ch_consumer',
  kafka_format = 'JSONEachRow',
  kafka_num_consumers = 1;

CREATE MATERIALIZED VIEW IF NOT EXISTS url_shortener.mv_link_visits
TO url_shortener.link_visits
AS
SELECT
  parseDateTimeBestEffort(ts) AS ts,
  link_id,
  user_id,
  kind
FROM url_shortener.kafka_link_events
WHERE kind = 'link.visited';
