-- name: GetCategories :many
SELECT * FROM categories;

-- name: CreateCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = $1;