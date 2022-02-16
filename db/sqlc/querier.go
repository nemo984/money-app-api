// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error)
	CreateCategory(ctx context.Context, name string) (Category, error)
	CreateExpense(ctx context.Context, arg CreateExpenseParams) (Expense, error)
	CreateIncome(ctx context.Context, arg CreateIncomeParams) (Income, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteBudget(ctx context.Context, budgetID int32) error
	DeleteCategory(ctx context.Context, categoryID int32) error
	DeleteExpense(ctx context.Context, userID int32) error
	DeleteIncome(ctx context.Context, userID int32) error
	DeleteNotification(ctx context.Context, userID int32) error
	DeleteUser(ctx context.Context, userID int32) error
	GetBudget(ctx context.Context, budgetID int32) (Budget, error)
	GetBudgets(ctx context.Context, userID int32) ([]Budget, error)
	GetCategories(ctx context.Context) ([]Category, error)
	GetExpense(ctx context.Context, expenseID int32) (Expense, error)
	GetExpenses(ctx context.Context, userID int32) ([]Expense, error)
	GetIncome(ctx context.Context, incomeID int32) (Income, error)
	GetIncomes(ctx context.Context, userID int32) ([]Income, error)
	GetNotifications(ctx context.Context, userID int32) ([]Notification, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserByID(ctx context.Context, userID int32) (User, error)
	UpdateBudget(ctx context.Context, arg UpdateBudgetParams) error
	UpdateNotification(ctx context.Context, arg UpdateNotificationParams) (Notification, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
