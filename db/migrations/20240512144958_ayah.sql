-- migrate:up
CREATE TABLE "ayah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "tajweed" TEXT,
    "juz_number" INTEGER
);

-- migrate:down
DROP TABLE "ayah";