package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/docs"
	"github.com/nemo984/money-app-api/service"
)

// A list of expenses
// swagger:response expensesResponse
type expensesResponse struct {
	// User's expenses
	// in:body
	Body []db.Expense
}

// Expense
// swagger:response expenseResponse
type expenseResponse struct {
	Body db.Expense
}

// swagger:route GET /me/expenses Expenses listExpenses
// Returns a list of expenses of the user
// responses:
//  200: expensesResponse
func (s *Server) getExpenses(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	expenses, err := s.service.GetExpenses(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// swagger:parameters createExpense
type CreateExpenseRequest struct {
	// The expense to create
	//
	// required: true
	// in: body
	Expense *createExpenseRequest `json:"expense"`
}

// swagger:model
type createExpenseRequest struct {
	// id of a category
	//
	// required: true
	// min: 1
	CategoryID int32 `json:"category_id" binding:"required,min=1"`
	// amount of the
	//
	// required: true
	// min: 1
	Amount string `json:"amount" binding:"required,min=1"`
	// frequency of the expense
	Frequency docs.DateFrequency `json:"frequency"`
	// note of the expense
	Note string `json:"note"`
}

// swagger:route POST /me/expenses Expenses createExpense
// Returns the created expense
//
// Consumes:
//  - application/json
//
// Produces:
//	- application/json
//
// responses:
//  201: expenseResponse
//  422: validationError
func (s *Server) createExpense(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req createExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateExpenseParams{
		UserID:     userPayload.UserID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Frequency:  db.DateFrequency(req.Frequency),
		Note: sql.NullString{
			String: req.Note,
			Valid:  true,
		},
	}
	expense, err := s.service.CreateExpense(c, args)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// swagger:parameters deleteExpense
type deleteExpenseURI struct {
	// The id of the expense to delete from the database
	// in: path
	// required: true
	// min: 1
	ExpenseID int32 `uri:"id" binding:"min=1"`
}

// swagger:route DELETE /me/expenses/{id} Expenses deleteExpense
// responses:
//  204: noContent
func (s *Server) deleteExpense(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var uri deleteExpenseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		handleError(c, err)
		return
	}

	if err := s.service.DeleteExpense(c, userPayload.UserID, uri.ExpenseID); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
