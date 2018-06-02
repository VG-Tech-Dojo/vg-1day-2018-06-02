-- +migrate Up
ALTER TABLE message ADD COLUMN room_id INTEGER;
