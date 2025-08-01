-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed
    FOREIGN KEY (feed_id)
    REFERENCES feed(id) ON DELETE CASCADE,
    CONSTRAINT uc_feed_url UNIQUE (feed_id, url)
);

-- +goose Down
DROP TABLE posts;