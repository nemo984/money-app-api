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

	users := router.Group("/users")
	{
		users.POST("", server.createUser)
		users.POST("/token", server.createUserToken)
		users.PATCH("/:id", server.updateUser)
		users.DELETE("/:id", server.deleteUser)
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
