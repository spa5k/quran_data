-- migrate:up
CREATE TABLE "sajdah" (
    id INTEGER not null primary key autoincrement,
    sajdah_number INTEGER not null,
    surah_number INTEGER not null,
    ayah_number INTEGER not null,
    recommended INTEGER not null,
    obligatory INTEGER not null
);

-- migrate:down
DROP TABLE "sajdah";
