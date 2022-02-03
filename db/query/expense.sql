-- name: GetExpenses :many
SELECT * FROM expenses
WHERE user_id = $1;

-- name: CreateExpense :one
INSERT INTO expenses (
  category_id, amount, frequency, note
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE user_id = $1;