package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

func (s *Server) getExpenses(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	expenses, err := s.service.GetExpenses(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}

type createExpenseRequest struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
	// Frequency  DateFrequency  `json:"frequency"` : later dude
	Note string `json:"note"`
}

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
		Frequency:  db.DateFrequencyMonth, // later
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

type deleteExpenseURI struct {
	ExpenseID int32 `uri:"id"`
}

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
