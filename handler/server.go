package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nemo984/money-app-api/service"
)

type Server struct {
	service service.Service
	router  *gin.Engine
}

func NewServer(service service.Service) *Server {
	server := &Server{service: service}
	router := gin.Default()
	apiRoute := router.Group("/api")

	users := apiRoute.Group("/users")
	{
		users.POST("", server.createUser)
		users.POST("/token", server.createUserToken)
	}

	apiRoute.GET("/google-login", server.GoogleLogin)
	apiRoute.GET("/google-callback", server.GoogleCallback)

	userRoute := apiRoute.Group("/me")
	{
		userRoute.Use(authenticatedToken())
		userRoute.GET("", server.getUser)
		userRoute.PATCH("", server.updateUser)
		userRoute.DELETE("", server.deleteUser)
	}

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func handleError(c *gin.Context, err error) {
	switch v := err.(type) {
	case service.AppError:
		c.JSON(v.StatusCode, errorResponse(v.Err))
	case error:
		c.JSON(http.StatusInternalServerError, errorResponse(v))
	}
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
