CREATE TABLE stats (
    id SERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL,
    user_id INTEGER,
    "date" DATE NOT NULL,
    clicks INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT uniq_stat_per_day UNIQUE(link_id, "date")
);
