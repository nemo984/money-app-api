package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetCategories(ctx context.Context) ([]db.Category, error) {
	categories, err := s.db.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

func (s *service) CreateCategory(ctx context.Context, name string) (db.Category, error) {
	category, err := s.db.CreateCategory(ctx, name)
	if err != nil {
		return db.Category{}, err
	}

	return category, nil
}

func (s *service) DeleteCategory(ctx context.Context, id int32) error {
	if err := s.db.DeleteCategory(ctx, id); err != nil {
		return err
	}
	
	return nil
}