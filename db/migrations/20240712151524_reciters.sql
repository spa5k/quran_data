-- migrate:up
CREATE TABLE "reciters" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    reciter_name TEXT NOT NULL,
    style TEXT NOT NULL,
    slug TEXT NOT NULL,
    translated_name TEXT NOT NULL,
    language_name TEXT NOT NULL,
    source TEXT NOT NULL,
    source_id INT NOT NULL
);
-- migrate:down
DROP TABLE "reciters";