package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
)

var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

type JWTClaims struct {
	jwt.MapClaims
	UserID int32 `json:"user_id"`
}

func CreateToken(userID int32) (string, error) {
	claims := JWTClaims{
		UserID: userID,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", AppError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return token, nil
}

func (s *service) VerifyToken(ctx context.Context, token string) (JWTClaims, error) {
	payload, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return JWTClaims{}, AppError{
			StatusCode: http.StatusUnauthorized,
			Err:        err,
		}
	}

	claims, ok := payload.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, AppError{
			StatusCode: http.StatusUnauthorized,
			Err:        errors.New("jwt claims mismatch"),
		}
	}

	_, err = s.GetUser(ctx, claims.UserID)
	if errors.Is(err, ErrUserNotFound) {
		return JWTClaims{}, ErrUserNotFound
	}
	if err != nil {
		return JWTClaims{}, AppError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("verify token get user error: %w", err),
		}
	}

	return *claims, nil
}
