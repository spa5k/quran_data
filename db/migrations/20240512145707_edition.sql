-- migrate:up
CREATE TABLE IF NOT EXISTS "edition" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL,
    "author" TEXT,
    "language" TEXT NOT NULL,
    "direction" TEXT NOT NULL,
    "source" TEXT,
    "type" TEXT NOT NULL,
    "enabled" INTEGER NOT NULL CHECK ("enabled" IN (0, 1))
);
CREATE UNIQUE INDEX "idx_name" ON "edition" ("name");

-- migrate:down
DROP table "edition";