-- migrate:up
CREATE TABLE IF NOT EXISTS "ayah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    FOREIGN KEY ("edition_id") REFERENCES "edition"("id"),
    CHECK ("ayah_number" > 0)
);
CREATE UNIQUE INDEX "ayah_surah_number_ayah_number_edition_id_unique" ON "ayah" ("surah_number", "ayah_number", "edition_id");

-- migrate:down
DROP TABLE "ayah";