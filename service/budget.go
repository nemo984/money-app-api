package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetBudgets(ctx context.Context, userID int32) ([]db.Budget, error) {
	budgets, err := s.db.GetBudgets(ctx, userID)
	if err != nil {
		return []db.Budget{}, err
	}

	return budgets, nil
}

func (s *service) CreateBudget(ctx context.Context, args db.CreateBudgetParams) (db.Budget, error) {
	budget, err := s.db.CreateBudget(ctx, args)
	if err != nil {
		return db.Budget{}, err
	}

	return budget, nil
}

func (s *service) DeleteBudget(ctx context.Context, userID, budgetID int32) error {
	budget, err := s.db.GetBudget(ctx, budgetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("no budget with that id"),
			}
		}
		return fmt.Errorf("db get budget error: %v", err)
	}

	if budget.UserID != userID {
		return AppError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("not your budget"),
		}
	}

	return s.db.DeleteBudget(ctx, budgetID)
}
