-- +goose Up 
ALTER table feeds
ADD COLUMN last_fetched_at TIMESTAMP;
-- +goose Down
ALTER table feeds DROP COLUMN last_fetched_at;