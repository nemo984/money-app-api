-- name: GetNotifications :many
SELECT * FROM notifications
WHERE user_id = $1;

-- name: CreateNotification :one
INSERT INTO notifications (
  user_id, description, type, priority
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateNotification :exec
UPDATE notifications 
SET read = $2
WHERE user_id = $1;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE user_id = $1;