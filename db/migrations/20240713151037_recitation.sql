-- migrate:up
CREATE TABLE IF NOT EXISTS recitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reciter_id INTEGER NOT NULL,
    surah_number INTEGER NOT NULL,
    recitation_data TEXT NOT NULL,
    FOREIGN KEY (reciter_id) REFERENCES reciters(id),
    UNIQUE (reciter_id, surah_number)
);

-- migrate:down

DROP TABLE IF EXISTS recitations;