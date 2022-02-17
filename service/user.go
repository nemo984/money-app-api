package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = AppError{
		StatusCode: http.StatusNotFound,
		Err:        errors.New("user doesn't exists"),
	}
	ErrUserAlreadyExists = AppError{
		StatusCode: http.StatusConflict,
		Err:        errors.New("user already exists"),
	}
)

func (s *service) GetUser(ctx context.Context, userID int32) (db.User, error) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, AppError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return user, nil
}

func (s *service) GetUserByName(ctx context.Context, username string) (db.User, error) {
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, AppError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error) {
	u, err := s.GetUserByName(ctx, args.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return db.User{}, err
	}
	if err == nil {
		return u, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), 4)
	if err != nil {
		log.Println("[error] bcrypt hashing pwd:", err.Error())
		return db.User{}, err
	}

	args.Password = string(hashedPassword)
	log.Printf("Creating User %#v\n", args)
	user, err := s.db.CreateUser(ctx, args)
	if err != nil {
		log.Println("[error] db create user error:", err.Error())
		return db.User{}, errors.New("cannot create user")
	}

	return user, nil
}

func (s *service) LoginUser(ctx context.Context, username, password string) (token string, err error) {
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return "", err
	}
	//comparing hashed ones
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", AppError{
			StatusCode: http.StatusUnauthorized,
			Err:        errors.New("incorrect password"),
		}
	}

	token, err = CreateToken(user.UserID)
	if err != nil {
		return "", err
	}
	return token, nil
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
	updateArgs := db.UpdateUserParams{
		UserID:           args.UserID,
		UsernameDoUpdate: args.Username != "",
		Username:         args.Username,
		NameDoUpdate:     args.Name != "",
		Name:             args.Name,
		PasswordDoUpdate: args.Password != "",
		Password:         args.Password,
		ProfileDoUpdate:  args.ProfileUrl != "",
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

func (s *service) UpdateUserPicture(ctx context.Context, args db.UpdateUserPictureParams) (db.User, error) {
	user, err := s.db.UpdateUserPicture(ctx, args)
	if err != nil {
		return db.User{}, nil
	}

	return user, nil
}
