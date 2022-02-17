package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

func (s *Server) getIncomes(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
	incomes, err := s.service.GetIncomes(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, incomes)
}

type createIncomeRequest struct {
	IncomeTypeID int32  `json:"income_type_id" binding:"required"`
	Description  string `json:"description"`
	Amount       string `json:"amount" binding:"required"`
	// Frequency    DateFrequency  `json:"frequency"`` later
}

func (s *Server) createIncome(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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
		Frequency: db.DateFrequencyMonth, // later
	}
	income, err := s.service.CreateIncome(c, args)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, income)
}

type deleteIncomeURI struct {
	IncomeID int32 `uri:"id"`
}

func (s *Server) deleteIncome(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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

func (s *Server) getIncomeTypes(c *gin.Context) {
	types, err := s.service.GetIncomeTypes(c)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, types)
}
