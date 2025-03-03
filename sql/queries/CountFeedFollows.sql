-- name: CountFeedFollows :one
SELECT COUNT(*)
FROM feed_follows
WHERE feed_id = $1;