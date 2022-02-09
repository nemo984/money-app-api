package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error) {
	if _, err := s.db.GetUser(ctx, args.Username); err == nil {
		return db.User{}, AppError{
			StatusCode: http.StatusConflict,
			Err:        errors.New("username already taken"),
		}
	}

	log.Printf("Creating User %#v\n", args)
	user, err := s.db.CreateUser(ctx, args)
	if err != nil {
		log.Println("[error] db create user error:", err.Error())
		return db.User{}, errors.New("cannot create user")
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, username string) error {
	if _, err := s.db.GetUser(ctx, username); err != nil {
		log.Println("[error] db get user error:", err.Error())
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("user with that username doesn't exists"),
			}
		}
		return err
	}

	err := s.db.DeleteUser(ctx, username)
	if err != nil {
		log.Println("[error] db delete user error:", err.Error())
		return errors.New("cannot delete user")
	}
	return nil
}
