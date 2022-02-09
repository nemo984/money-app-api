package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Service interface {
	CreateUser(context.Context, db.CreateUserParams) (db.User, error)
	LoginUser(ctx context.Context, username string, password string) (int32, error)
	DeleteUser(context.Context, int32) error
	UpdateUser(context.Context, db.UpdateUserParams) error
}

type service struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &service{db: db}
}

type AppError struct {
	StatusCode int
	Err        error
}

func (a AppError) Error() string {
	return a.Err.Error()
}
