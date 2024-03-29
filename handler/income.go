package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/docs"
	"github.com/nemo984/money-app-api/service"
)

// A list of incomes
// swagger:response incomesResponse
type incomesResponse struct {
	// User's expenses
	// in:body
	Body []db.Income
}

// Income
// swagger:response incomeResponse
type incomeResponse struct {
	Body db.Income
}

// swagger:route GET /me/incomes Incomes listIncomes
// Returns a list of incomes of the user
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: incomesResponse
func (h *handler) getIncomes(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	incomes, err := h.service.GetIncomes(c, userPayload.UserID)
	if err != nil {
		return nil, 0, err
	}

	return incomes, http.StatusOK, nil
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
	// min: 1
	IncomeTypeID int32 `json:"income_type_id" binding:"required,min=1"`
	// description of for the income
	//
	// maximum length: 255
	Description string `json:"description" binding:"max=255"`
	// amount of the income
	//
	// required: true
	// min: 1
	Amount string `json:"amount" binding:"required,min=1"`
	// frequency of the income
	Frequency docs.DateFrequency `json:"frequency"`
}

// swagger:route POST /me/incomes Incomes createIncome
// Returns the created income
//
// Consumes:
//  - application/json
//
// Produces:
//	- application/json
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  201: incomeResponse
//  422: validationError
func (h *handler) createIncome(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req createIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := validationErrors(&req, err)
		return nil, 0, err
	}

	args := db.CreateIncomeParams{
		UserID:       userPayload.UserID,
		Amount:       req.Amount,
		IncomeTypeID: req.IncomeTypeID,
		Description: sql.NullString{
			String: req.Description,
			Valid:  true,
		},
		Frequency: db.DateFrequency(req.Frequency),
	}
	income, err := h.service.CreateIncome(c, args)
	if err != nil {
		return nil, 0, err
	}

	return income, http.StatusCreated, nil
}

// swagger:parameters deleteIncome
type deleteIncomeURI struct {
	// The id of the income to delete from the database
	// in: path
	// required: true
	// min: 1
	IncomeID int32 `uri:"id" binding:"min=1"`
}

// swagger:route DELETE /me/incomes/{id} Incomes deleteIncome
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  204: noContent
func (h *handler) deleteIncome(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var uri deleteIncomeURI
	if err := c.ShouldBindUri(&uri); err != nil {
		err := validationErrors(&uri, err)
		return nil, 0, err
	}

	if err := h.service.DeleteIncome(c, userPayload.UserID, uri.IncomeID); err != nil {
		return nil, 0, err
	}

	return "", http.StatusNoContent, nil
}

// A list of income types
// swagger:response incomeTypesResponse
type incomeTypesResponse struct {
	// List of income types
	// in:body
	Body []db.IncomeType
}

// swagger:route GET /income-types Incomes listIncomeTypes
// List the available income types
// responses:
//  200: incomeTypesResponse
func (h *handler) getIncomeTypes(c *gin.Context) (interface{}, int, error) {
	types, err := h.service.GetIncomeTypes(c)
	if err != nil {
		return nil, 0, err
	}

	return types, http.StatusOK, nil
}
