package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
)

// A list of categories
// swagger:response categoriesResponse
type categoriesResponse struct {
	// Categories
	// in:body
	Body []db.Category
}

// swagger:route GET /categories Categories listCategories
// Returns a list of categories
// responses:
//  200: categoriesResponse
func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.service.GetCategories(c)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}
