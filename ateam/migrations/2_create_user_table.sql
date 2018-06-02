-- +migrate Up
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT DEFAULT "",
    point INTEGER NOT NULL DEFAULT 0,
    created TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    updated TIMESTAMP NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

-- +migrate Down
DROP TABLE users;
