package service

import (
	"context"
	"log"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error) {
	log.Printf("Creating User %#v\n", args)
	user, err := s.db.CreateUser(ctx, args)
	if err != nil {
		log.Println("[error] db create user error")
		//do some custom error? status code?
		return db.User{}, err
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, username string) error {
	err := s.db.DeleteUser(ctx, username)
	if err != nil {
		log.Println("[error] db delete user error")
		return err
	}
	return nil
}
