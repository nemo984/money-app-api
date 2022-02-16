package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

func (s *Server) getBudgets(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	budgets, err := s.service.GetBudgets(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, budgets)
}

type createBudgetRequest struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
	Days       int    `json:"days" binding:"required"`
}

func (s *Server) createBudget(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	var req createBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
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

type deleteBudgetURI struct {
	BudgetID int32 `uri:"id"`
}

func (s *Server) deleteBudget(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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
