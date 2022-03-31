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
func (h *handler) getCategories(c *gin.Context) (interface{}, int, error) {
	categories, err := h.service.GetCategories(c)
	if err != nil {
		return nil, 0, err
	}

	return categories, http.StatusOK, nil
}
