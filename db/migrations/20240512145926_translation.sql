-- migrate:up
CREATE TABLE "translation" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "juz_number" INTEGER
);

-- migrate:down
DROP table "translation";
