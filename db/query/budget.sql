-- name: GetBudgets :many
SELECT * FROM budgets
WHERE user_id = $1;

-- name: GetBudget :one
SELECT * FROM budgets
WHERE budget_id = $1;

-- name: CreateBudget :one
INSERT INTO budgets (
  category_id, amount, end_date, user_id
) VALUES (
  $2, $3, $4, $1
)
RETURNING *;

-- name: UpdateBudget :exec
UPDATE budgets
SET category_id = $2,
    amount = $3,
    end_date = $4
WHERE user_id = $1;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE budget_id = $1;