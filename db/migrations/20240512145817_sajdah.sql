-- migrate:up
CREATE TABLE IF NOT EXISTS "sajdah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "sajdah_number" INTEGER NOT NULL,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "recommended" INTEGER NOT NULL CHECK ("recommended" IN (0, 1)),
    "obligatory" INTEGER NOT NULL CHECK ("obligatory" IN (0, 1))
);

-- migrate:down
DROP TABLE "sajdah";
