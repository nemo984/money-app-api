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
	JWTTokenCookieKey = "jwt-token"

	authorizationHeader  = "Authorization"
	authorizationBearer  = "bearer"
	AuthorizationPayload = "payload"
)

func (s *Server) authenticatedToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtToken string
		//check cookie first
		if cookie, err := c.Request.Cookie(JWTTokenCookieKey); err == nil {
			jwtToken = cookie.Value
		} else {
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

			jwtToken = fields[1]
		}

		claims, err := s.service.VerifyToken(c, jwtToken)
		if err != nil {
			handleAbortError(c, err)
			return
		}

		c.Set(AuthorizationPayload, claims)
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
