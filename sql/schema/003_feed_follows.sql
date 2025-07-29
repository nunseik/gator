-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_feed
    FOREIGN KEY (feed_id)
    REFERENCES feed(id) ON DELETE CASCADE,
    CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uc_feed_user UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;