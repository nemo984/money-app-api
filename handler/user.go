package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

//TODO: validator
type usernamePasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserResponse struct {
	UserID int32 `json:"id"`
}

type updateUserRequest struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	ProfileURL string `json:"profile_url"`
}

func (s *Server) createUserToken(c *gin.Context) {
	var req usernamePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	token, err := s.service.LoginUser(c, req.Username, req.Password)
	if err != nil {
		handleError(c, err)
		return
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

func (s *Server) getUser(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	user, err := s.service.GetUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) updateUser(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("%#v\n", req)

	user, err := s.service.UpdateUser(c, service.UpdateUserParams{
		UserID:     userPayload.UserID,
		Username:   req.Username,
		Name:       req.Name,
		Password:   req.Password,
		ProfileUrl: req.ProfileURL,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) deleteUser(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	err := s.service.DeleteUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
