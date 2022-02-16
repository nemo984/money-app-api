package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nemo984/money-app-api/service"
)

const (
	authorizationHeader  = "Authorization"
	authorizationBearer  = "bearer"
	authorizationPayload = "payload"
)

func authenticatedToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is missing")
			handleAbortError(c, err)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid auth header format")
			handleAbortError(c, err)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationBearer {
			err := fmt.Errorf("unsupported auth type, must be %s", authorizationBearer)
			handleAbortError(c, err)
			return
		}

		claims, err := service.VerifyToken(fields[1])
		if err != nil {
			handleAbortError(c, err)
			return
		}

		c.Set(authorizationPayload, claims)
		c.Next()
	}
}

func handleAbortError(c *gin.Context, err error) {
	switch v := err.(type) {
	case service.AppError:
		c.AbortWithStatusJSON(v.StatusCode, errorResponse(v.Err))
	case error:
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(v))
	}
}