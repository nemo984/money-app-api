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
SET 
    username = CASE WHEN @username_do_update::boolean
        THEN @username::VARCHAR(20) ELSE username END,

    name = CASE WHEN @name_do_update::boolean
        THEN @name::VARCHAR(20) ELSE name END,

    password = CASE WHEN @password_do_update::boolean
        THEN @password::VARCHAR ELSE password END,

    profile_url = CASE WHEN @profile_do_update::boolean
        THEN @profile_url::VARCHAR ELSE profile_url END

WHERE user_id = @user_id
RETURNING *;

-- name: UpdateUserPicture :one
UPDATE users
SET profile_url = $2
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;