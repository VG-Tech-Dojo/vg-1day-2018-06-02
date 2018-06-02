-- +migrate Up
CREATE TABLE room (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    birth TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    updated TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

-- +migrate Down
DROP TABLE room;
