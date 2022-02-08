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

	//setup routes func
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "cool")
	})
	router.POST("/users", server.createUser)
	router.DELETE("/users", server.deleteUser)

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
