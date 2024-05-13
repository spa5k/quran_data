-- migrate:up
CREATE TABLE IF NOT EXISTS "ayah_info" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "ayah_key" TEXT NOT NULL,
    "hizb" INTEGER NOT NULL,
    "rub_el_hizb" INTEGER NOT NULL,
    "ruku" INTEGER NOT NULL,
    "manzil" INTEGER NOT NULL,
    "page" INTEGER NOT NULL,
    "juz" INTEGER NOT NULL,
    CHECK ("ayah_number" > 0 AND "surah_number" > 0)
);
CREATE UNIQUE INDEX "ayah_info_surah_number_ayah_number_unique" ON "ayah_info" ("surah_number", "ayah_number");

-- migrate:down
DROP TABLE "ayah_info";