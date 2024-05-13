-- migrate:up
CREATE TABLE "juz" (
    id INTEGER not null primary key autoincrement,
    juz_number INTEGER not null,
    start_surah INTEGER not null,
    start_ayah INTEGER not null,
    end_surah INTEGER not null,
    end_ayah INTEGER not null
);

CREATE UNIQUE INDEX "juz_juz_number_unique" ON "juz" ("juz_number");

-- migrate:down
DROP TABLE "juz";