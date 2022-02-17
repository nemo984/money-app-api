package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Service interface {
	GetUser(context.Context, int32) (db.User, error)
	CreateUser(context.Context, db.CreateUserParams) (db.User, error)
	LoginUser(ctx context.Context, username string, password string) (token string, err error)
	VerifyToken(ctx context.Context, token string) (JWTClaims, error)
	DeleteUser(context.Context, int32) error
	UpdateUser(context.Context, UpdateUserParams) (db.User, error)

	GetCategories(context.Context) ([]db.Category, error)
	CreateCategory(context.Context, string) (db.Category, error)
	DeleteCategory(context.Context, int32) error

	GetExpenses(context.Context, int32) ([]db.Expense, error)
	CreateExpense(context.Context, db.CreateExpenseParams) (db.Expense, error)
	DeleteExpense(ctx context.Context, userID, expenseID int32) error

	GetBudgets(ctx context.Context, userID int32) ([]db.Budget, error)
	CreateBudget(ctx context.Context, args db.CreateBudgetParams) (db.Budget, error)
	DeleteBudget(ctx context.Context, userID, budgetID int32) error

	GetIncomes(ctx context.Context, userID int32) ([]db.Income, error)
	CreateIncome(ctx context.Context, args db.CreateIncomeParams) (db.Income, error)
	DeleteIncome(ctx context.Context, userID, incomeID int32) error

	GetIncomeTypes(ctx context.Context) ([]db.IncomeType, error)
	CreateIncomeType(ctx context.Context, name string) (db.IncomeType, error)
	DeleteIncomeType(ctx context.Context, id int32) error

	GetNotifications(ctx context.Context, userID int32) ([]db.Notification, error)
	CreateNotification(ctx context.Context, args db.CreateNotificationParams) (db.Notification, error)
	UpdateNotification(ctx context.Context, userID int32, args db.UpdateNotificationParams) (db.Notification, error)
	UpdateNotifications(ctx context.Context, args db.UpdateNotificationsParams) ([]db.Notification, error)
}

type service struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &service{db: db}
}

type AppError struct {
	StatusCode int
	Err        error
}

func (a AppError) Error() string {
	return a.Err.Error()
}
