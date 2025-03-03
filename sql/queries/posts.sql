-- name: CreatePost :one
INSERT INTO posts (
        id,
        created_at,
        updated_at,
        name,
        -- ✅ Ensure "name" is included here
        title,
        description,
        published_at,
        url,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id,
    created_at,
    updated_at,
    name,
    title,
    description,
    published_at,
    url,
    feed_id;
-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
    JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2;