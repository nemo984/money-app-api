package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

// A list of incomes
// swagger:response incomesResponse
type incomesResponse struct {
	// User's expenses
	// in:body
	Body []db.Expense
}

// Income
// swagger:response incomeResponse
type incomeResponse struct {
	Body db.Budget
}

// swagger:route GET /me/incomes incomes listIncomes
// Returns a list of incomes of the user
// responses:
//  200: incomesResponse
func (s *Server) getIncomes(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	incomes, err := s.service.GetIncomes(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, incomes)
}

// swagger:parameters createIncome
type CreateIncomeRequest struct {
	// The budget to create
	//
	// required: true
	// in: body
	Income *createIncomeRequest `json:"income"`
}

// swagger:model
type createIncomeRequest struct {
	// id of the type of income
	//
	// required: true
	IncomeTypeID int32 `json:"income_type_id" binding:"required"`
	// description of for the income
	Description string `json:"description"`
	// amount of the income
	//
	// required: true
	Amount string `json:"amount" binding:"required"`
	// frequency of the income
	Frequency db.DateFrequency `json:"frequency"`
}

// swagger:route POST /me/incomes incomes createIncome
// Returns the created income
//
// Consumes:
//  - application/json
//
// Produces:
//	- application/json
//
// responses:
//  201: incomeResponse
//  422: validationError
func (s *Server) createIncome(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req createIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateIncomeParams{
		UserID:       userPayload.UserID,
		Amount:       req.Amount,
		IncomeTypeID: req.IncomeTypeID,
		Description: sql.NullString{
			String: req.Description,
			Valid:  true,
		},
		Frequency: req.Frequency,
	}
	income, err := s.service.CreateIncome(c, args)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, income)
}

// swagger:parameters deleteIncome
type deleteIncomeURI struct {
	// The id of the income to delete from the database
	// in: path
	// required: true
	IncomeID int32 `uri:"id"`
}

// swagger:route DELETE /me/incomes/{id} incomes deleteIncome
// responses:
//  204: noContent
func (s *Server) deleteIncome(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var uri deleteIncomeURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := s.service.DeleteIncome(c, userPayload.UserID, uri.IncomeID); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// A list of income types
// swagger:response incomeTypesResponse
type incomeTypesResponse struct {
	// List of income types
	// in:body
	Body []db.IncomeType
}

// swagger:route GET /income-types incomes listIncomeTypes
// List the available income types
// responses:
//  200: incomeTypesResponse
func (s *Server) getIncomeTypes(c *gin.Context) {
	types, err := s.service.GetIncomeTypes(c)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, types)
}
