-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPosts :many
SELECT 
    title
FROM posts
ORDER BY 
    created_at DESC
LIMIT
    $1;
