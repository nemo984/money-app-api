package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Server struct {
	db     db.Querier
	router *gin.Engine
}

func NewServer(db db.Querier) *Server {
	server := &Server{db: db}
	router := gin.Default()
	
	//setup routes func
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "cool")
	})
	router.POST("/users", server.createUser)

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H{
	return gin.H{
		"error": err.Error(),
	}
}