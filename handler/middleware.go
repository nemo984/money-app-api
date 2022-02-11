package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader  = "Authorization"
	authorizationBearer  = "bearer"
	authorizationPayload = "payload"
)

//TODO: token package?
type JWTClaims struct {
	jwt.MapClaims
	UserID int32 `json:"user_id"`
}

func authenticatedToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid auth header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported auth type, must be %s", authorizationBearer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		token := fields[1]
		payload, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		claims := payload.Claims.(*JWTClaims)
		c.Set(authorizationPayload, claims)
		c.Next()
	}
}
