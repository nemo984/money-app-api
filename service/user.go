package service

import (
	"context"
	"errors"
	"log"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error) {
	if _, err := s.db.GetUser(ctx, args.Username); err == nil {
		//TODO: custom error with conflict status code
		return db.User{}, errors.New("user with that username already exists")
	}

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
	if _, err := s.db.GetUser(ctx, username); err != nil {
		//TODO: custom error with 404 code
		return errors.New("user with that username doesn't exists")
	}

	err := s.db.DeleteUser(ctx, username)
	if err != nil {
		log.Println("[error] db delete user error")
		return err
	}
	return nil
}
