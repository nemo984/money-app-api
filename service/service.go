package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Service interface {
	CreateUser(context.Context, db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type service struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &service{db: db}
}
