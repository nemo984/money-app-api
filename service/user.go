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

func (s *service) LoginUser(ctx context.Context, username, password string) (int32, error) {
	//TODO: function for this
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		log.Println("[error] db get user error:", err.Error())
		if err == sql.ErrNoRows {
			return 0, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("cannot find user with that username"),
			}
		}
		return 0, err
	}
	//comparing hashed ones
	if user.Password != password {
		return 0, AppError{
			StatusCode: http.StatusUnauthorized,
			Err: errors.New("incorrect password"),
		}
	} 
	
	return user.UserID, nil
}

func (s *service) DeleteUser(ctx context.Context, userID int32) error {
	if _, err := s.db.GetUserByID(ctx, userID); err != nil {
		log.Println("[error] db get user error:", err.Error())
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("cannot find user with that id"),
			}
		}
		return err
	}

	err := s.db.DeleteUser(ctx, userID)
	if err != nil {
		log.Println("[error] db delete user error:", err.Error())
		return errors.New("cannot delete user")
	}
	return nil
}

func (s *service) UpdateUser(ctx context.Context, args db.UpdateUserParams) error {
	if _, err := s.db.UpdateUser(ctx, args); err != nil {
		if err == sql.ErrNoRows {
			return AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("user doesn't exists"),
			}
		}
		return err
	}
	return nil
}
