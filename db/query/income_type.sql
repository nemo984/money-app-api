-- name: GetIncomeTypes :many
SELECT * FROM income_types;

-- name: CreateIncomeType :one
INSERT INTO income_types (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteIncomeType :exec
DELETE FROM income_types
WHERE income_type_id = $1;