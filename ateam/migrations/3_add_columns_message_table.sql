-- +migrate Up
ALTER TABLE message ADD userid INTEGER;
ALTER TABLE message ADD messagetype INTEGER;

-- +migrate Down
DROP TABLE message;
