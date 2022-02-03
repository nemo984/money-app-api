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

-- name: UpdateUser :exec
UPDATE users 
SET name = $2,
    password = $3,
    profile_url = $4
WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;