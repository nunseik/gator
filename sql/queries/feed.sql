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

-- name: MarkFeedFetched :exec
UPDATE feed
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feed
WHERE last_fetched_at < NOW() - INTERVAL '1 hour'
   OR last_fetched_at IS NULL
-- Use NULLS FIRST to prioritize feeds that have never been fetched
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;