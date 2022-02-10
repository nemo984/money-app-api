package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
)

const SECRET_KEY = "TODO: use_env_later"

//TODO: validator
type usernamePasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserResponse struct {
	UserID int32 `json:"id"`
}

type updateUserRequest struct {
	Username *string `json:"username"`
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	ProfileUrl *string `json:"profile_url"`
}

type idURI struct {
	UserID int32 `uri:"id" binding:"required"`
}

func (s *Server) createUserToken(c *gin.Context) {
	var req usernamePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := s.service.LoginUser(c, req.Username, req.Password)
	if err != nil {	
		handleError(c, err)
		return
	}

	claims := JWTClaims{
		UserID: userID,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(SECRET_KEY))
	if err != nil {
		handleError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) createUser(c *gin.Context) {
	var req usernamePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//TODO: hash the pwd inside service
	user, err := s.service.CreateUser(c, db.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createUserResponse{
		UserID: user.UserID,
	})
}

func (s *Server) updateUser(c *gin.Context) {
	var uri idURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := s.service.UpdateUser(c, db.UpdateUserParams{
		UserID: uri.UserID,
		Username: *req.Username,
		// Name: sql.NullString{

		// },
		Password: *req.Password,
		// ProfileUrl: *req.ProfileUrl,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) deleteUser(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(*JWTClaims)
	err := s.service.DeleteUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
