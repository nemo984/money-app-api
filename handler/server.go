package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nemo984/money-app-api/service"
)

type Server struct {
	service service.Service
	router  *gin.Engine
}

func NewServer(service service.Service) *Server {
	server := &Server{service: service}
	router := gin.Default()
	apiRoute := router.Group("/api")

	apiRoute.GET("/google-login", server.GoogleLogin)
	apiRoute.GET("/google-callback", server.GoogleCallback)

	users := apiRoute.Group("/users")
	{
		users.POST("", server.createUser)
		users.POST("/token", server.createUserToken)
	}

	userRoute := apiRoute.Group("/me")
	{
		userRoute.Use(authenticatedToken())
		userRoute.GET("", server.getUser)
		userRoute.PATCH("", server.updateUser)
		userRoute.DELETE("", server.deleteUser)

		expensesRoute := userRoute.Group("/expenses")
		{
			expensesRoute.GET("", server.getExpenses)
			expensesRoute.POST("", server.createExpense)
			expensesRoute.DELETE("/:id", server.deleteExpense)
		}

		budgetsRoute := userRoute.Group("/budgets")
		{
			budgetsRoute.GET("", server.getBudgets)
			budgetsRoute.POST("", server.createBudget)
			budgetsRoute.DELETE("/:id", server.deleteBudget)
		}
	}

	categoriesRoute := apiRoute.Group("/categories")
	{
		categoriesRoute.GET("", server.getCategories)
	}

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
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

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
