package service

import (
	"errors"
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

func VerifyToken(token string) (JWTClaims, error) {
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
	return *claims, nil
}

