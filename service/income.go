package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetIncomes(ctx context.Context, userID int32) ([]db.Income, error) {
	incomes, err := s.db.GetIncomes(ctx, userID)
	if err != nil {
		return []db.Income{}, err
	}

	return incomes, nil
}

func (s *service) CreateIncome(ctx context.Context, args db.CreateIncomeParams) (db.Income, error) {
	income, err := s.db.CreateIncome(ctx, args)
	if err != nil {
		return db.Income{}, err
	}

	return income, nil
}

func (s *service) DeleteIncome(ctx context.Context, userID, incomeID int32) error {
	income, err := s.db.GetIncome(ctx, incomeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("no income with that id"),
			}
		}
		return fmt.Errorf("db get income error: %v", err)
	}

	if income.UserID != userID {
		return AppError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("not your expense"),
		}
	}

	return s.db.DeleteIncome(ctx, incomeID)
}
