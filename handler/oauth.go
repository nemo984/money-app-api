package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

type UserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (s *Server) GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		c.Redirect(http.StatusTemporaryRedirect, "/api/auth-failed")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/api/auth-failed")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println("get failed oauth2:", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/api/auth-failed")
		return
	}

	defer response.Body.Close()
	var userInfo UserInfo
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		log.Println("json decode error:", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/api/auth-failed")
		return
	}
	//user is created if doesn't exists
	user, err := s.service.CreateUser(c, db.CreateUserParams{
		Username: userInfo.Email,
		Name: sql.NullString{
			String: userInfo.Name,
			Valid: true,
		},
		Password: userInfo.Name, //TODO: later
		ProfileUrl: sql.NullString{
			String: userInfo.Picture,
			Valid:  true,
		},
	})
	if err != nil && !errors.Is(err, service.ErrUserAlreadyExists) {
		handleError(c, err)
		return
	}
	
	jwtToken, err := service.CreateToken(user.UserID)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}
