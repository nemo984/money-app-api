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

	apiRoute.POST("/token", handler.createUserToken)
	apiRoute.GET("/google-login", handler.GoogleLogin)
	apiRoute.GET("/google-callback", handler.GoogleCallback)

	users := apiRoute.Group("/users")
	{
		users.POST("", handler.createUser)
	}

	apiRoute.GET("/notifications-ws", handler.WSNotificationHandler)

	userRoute := apiRoute.Group("/me")
	userRoute.Use(handler.authenticatedToken())
	{
		userRoute.GET("", handler.getUser)
		userRoute.PATCH("", handler.updateUser)
		userRoute.DELETE("", handler.deleteUser)
		userRoute.PUT("/picture", handler.uploadProfilePicture)

		expensesRoute := userRoute.Group("/expenses")
		{
			expensesRoute.GET("", handler.getExpenses)
			expensesRoute.POST("", handler.createExpense)
			expensesRoute.DELETE("/:id", handler.deleteExpense)
		}

		budgetsRoute := userRoute.Group("/budgets")
		{
			budgetsRoute.GET("", handler.getBudgets)
			budgetsRoute.POST("", handler.createBudget)
			budgetsRoute.DELETE("/:id", handler.deleteBudget)
		}

		incomesRoute := userRoute.Group("/incomes")
		{
			incomesRoute.GET("", handler.getIncomes)
			incomesRoute.POST("", handler.createIncome)
			incomesRoute.DELETE("/:id", handler.deleteIncome)
		}

		notificationsRoute := userRoute.Group("/notifications")
		{
			notificationsRoute.GET("", handler.getNotifications)
			notificationsRoute.PATCH("", handler.updateAllNotifications)
			notificationsRoute.PATCH("/:id", handler.updateNotification)
		}
	}

	categoriesRoute := apiRoute.Group("/categories")
	{
		categoriesRoute.GET("", handler.getCategories)
	}

	incomeTypesRoute := apiRoute.Group("/income-types")
	{
		incomeTypesRoute.GET("", handler.getIncomeTypes)
	}

	handler.router = router
	return handler
}

func (s *handler) Start(addr string) error {
	go s.hub.Run()
	return s.router.Run(addr)
}

func handleError(c *gin.Context, err error) {
	switch v := err.(type) {
	case service.AppError:
		c.JSON(v.StatusCode, errorResponse(v.Err))
	case error:
		c.JSON(http.StatusInternalServerError, errorResponse(v))
	}
}

func handleValidationError(c *gin.Context, t interface{}, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := util.ListOfErrors(t, ve)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errs,
		})
		return
	}
	c.JSON(http.StatusBadRequest, errorResponse(err))
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
