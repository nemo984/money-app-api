package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

// A list of budgets
// swagger:response budgetsResponse
type budgetsResponse struct {
	// User's budgets
	// in:body
	Body []db.Budget
}

// Budget
// swagger:response budgetResponse
type budgetResponse struct {
	Body db.Budget
}

// swagger:route GET /me/budgets Budgets listBudgets
// Returns a list of budgets of the user
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: budgetsResponse
func (s *Server) getBudgets(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	budgets, err := s.service.GetBudgets(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// swagger:parameters createBudget
type CreateBudgetRequest struct {
	// The budget to create
	//
	// required: true
	// in: body
	Budget *createBudgetRequest `json:"budget"`
}

// swagger:model
type createBudgetRequest struct {
	// id of a category
	// required: true
	// min: 1
	CategoryID int32 `json:"category_id" binding:"required,min=1"`
	// amount of the budget
	// required: true
	// min: 1
	Amount string `json:"amount" binding:"required,min=1"`
	// Numbers of days to budget
	// required: true
	// min: 1
	Days int `json:"days" binding:"required,min=1"`
}

// swagger:route POST /me/budgets Budgets createBudget
// Returns the created budget
//
// Consumes:
//  - application/json
//
// Produces:
//	-application/json
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  201: budgetResponse
//  422: validationError
func (s *Server) createBudget(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req createBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	args := db.CreateBudgetParams{
		UserID:     userPayload.UserID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		EndDate: sql.NullTime{
			Time:  time.Now().AddDate(0, 0, req.Days),
			Valid: true,
		},
	}
	budget, err := s.service.CreateBudget(c, args)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// swagger:parameters deleteBudget
type deleteBudgetURI struct {
	// The id of the budget to delete from the database
	// in: path
	// required: true
	// min: 1
	BudgetID int32 `uri:"id" binding:"min=1"`
}

// swagger:route DELETE /me/budgets/{id} Budgets deleteBudget
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  204: noContent
func (s *Server) deleteBudget(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var uri deleteBudgetURI
	if err := c.ShouldBindUri(&uri); err != nil {
		handleError(c, err)
		return
	}

	if err := s.service.DeleteBudget(c, userPayload.UserID, uri.BudgetID); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
