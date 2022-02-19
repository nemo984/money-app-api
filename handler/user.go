package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

// User
// swagger:response userResponse
type userResponse struct {
	Body db.User
}

// swagger:parameters loginUser createUser
type UsernamePasswordRequest struct {
	// The budget to create
	//
	// required: true
	// in: body
	User *usernamePasswordRequest `json:"user"`
}

// swagger:model
type usernamePasswordRequest struct {
	// required: true
	// example: Hello
	Username string `json:"username" binding:"required"`
	// required: true
	// min length: 6
	// example: Hello123
	Password string `json:"password" binding:"required,min=6"`
}

// Token
// swagger:response tokenResponse
type tokenResponse struct {
	Body struct {
		Token string `json:"token"`
	}
}

// swagger:route POST /token Users loginUser
// Returns an auth token for the user
// responses:
//  200: tokenResponse
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

// swagger:route POST /users Users createUser
// Returns the created user
// responses:
//  200: userResponse
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

	c.JSON(http.StatusCreated, user)
}

// swagger:route GET /me Users getUser
// Returns user of the auth token
// responses:
//  200: userResponse
func (s *Server) getUser(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	user, err := s.service.GetUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// swagger:parameters updateUser
type UpdateUserRequest struct {
	// The fields to update for the user
	//
	// required: true
	// in: body
	User *updateUserRequest `json:"user"`
}

// swagger:model
type updateUserRequest struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	ProfileURL string `json:"profile_url"`
}

// swagger:route PATCH /me Users updateUser
// Update user of the auth token with provided fields and returns the user
// responses:
//  200: userResponse
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

const (
	imagesFilePath = "images/user-profile-pics/"
	maxBodyBytes   = 2 * 1024 * 1024 // 2 Mb
)

var allowedFileExtensions = []string{".jpg", ".jpeg", ".png"}

// swagger:parameters updateUserPicture
type uploadReq struct {
	// Image file for profile picture < 2Mb
	// Required: true
	// In: formData
	// swagger:file
	File os.File `json:"file"`
}

// swagger:route PUT /me/picture Users updateUserPicture
// Update the profile picture of the user and returns the user
// Consumes:
//  - multipart/form-data
// responses:
//  200: userResponse
func (s *Server) uploadProfilePicture(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	file, err := c.FormFile("file")
	if err != nil {
		err := fmt.Errorf("get form err: %s", err.Error())
		log.Println(err)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fileExtension := filepath.Ext(file.Filename)
	var allowed bool
	for _, ext := range allowedFileExtensions {
		if fileExtension == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		err := errors.New("only .jpg, .jpeg, .png file extensions are allowed ")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newFilename := uuid.New().String() + fileExtension
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

// swagger:route DELETE /me Users deleteUser
// Delete user of the auth token
// responses:
//  204: noContent
func (s *Server) deleteUser(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	err := s.service.DeleteUser(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
