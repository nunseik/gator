-- name: CreateFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :many
SELECT name, url, user_id FROM feed;

-- name: GetFeedByURL :one
SELECT * FROM feed WHERE url = $1;

-- name: GetFeedById :one
SELECT * FROM feed WHERE id = $1;
