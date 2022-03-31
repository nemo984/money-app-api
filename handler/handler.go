package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-playground/validator/v10"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/notification"
	"github.com/nemo984/money-app-api/service"
	"github.com/nemo984/money-app-api/util"
)

type handler struct {
	hub     NotificationHub
	service service.Service
	router  *gin.Engine
}

type NotificationHub interface {
	Run()
	Notify(userID int32, notification db.Notification)
	Register(user *notification.User)
	Unregister(user *notification.User)
}

func New(service service.Service, hub NotificationHub) *handler {
	handler := &handler{service: service, hub: hub}
	router := gin.Default()
	apiRoute := router.Group("/api")

	//documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml", Title: "Money App API Documentation"}
	sh := middleware.Redoc(opts, nil)
	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/docs", gin.WrapH(sh))

	apiRoute.StaticFS("/images/user-profile-pics", util.Fs{Dir: http.Dir("./images/user-profile-pics")})

	apiRoute.POST("/token", responseHandler(handler.createUserToken))
	apiRoute.GET("/google-login", handler.GoogleLogin)
	apiRoute.GET("/google-callback", handler.GoogleCallback)

	users := apiRoute.Group("/users")
	{
		users.POST("", responseHandler(handler.createUser))
	}

	apiRoute.GET("/notifications-ws", responseHandler(handler.WSNotificationHandler))

	userRoute := apiRoute.Group("/me")
	userRoute.Use(handler.authenticatedToken())
	{
		userRoute.GET("", responseHandler(handler.getUser))
		userRoute.PATCH("", responseHandler(handler.updateUser))
		userRoute.DELETE("", responseHandler(handler.deleteUser))
		userRoute.PUT("/picture", responseHandler(handler.uploadProfilePicture))

		expensesRoute := userRoute.Group("/expenses")
		{
			expensesRoute.GET("", responseHandler(handler.getExpenses))
			expensesRoute.POST("", responseHandler(handler.createExpense))
			expensesRoute.DELETE("/:id", responseHandler(handler.deleteExpense))
		}

		budgetsRoute := userRoute.Group("/budgets")
		{
			budgetsRoute.GET("", responseHandler(handler.getBudgets))
			budgetsRoute.POST("", responseHandler(handler.createBudget))
			budgetsRoute.DELETE("/:id", responseHandler(handler.deleteBudget))
		}

		incomesRoute := userRoute.Group("/incomes")
		{
			incomesRoute.GET("", responseHandler(handler.getIncomes))
			incomesRoute.POST("", responseHandler(handler.createIncome))
			incomesRoute.DELETE("/:id", responseHandler(handler.deleteIncome))
		}

		notificationsRoute := userRoute.Group("/notifications")
		{
			notificationsRoute.GET("", responseHandler(handler.getNotifications))
			notificationsRoute.PATCH("", responseHandler(handler.updateAllNotifications))
			notificationsRoute.PATCH("/:id", responseHandler(handler.updateNotification))
		}
	}

	categoriesRoute := apiRoute.Group("/categories")
	{
		categoriesRoute.GET("", responseHandler(handler.getCategories))
	}

	incomeTypesRoute := apiRoute.Group("/income-types")
	{
		incomeTypesRoute.GET("", responseHandler(handler.getIncomeTypes))
	}

	handler.router = router
	return handler
}

func (h *handler) Start(addr string) error {
	go h.hub.Run()
	return h.router.Run(addr)
}

func responseHandler(h func(*gin.Context) (interface{}, int, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, status, err := h(c)
		if err != nil {
			switch v := err.(type) {
			case service.AppError:
				c.JSON(v.StatusCode, errorResponse(v.Err))
			case validationError:
				c.JSON(http.StatusUnprocessableEntity, v.errors)
			case error:
				code := http.StatusInternalServerError
				if status != 0 {
					code = status
				}
				c.JSON(code, errorResponse(v))
			}
			return
		}
		c.JSON(status, data)
	}
}

//validationError contains a map of strings for each attribute that fails to validate with a cause
type validationError struct {
	errors []map[string]string
}

func (v validationError) Error() string {
	return ""
}

func validationErrors(t interface{}, err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := util.ListOfErrors(t, ve)
		return validationError{
			errors: errs,
		}
	}
	return err
}

func handleError(c *gin.Context, err error) {
	switch v := err.(type) {
	case service.AppError:
		c.JSON(v.StatusCode, errorResponse(v.Err))
	case error:
		c.JSON(http.StatusInternalServerError, errorResponse(v))
	}
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
