CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE IF NOT EXISTS "ayah" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "tajweed" TEXT,
    "juz_number" INTEGER
);
CREATE TABLE IF NOT EXISTS "ayah_info" (
    id INTEGER not null primary key autoincrement,
    surah_number INTEGER not null,
    ayah_number INTEGER not null,
    ayah_key TEXT not null,
    hizb INTEGER not null,
    rub_el_hizb INTEGER not null,
    ruku INTEGER not null,
    manzil INTEGER not null,
    page INTEGER not null,
    juz INTEGER not null
);
CREATE TABLE IF NOT EXISTS "edition" (
    id INTEGER not null primary key autoincrement,
    name TEXT not null,
    author TEXT,
    language TEXT not null,
    direction TEXT not null,
    source TEXT,
    type TEXT not null,
    enabled INTEGER not null
);
CREATE UNIQUE INDEX idx_name ON edition(name);
CREATE TABLE IF NOT EXISTS "juz" (
    id INTEGER not null primary key autoincrement,
    juz_number INTEGER not null,
    start_surah INTEGER not null,
    start_ayah INTEGER not null,
    end_surah INTEGER not null,
    end_ayah INTEGER not null
);
CREATE TABLE IF NOT EXISTS "sajdah" (
    id INTEGER not null primary key autoincrement,
    sajdah_number INTEGER not null,
    surah_number INTEGER not null,
    ayah_number INTEGER not null,
    recommended INTEGER not null,
    obligatory INTEGER not null
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
    "page_end" INTEGER NOT NULL,
    "juz_number" INTEGER
);
CREATE TABLE IF NOT EXISTS "translation" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "surah_number" INTEGER NOT NULL,
    "ayah_number" INTEGER NOT NULL,
    "edition_id" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "juz_number" INTEGER
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240512144958'),
  ('20240512145616'),
  ('20240512145707'),
  ('20240512145749'),
  ('20240512145817'),
  ('20240512145841'),
  ('20240512145926');
