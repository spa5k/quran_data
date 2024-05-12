-- migrate:up
CREATE TABLE "edition" (
    id INTEGER not null primary key autoincrement,
    name TEXT not null,
    author TEXT,
    language TEXT not null,
    direction TEXT not null,
    source TEXT,
    type TEXT not null,
    enabled INTEGER not null
);

CREATE UNIQUE INDEX idx_name ON edition(name);

-- migrate:down
DROP table "edition";