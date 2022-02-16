package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.service.GetCategories(c)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}