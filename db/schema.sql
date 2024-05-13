CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
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
CREATE TABLE IF NOT EXISTS "juz" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "juz_number" INTEGER NOT NULL,
    "start_surah" INTEGER NOT NULL,
    "start_ayah" INTEGER NOT NULL,
    "end_surah" INTEGER NOT NULL,
    "end_ayah" INTEGER NOT NULL
);
CREATE UNIQUE INDEX "juz_juz_number_unique" ON "juz" ("juz_number");
CREATE TABLE IF NOT EXISTS "sajdah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "sajdah_number" INTEGER NOT NULL,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "recommended" INTEGER NOT NULL CHECK ("recommended" IN (0, 1)),
    "obligatory" INTEGER NOT NULL CHECK ("obligatory" IN (0, 1))
);
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
CREATE TABLE IF NOT EXISTS "translation" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "juz_number" INTEGER,
    FOREIGN KEY ("edition_id") REFERENCES "edition"("id"),
    FOREIGN KEY ("surah_number", "ayah_number") REFERENCES "ayah"("surah_number", "ayah_number")
);
CREATE UNIQUE INDEX "translation_surah_number_ayah_number_edition_id_unique" ON "translation" ("surah_number", "ayah_number", "edition_id");
CREATE TABLE IF NOT EXISTS "tajweed" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "tajweed" TEXT NOT NULL,
    UNIQUE ("surah_number", "ayah_number")
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240512144958'),
  ('20240512145616'),
  ('20240512145707'),
  ('20240512145749'),
  ('20240512145817'),
  ('20240512145841'),
  ('20240512145926'),
  ('20240513143035');
