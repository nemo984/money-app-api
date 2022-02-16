-- name: GetIncomes :many
SELECT * FROM incomes
WHERE user_id = $1;

-- name: GetIncome :one
SELECT * FROM incomes
WHERE income_id = $1;

-- name: CreateIncome :one
INSERT INTO incomes (
  income_type_id, description, amount, frequency, user_id
) VALUES (
  $2, $3, $4, $5, $1
)
RETURNING *;

-- name: DeleteIncome :exec
DELETE FROM incomes
WHERE income_id = $1;