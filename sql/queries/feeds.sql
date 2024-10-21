-- name: CreateFeeds :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds 
WHERE url = $1 LIMIT 1;


-- name: GetNextFeedToFetch :one
SELECT * FROM feeds 
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = Now(),
updated_at = Now()
WHERE id = $1
RETURNING *;