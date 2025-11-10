-- +goose Up
ALTER TABLE users ALTER COLUMN country SET DEFAULT 'TN';


-- +goose Down
ALTER TABLE users ALTER COLUMN country DROP DEFAULT;
