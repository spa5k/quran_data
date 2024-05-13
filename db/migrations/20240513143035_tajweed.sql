-- migrate:up
CREATE TABLE IF NOT EXISTS "tajweed" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "tajweed" TEXT NOT NULL,
    UNIQUE ("surah_number", "ayah_number")
);
-- migrate:down

DROP  TABLE "tajweed";