-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, name, password, profile_url
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET username = COALESCE($2, username),
    name = COALESCE($3, name),
    password = COALESCE($4, password),
    profile_url = COALESCE($5, profile_url)
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;