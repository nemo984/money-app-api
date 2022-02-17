package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	user, err := s.service.GetUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) updateUser(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.service.UpdateUser(c, service.UpdateUserParams{
		UserID:   userPayload.UserID,
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

const imagesFilePath = "images/user-profile-pics/"

func (s *Server) uploadProfilePicture(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	file, err := c.FormFile("file")
	if err != nil {
		err := fmt.Errorf("get form err: %s", err.Error())
		log.Println(err)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//TODO: validate extension, resize image?
	newFilename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := imagesFilePath + newFilename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		err := fmt.Errorf("upload file err: %s", err.Error())
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.service.UpdateUserPicture(c, db.UpdateUserPictureParams{
		UserID: userPayload.UserID,
		ProfileUrl: sql.NullString{
			String: filepath,
			Valid:  true,
		},
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) deleteUser(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	err := s.service.DeleteUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
