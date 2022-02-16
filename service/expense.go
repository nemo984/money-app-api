package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetExpenses(ctx context.Context, userID int32) ([]db.Expense, error) {
	expenses, err := s.db.GetExpenses(ctx, userID)
	if err != nil {
		return []db.Expense{}, err
	}

	return expenses, nil
}

func (s *service) CreateExpense(ctx context.Context, args db.CreateExpenseParams) (db.Expense, error) {
	expense, err := s.db.CreateExpense(ctx, args)
	if err != nil {
		return db.Expense{}, err
	}

	return expense, nil
}

func (s *service) DeleteExpense(ctx context.Context, userID, expenseID int32) error {
	expense, err := s.db.GetExpense(ctx, expenseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("no expense with that id"),
			}
		}
		return fmt.Errorf("db get expense error: %v", err)
	}

	if expense.UserID != userID {
		return AppError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("not your expense"),
		}
	}

	return s.db.DeleteExpense(ctx, expenseID)
}
