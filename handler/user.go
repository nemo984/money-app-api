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
	Body User
}

type User struct {
	UserID int32 `json:"user_id"`
	// example: apodqila
	Username string `json:"username"`
	// example: John Smith
	Name       sql.NullString `json:"name"`
	ProfileUrl sql.NullString `json:"profile_url"`
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
//  401: userLoginError
func (h *handler) createUserToken(c *gin.Context) (interface{}, int, error) {
	var req usernamePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := validationErrors(&req, err)
		return nil, 0, err
	}

	token, err := h.service.LoginUser(c, req.Username, req.Password)
	if err != nil {
		return nil, 0, err
	}

	return gin.H{
		"token": token,
	}, http.StatusOK, nil
}

// swagger:route POST /users Users createUser
// Returns the created user
// responses:
//  200: userResponse
//  409: usernameTakenError
func (h *handler) createUser(c *gin.Context) (interface{}, int, error) {
	var req usernamePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := validationErrors(&req, err)
		return nil, 0, err
	}

	user, err := h.service.CreateUser(c, db.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, 0, err
	}

	return User{
		UserID:     user.UserID,
		Username:   user.Username,
		Name:       user.Name,
		ProfileUrl: user.ProfileUrl,
	}, http.StatusCreated, nil
}

// swagger:route GET /me Users getUser
// Returns user of the auth token
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: userResponse
func (h *handler) getUser(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	user, err := h.service.GetUser(c, userPayload.UserID)
	if err != nil {
		return nil, 0, err
	}

	return User{
		UserID:     user.UserID,
		Username:   user.Username,
		Name:       user.Name,
		ProfileUrl: user.ProfileUrl,
	}, http.StatusOK, nil
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
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// swagger:route PATCH /me Users updateUser
// Update user of the auth token with provided fields and returns the user
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: userResponse
func (h *handler) updateUser(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := validationErrors(&req, err)
		return nil, 0, err
	}

	user, err := h.service.UpdateUser(c, service.UpdateUserParams{
		UserID:   userPayload.UserID,
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, 0, err
	}

	return User{
		UserID:     user.UserID,
		Username:   user.Username,
		Name:       user.Name,
		ProfileUrl: user.ProfileUrl,
	}, http.StatusOK, nil
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
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: userResponse
func (h *handler) uploadProfilePicture(c *gin.Context) (interface{}, int, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	file, err := c.FormFile("file")
	if err != nil {
		err := fmt.Errorf("get form err: %s", err.Error())
		log.Println(err)
		return nil, http.StatusBadRequest, err
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
		return nil, http.StatusBadRequest, err
	}

	newFilename := uuid.New().String() + fileExtension
	filepath := imagesFilePath + newFilename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		err := fmt.Errorf("upload file err: %s", err.Error())
		log.Println(err)
		return nil, http.StatusInternalServerError, err
	}

	user, err := h.service.UpdateUserPicture(c, db.UpdateUserPictureParams{
		UserID: userPayload.UserID,
		ProfileUrl: sql.NullString{
			String: filepath,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, 0, err
	}

	return User{
		UserID:     user.UserID,
		Username:   user.Username,
		Name:       user.Name,
		ProfileUrl: user.ProfileUrl,
	}, http.StatusOK, nil
}

// swagger:route DELETE /me Users deleteUser
// Delete user of the auth token
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  204: noContent
func (h *handler) deleteUser(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	err := h.service.DeleteUser(c, userPayload.UserID)
	if err != nil {
		return nil, 0, err
	}

	return "", http.StatusNoContent, nil
}
