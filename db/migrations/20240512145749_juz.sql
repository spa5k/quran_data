-- migrate:up
CREATE TABLE IF NOT EXISTS "juz" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "juz_number" INTEGER NOT NULL,
    "start_surah" INTEGER NOT NULL,
    "start_ayah" INTEGER NOT NULL,
    "end_surah" INTEGER NOT NULL,
    "end_ayah" INTEGER NOT NULL
);
CREATE UNIQUE INDEX "juz_juz_number_unique" ON "juz" ("juz_number");

-- migrate:down
DROP TABLE "juz";