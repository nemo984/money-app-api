package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetIncomeTypes(ctx context.Context) ([]db.IncomeType, error) {
	types, err := s.db.GetIncomeTypes(ctx)
	if err != nil {
		return []db.IncomeType{}, err
	}

	return types, nil
}

func (s *service) CreateIncomeType(ctx context.Context, name string) (db.IncomeType, error) {
	incomeType, err := s.db.CreateIncomeType(ctx, name)
	if err != nil {
		return db.IncomeType{}, err
	}

	return incomeType, nil
}

func (s *service) DeleteIncomeType(ctx context.Context, id int32) error {
	if err := s.db.DeleteIncomeType(ctx, id); err != nil {
		return err
	}

	return nil
}
