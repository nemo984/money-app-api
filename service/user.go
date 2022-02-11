package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetUser(ctx context.Context, userID int32) (db.User, error) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("user doesn't exist"),
			}
		}
		return db.User{}, err
	}
	return user, nil
}

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
		if err == sql.ErrNoRows {
			return 0, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("cannot find user with that username"),
			}
		}
		log.Println("[error] db get user error:", err.Error())
		return 0, err
	}
	//comparing hashed ones
	if user.Password != password {
		return 0, AppError{
			StatusCode: http.StatusUnauthorized,
			Err:        errors.New("incorrect password"),
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

type UpdateUserParams struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	ProfileUrl string `json:"profile_url"`
	UserID     int32  `json:"user_id"`
}

func (s *service) UpdateUser(ctx context.Context, args UpdateUserParams) (db.User, error) {
	var (
		updateUsername   bool
		updateName       bool
		updatePassword   bool
		updateProfileURL bool
	)

	if args.Username != "" {
		updateUsername = true
	}
	if args.Name != "" {
		updateName = true
	}
	if args.Password != "" {
		updatePassword = true
	}
	if args.ProfileUrl != "" {
		updateProfileURL = true
	}

	updateArgs := db.UpdateUserParams{
		UserID:           args.UserID,
		UsernameDoUpdate: updateUsername,
		Username:         args.Username,
		NameDoUpdate:     updateName,
		Name:             args.Name,
		PasswordDoUpdate: updatePassword,
		Password:         args.Password,
		ProfileDoUpdate:  updateProfileURL,
		ProfileUrl:       args.ProfileUrl,
	}

	user, err := s.db.UpdateUser(ctx, updateArgs)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("user doesn't exists"),
			}
		}
		log.Println("[error] db update user error:", err.Error())
		return db.User{}, err
	}
	return user, nil
}
