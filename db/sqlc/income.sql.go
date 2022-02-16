// Code generated by sqlc. DO NOT EDIT.
// source: income.sql

package db

import (
	"context"
	"database/sql"
)

const createIncome = `-- name: CreateIncome :one
INSERT INTO incomes (
  income_type_id, description, amount, frequency, user_id
) VALUES (
  $2, $3, $4, $5, $1
)
RETURNING income_id, income_type_id, description, amount, created_at, frequency, user_id
`

type CreateIncomeParams struct {
	UserID       int32          `json:"user_id"`
	IncomeTypeID int32          `json:"income_type_id"`
	Description  sql.NullString `json:"description"`
	Amount       string         `json:"amount"`
	Frequency    DateFrequency  `json:"frequency"`
}

func (q *Queries) CreateIncome(ctx context.Context, arg CreateIncomeParams) (Income, error) {
	row := q.db.QueryRowContext(ctx, createIncome,
		arg.UserID,
		arg.IncomeTypeID,
		arg.Description,
		arg.Amount,
		arg.Frequency,
	)
	var i Income
	err := row.Scan(
		&i.IncomeID,
		&i.IncomeTypeID,
		&i.Description,
		&i.Amount,
		&i.CreatedAt,
		&i.Frequency,
		&i.UserID,
	)
	return i, err
}

const deleteIncome = `-- name: DeleteIncome :exec
DELETE FROM incomes
WHERE income_id = $1
`

func (q *Queries) DeleteIncome(ctx context.Context, incomeID int32) error {
	_, err := q.db.ExecContext(ctx, deleteIncome, incomeID)
	return err
}

const getIncome = `-- name: GetIncome :one
SELECT income_id, income_type_id, description, amount, created_at, frequency, user_id FROM incomes
WHERE income_id = $1
`

func (q *Queries) GetIncome(ctx context.Context, incomeID int32) (Income, error) {
	row := q.db.QueryRowContext(ctx, getIncome, incomeID)
	var i Income
	err := row.Scan(
		&i.IncomeID,
		&i.IncomeTypeID,
		&i.Description,
		&i.Amount,
		&i.CreatedAt,
		&i.Frequency,
		&i.UserID,
	)
	return i, err
}

const getIncomes = `-- name: GetIncomes :many
SELECT income_id, income_type_id, description, amount, created_at, frequency, user_id FROM incomes
WHERE user_id = $1
`

func (q *Queries) GetIncomes(ctx context.Context, userID int32) ([]Income, error) {
	rows, err := q.db.QueryContext(ctx, getIncomes, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Income
	for rows.Next() {
		var i Income
		if err := rows.Scan(
			&i.IncomeID,
			&i.IncomeTypeID,
			&i.Description,
			&i.Amount,
			&i.CreatedAt,
			&i.Frequency,
			&i.UserID,
		); err != nil {
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
