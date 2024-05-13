-- migrate:up
CREATE TABLE "translation" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "juz_number" INTEGER
);

CREATE UNIQUE INDEX "translation_surah_number_ayah_number_edition_id_unique" ON "translation" ("surah_number", "ayah_number", "edition_id");

-- migrate:down
DROP table "translation";
