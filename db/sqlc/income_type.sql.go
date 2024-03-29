// Code generated by sqlc. DO NOT EDIT.
// source: income_type.sql

package db

import (
	"context"
)

const createIncomeType = `-- name: CreateIncomeType :one
INSERT INTO income_types (
  name
) VALUES (
  $1
)
RETURNING income_type_id, name
`

func (q *Queries) CreateIncomeType(ctx context.Context, name string) (IncomeType, error) {
	row := q.db.QueryRowContext(ctx, createIncomeType, name)
	var i IncomeType
	err := row.Scan(&i.IncomeTypeID, &i.Name)
	return i, err
}

const deleteIncomeType = `-- name: DeleteIncomeType :exec
DELETE FROM income_types
WHERE income_type_id = $1
`

func (q *Queries) DeleteIncomeType(ctx context.Context, incomeTypeID int32) error {
	_, err := q.db.ExecContext(ctx, deleteIncomeType, incomeTypeID)
	return err
}

const getIncomeTypes = `-- name: GetIncomeTypes :many
SELECT income_type_id, name FROM income_types
`

func (q *Queries) GetIncomeTypes(ctx context.Context) ([]IncomeType, error) {
	rows, err := q.db.QueryContext(ctx, getIncomeTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []IncomeType{}
	for rows.Next() {
		var i IncomeType
		if err := rows.Scan(&i.IncomeTypeID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
