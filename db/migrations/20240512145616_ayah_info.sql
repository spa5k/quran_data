-- migrate:up
CREATE TABLE "ayah_info" (
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

CREATE UNIQUE INDEX "ayah_info_surah_number_ayah_number_unique" ON "ayah_info" ("surah_number", "ayah_number");

-- migrate:down
DROP TABLE "ayah_info";