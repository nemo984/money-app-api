-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, name, password, profile_url
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET username = $2,
    name = $3,
    password = $4,
    profile_url = $5
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;