-- +migrate Up
ALTER TABLE message ADD (
    userid INTEGER,
    messagetype INTEGER
);

-- +migrate Down
DROP TABLE message;
