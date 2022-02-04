package db

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"strconv"
// 	"testing"

// 	"github.com/nemo984/money-app-api/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomIncome(t *testing.T, userID int32) Expense {
// 	arg := CreateIncomeParams{
// 		UserID: userID,
// 		IncomeTypeName:
// 	}

// 	expense, err := testQueries.CreateExpense(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, expense)

// 	require.Equal(t, arg.UserID, expense.UserID)
// 	require.Equal(t, arg.CategoryID, expense.CategoryID)
// 	argAmount, err := strconv.Atoi(arg.Amount)
// 	require.NoError(t, err)
// 	expAmount, err := strconv.ParseFloat(expense.Amount, 64)
// 	require.NoError(t, err)
// 	require.Equal(t, argAmount, int(expAmount))
// 	require.Equal(t, arg.Frequency, expense.Frequency)
// 	require.Equal(t, arg.Note, expense.Note)

// 	require.NotZero(t, expense.ExpenseID)
// 	require.NotZero(t, expense.CreatedAt)

// 	return expense
// }

// func TestCreateExpense(t *testing.T) {
// 	user := createRandomUser(t)
// 	createRandomExpense(t, user.UserID)
// }

// func TestDeleteExpense(t *testing.T) {
// 	user := createRandomUser(t)
// 	e := createRandomExpense(t, user.UserID)
// 	err := testQueries.DeleteExpense(context.Background(), e.ExpenseID)
// 	require.NoError(t, err)
// }