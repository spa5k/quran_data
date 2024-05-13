-- migrate:up
CREATE TABLE IF NOT EXISTS "surah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "name_simple" TEXT NOT NULL,
    "name_complex" TEXT NOT NULL,
    "name_arabic" TEXT NOT NULL,
    "ayah_start" INTEGER NOT NULL,
    "ayah_end" INTEGER NOT NULL,
    "revelation_place" TEXT NOT NULL,
    "page_start" INTEGER NOT NULL,
    "page_end" INTEGER NOT NULL
);
CREATE UNIQUE INDEX "surah_surah_number_unique" ON "surah" ("surah_number");

-- migrate:down
DROP table "surah";
